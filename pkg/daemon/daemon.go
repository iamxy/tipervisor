package daemon

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/pingcap/tipervisor/pkg/sink"
	"github.com/pingcap/tipervisor/pkg/util"
	"github.com/pingcap/tipervisor/pkg/util/log"
	"github.com/pkg/errors"
)

// Daemon is a single process manager that controls the process start or stop
// and supervises the process running
type Daemon struct {
	lock, lockOnce uint32
	proc           *process
	state          ProcessState
	mu             sync.Mutex
	runch          chan struct{}
	sigch          chan SignalRequest
	runStat        *RunStat
	cfg            *Config
	logSinkFactory sink.LogSinkFactory
}

// New creates a new daemon instance
func New(cfg *Config, lsf sink.LogSinkFactory) (*Daemon, error) {
	var err error
	if err = checkRunning(cfg); err != nil {
		return nil, err
	}
	if err = checkUser(cfg); err != nil {
		return nil, err
	}
	d := &Daemon{
		state:          ProcStatStopped,
		runch:          make(chan struct{}, 1),
		sigch:          make(chan SignalRequest),
		runStat:        &RunStat{},
		cfg:            cfg,
		logSinkFactory: lsf,
	}
	return d, nil
}

func checkRunning(cfg *Config) error {
	if cfg.Name == "" {
		return errors.New("daemon name can not be empty")
	}
	if cfg.StatusDir == "" {
		return errors.New("daemon status dir can not be empty")
	}
	if !util.IsDir(cfg.StatusDir) {
		return errors.Errorf("daemon status dir [%s] not exists", cfg.StatusDir)
	}
	cfg.pidfile = filepath.Join(cfg.StatusDir, fmt.Sprintf("%s.pid", cfg.Name))
	if util.IsFile(cfg.pidfile) {
		var (
			pid int
			err error
		)
		if pid, err = readPidFile(cfg.pidfile); err != nil {
			return err
		}
		// ensure that no process is running
		if isRunning(pid) {
			if err := syscall.Kill(pid, syscall.SIGKILL); err != nil {
				return errors.Wrapf(err, "sending signal failed")
			}
		}
	}
	return nil
}

func checkUser(cfg *Config) error {
	if cfg.User == "" {
		// cfg.user = nil
		return nil
	}
	usr, err := user.Lookup(cfg.User)
	if err != nil {
		if _, ok := err.(user.UnknownUserError); ok {
			return errors.Errorf("user [%s] not exists", cfg.User)
		}
		return errors.Wrapf(err, "error looking up user [%s]", cfg.User)
	}
	cfg.user = usr
	return nil
}

// ReadPidFile read pid from file if error returns pid 0
func readPidFile(pidfile string) (int, error) {
	data, err := ioutil.ReadFile(pidfile)
	if err != nil {
		return 0, errors.Wrap(err, "read pid file failed")
	}
	lines := strings.Split(string(data), "\n")
	pid, err := strconv.Atoi(lines[0])
	if err != nil {
		return 0, errors.Wrap(err, "convert pid failed")
	}
	return pid, nil
}

func writePidFile(pidfile string, pid int) error {
	data := []byte(fmt.Sprintf("%d", pid))
	if err := ioutil.WriteFile(pidfile, data, 0644); err != nil {
		return errors.Wrapf(err, "write pid file failed")
	}
	return nil
}

func isRunning(pid int) bool {
	// On Unix systems, FindProcess always succeeds and returns a
	// Process for the given pid, regardless of whether the process exists.
	proc, _ := os.FindProcess(pid)
	if err := proc.Signal(syscall.Signal(0)); err != nil {
		return false
	}
	return true
}

func (d *Daemon) newProcess() *process {
	// initialize new process struct
	p := &process{
		Config:  d.cfg,
		logSink: d.logSinkFactory.NewLogSink(),
	}
	return p
}

func (d *Daemon) run() error {
	var err error

	// if already locked, not need to start new process
	if atomic.SwapUint32(&d.lock, uint32(1)) != 0 {
		return nil
	}

	p := d.newProcess()
	if err = p.Start(); err != nil {
		return err
	}

	if err = writePidFile(d.cfg.pidfile, p.Pid()); err != nil {
		if kerr := p.MustKill(); kerr != nil {
			log.Warnf("%+v", kerr)
		}
		return err
	}

	d.proc = p
	d.changeToState(ProcStatRunning)
	d.runStat.Lock()
	// increase the counter of run times
	d.runStat.RunCount++
	d.runStat.StartTime = p.sTime
	d.runStat.Pid = p.Pid()
	d.runStat.Unlock()

	return nil
}

// Supervise supervises the process running,
// when detects the process exited abnormally, starts it immediately
func (d *Daemon) Supervise(ctx context.Context) {
	go func(ctx context.Context) {
		err := d.supervise(ctx)
		if err != nil {
			log.WithField("daemon", d.cfg.Name).Errorf("supervise error exit: %+v", err)
		}
	}(ctx)
	// sleep for one second, waiting for process to run up
	time.Sleep(1 * time.Second)
}

func (d *Daemon) supervise(ctx context.Context) error {
	var (
		err error
		sig SignalRequest
	)

	if err = d.run(); err != nil {
		// reset lock
		atomic.StoreUint32(&d.lock, d.lockOnce)
		return err
	}

	for {
		select {
		case <-d.runch:
			if err = d.run(); err != nil {
				// reset lock
				atomic.StoreUint32(&d.lock, d.lockOnce)
				return err
			}
		case sig = <-d.sigch:
			d.handleSignal(sig)
		case <-ctx.Done():
			s := d.ProcessState()
			switch s {
			case ProcStatRunning:
				d.changeToState(ProcStatTerminating)
				d.lockOnce = 1
				err = d.proc.Kill()
			case ProcStatStopping, ProcStatRestarting, ProcStatKilling:
				d.changeToState(ProcStatTerminating)
				d.lockOnce = 1
			default:
				// exit normally
				return nil
			}
		case perr := <-d.proc.errch:
			d.runStat.Lock()
			d.runStat.LastStartTime = d.proc.sTime
			d.runStat.LastEndTime = d.proc.eTime
			d.runStat.LastUpTime = d.proc.eTime.Sub(d.proc.sTime)
			d.runStat.LastUserTime = d.proc.cmd.ProcessState.UserTime()
			d.runStat.LastSysTime = d.proc.cmd.ProcessState.SystemTime()
			d.runStat.LastExitErr = perr
			d.runStat.Pid = 0
			d.runStat.Unlock()

			s := d.ProcessState()
			switch s {
			case ProcStatStopping, ProcStatRestarting:
				d.changeToState(ProcStatStopped)
				d.runStat.Lock()
				d.runStat.StoppedCount++
				d.runStat.LastTerminateState = ProcStatStopped
				d.runStat.Unlock()
			case ProcStatKilling:
				d.changeToState(ProcStatKilled)
				d.runStat.Lock()
				d.runStat.KilledCount++
				d.runStat.LastTerminateState = ProcStatKilled
				d.runStat.Unlock()
			case ProcStatRunning:
				d.changeToState(ProcStatExited)
				d.runStat.Lock()
				d.runStat.ExitedCount++
				d.runStat.LastTerminateState = ProcStatExited
				d.runStat.Unlock()
			case ProcStatTerminating:
				// exit normally
				return nil
			default:
				return errors.Errorf("process exit from unexpected state: %v", s)
			}

			// if lockOnce is not true, the process will be automatically started
			if d.lockOnce == 0 {
				d.changeToState(ProcStatStarting)
			}
			// reset lock, if lockOnce is true, the process will not automatically restart next time
			atomic.StoreUint32(&d.lock, d.lockOnce)
			d.runch <- struct{}{}
		}
	}
}
