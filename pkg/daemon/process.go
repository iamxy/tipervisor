package daemon

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"time"

	"github.com/pingcap/tipervisor/pkg/sink"
	"github.com/pingcap/tipervisor/pkg/util/log"
	"github.com/pkg/errors"
)

const minRunningDuration time.Duration = 5 * time.Second

// Process defines operations with process
type Process interface {
	Pid() int
	Start() error
	Wait() error
	Signal(syscall.Signal) error
	Kill() error
	MustKill() error
}

type process struct {
	*Config
	logSink sink.LogSink

	cmd          *exec.Cmd
	errch        chan error
	sTime, eTime time.Time
}

// Start runs the process
func (p *process) Start() error {
	p.cmd = exec.Command(p.Cmd, p.Args...)

	if p.Cwd != "" {
		p.cmd.Dir = p.Cwd
	}

	if p.Env != nil {
		env := os.Environ()
		for k, v := range p.Env {
			env = append(env, fmt.Sprintf("%s=%s", k, v))
		}
		p.cmd.Env = env
	}

	sysProcAttr := new(syscall.SysProcAttr)

	if p.user != nil {
		uid, err := strconv.Atoi(p.user.Uid)
		if err != nil {
			return errors.Wrap(err, "invalid uid")
		}
		gid, err := strconv.Atoi(p.user.Gid)
		if err != nil {
			return errors.Wrap(err, "invalid gid")
		}
		sysProcAttr.Credential = &syscall.Credential{
			Uid: uint32(uid),
			Gid: uint32(gid),
		}
	}

	sysProcAttr.Setpgid = true
	sysProcAttr.Pgid = 0
	p.cmd.SysProcAttr = sysProcAttr

	var (
		prOut, pwOut *os.File
		prErr, pwErr *os.File
		e            error
	)

	if p.logSink != nil {
		prOut, pwOut, e = os.Pipe()
		if e != nil {
			return errors.Wrap(e, "create stdout pipe failed")
		}

		prErr, pwErr, e = os.Pipe()
		if e != nil {
			return errors.Wrap(e, "create stderr pipe failed")
		}

		p.cmd.Stdout = pwOut
		p.cmd.Stderr = pwErr
		p.logSink.Start(prOut, prErr)
	}

	if err := p.cmd.Start(); err != nil {
		return errors.Wrap(err, "start process failed")
	}

	p.sTime = time.Now()
	p.errch = make(chan error, 1)

	go func() {
		err := p.cmd.Wait()
		p.eTime = time.Now()
		// prevent the process from being restarted repeatedly
		uptime := p.eTime.Sub(p.sTime)
		if uptime < minRunningDuration {
			wait := minRunningDuration - uptime
			time.Sleep(wait)
		}
		// stop log sink
		if p.logSink != nil {
			p.logSink.Stop()
		}
		// process exit notification
		p.errch <- err
	}()

	return nil
}

// Wait is waiting for process to end, and return its exit code
// return nil if exit code is zero
func (p *process) Wait() error {
	if p.errch == nil {
		return nil
	}
	err := <-p.errch
	return err
}

// Pid return process pid
func (p *process) Pid() int {
	if p.cmd == nil || p.cmd.Process == nil {
		return 0
	}
	return p.cmd.Process.Pid
}

// Signal sends a signal to the process
func (p *process) Signal(sig syscall.Signal) error {
	err := syscall.Kill(p.cmd.Process.Pid, sig)
	return errors.Wrapf(err, "sending signal failed")
}

// Kill the entire process group
func (p *process) Kill() error {
	processGroup := 0 - p.cmd.Process.Pid
	err := syscall.Kill(processGroup, syscall.SIGTERM)
	return errors.Wrap(err, "killing process failed")
}

// MustKill try to kill the process and wait for it to die in a timeout period
func (p *process) MustKill() error {
	if err := p.Kill(); err != nil {
		return err
	}
	select {
	case <-p.errch:
		// exited
	case <-time.After(1 * time.Minute):
		log.Warnf("kill process [%d] waiting timeout, try to send a kill -9 signal", p.Pid())
		if err := p.Signal(syscall.SIGKILL); err != nil {
			return errors.Wrap(err, "killing process -9 failed")
		}
		// make sure process exited
		<-p.errch
	}
	return nil
}
