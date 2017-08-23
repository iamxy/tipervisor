package daemon

import (
	"fmt"
	"io/ioutil"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
)

const minRunInterval time.Duration = 5 * time.Second

type demonRunStat struct {
	sync.RWMutex
	lastStartTime      time.Time
	lastEndTime        time.Time
	lastUptime         time.Duration
	lastUserTime       time.Duration
	lastSysTime        time.Duration
	lastTerminateState ProcessState
	lastExitErr        error
	startTime          time.Time
	runCount           uint32
	stoppedCount       uint32
	restartedCount     uint32
	exitedCount        uint32
	killedCount        uint32
}

// Daemon is a single process manager that controls the process start or stop
// and supervises the process running
type Daemon struct {
	pidFilePath    string
	lock, lockOnce uint32
	proc           *process
	state          ProcessState
	mu             sync.Mutex
	runch          chan struct{}
	sigch          chan SignalRequest
	runStat        demonRunStat
}

// New do ...
func New() (*Daemon, error) {
	return nil, nil
}

func (d *Daemon) newProcess() *process {
	// initialize new process struct
	p := &process{}

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

	pid := fmt.Sprintf("%d", p.Pid())
	if err = ioutil.WriteFile(d.pidFilePath, []byte(pid), 0644); err != nil {
		if kerr := p.MustKill(); kerr != nil {
			// logger warn
		}
		return errors.Wrap(err, "write pid file failed")
	}

	d.proc = p
	d.changeToState(ProcStatRunning)
	d.runStat.Lock()
	// increase the counter of run times
	d.runStat.runCount++
	d.runStat.startTime = p.sTime
	d.runStat.Unlock()

	return nil
}

// Supervise supervises the process running,
// when detects the process exited abnormally, starts it immediately
func (d *Daemon) Supervise() error {
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
		case perr := <-d.proc.errch:
			uptime := d.proc.eTime.Sub(d.proc.sTime)
			if uptime < minRunInterval {
				wait := minRunInterval - uptime
				time.Sleep(wait)
			}

			d.runStat.Lock()
			d.runStat.lastStartTime = d.proc.sTime
			d.runStat.lastEndTime = d.proc.eTime
			d.runStat.lastUptime = uptime
			d.runStat.lastUserTime = d.proc.cmd.ProcessState.UserTime()
			d.runStat.lastSysTime = d.proc.cmd.ProcessState.SystemTime()
			d.runStat.lastExitErr = perr
			d.runStat.Unlock()

			s := d.ProcessState()
			if s == ProcStatStopping || s == ProcStatRestarting {
				d.changeToState(ProcStatStopped)
				d.runStat.Lock()
				d.runStat.stoppedCount++
				d.runStat.lastTerminateState = ProcStatStopped
				d.runStat.Unlock()
			} else if s == ProcStatKilling {
				d.changeToState(ProcStatKilled)
				d.runStat.Lock()
				d.runStat.killedCount++
				d.runStat.lastTerminateState = ProcStatKilled
				d.runStat.Unlock()
			} else {
				// from ProcStatRunning state
				d.changeToState(ProcStatExited)
				d.runStat.Lock()
				d.runStat.exitedCount++
				d.runStat.lastTerminateState = ProcStatExited
				d.runStat.Unlock()
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
