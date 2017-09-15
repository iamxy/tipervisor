package daemon

import (
	"sync/atomic"
	"syscall"
	"time"

	"github.com/pkg/errors"
)

// Signal defines the types of signal which is sent to daemon to impact process running
type Signal int

// Enum values of the Signal type
const (
	SignalAlrm Signal = iota
	SignalCont
	SignalDown
	SignalHup
	SignalRestart
	SignalInterrupt
	SignalTtin
	SignalKill
	SignalTtou
	SignalStop
	SignalQuit
	SignalUp
	SignalUsr1
	SignalUsr2
	SignalWinch
)

func (s Signal) String() string {
	switch s {
	case SignalAlrm:
		return "ALRM"
	case SignalCont:
		return "CONT"
	case SignalDown:
		return "DOWN"
	case SignalHup:
		return "HUP"
	case SignalRestart:
		return "RESTART"
	case SignalInterrupt:
		return "INTERRUPT"
	case SignalTtin:
		return "TTIN"
	case SignalKill:
		return "KILL"
	case SignalTtou:
		return "TTOU"
	case SignalStop:
		return "STOP"
	case SignalQuit:
		return "QUIT"
	case SignalUp:
		return "UP"
	case SignalUsr1:
		return "USR1"
	case SignalUsr2:
		return "USR2"
	case SignalWinch:
		return "WINCH"
	default:
		return "UNKNOWN"
	}
}

// SignalRequest includes the specific signal value and a response channel to send back result
type SignalRequest struct {
	signal Signal
	respc  chan error
}

// Signal sends a given signal, and waiting for the daemon return
func (d *Daemon) Signal(s Signal) error {
	rc := make(chan error, 1)
	d.sigch <- SignalRequest{
		signal: s,
		respc:  rc,
	}
	select {
	case err := <-rc:
		if err != nil {
			return errors.Wrapf(err, "handling signal [%v] return error", s)
		}
	case <-time.After(30 * time.Second):
		return errors.Errorf("waiting timeout 30s for signal [%v]", s)
	}
	return nil
}

func (d *Daemon) handleSignal(r SignalRequest) {
	var err error

	switch r.signal {
	case SignalAlrm:
		err = d.proc.Signal(syscall.SIGALRM)
	case SignalCont:
		err = d.proc.Signal(syscall.SIGCONT)
	case SignalDown:
		s := d.ProcessState()
		switch s {
		case ProcStatRunning:
			d.changeToState(ProcStatStopping)
			d.lockOnce = 1
			err = d.proc.Kill()
		default:
			err = errors.Errorf("can't stop process from state [%v]", s)
		}
	case SignalHup:
		err = d.proc.Signal(syscall.SIGHUP)
	case SignalRestart:
		s := d.ProcessState()
		switch s {
		case ProcStatRunning:
			d.changeToState(ProcStatRestarting)
			d.lockOnce = 0
			err = d.proc.Kill()
		default:
			err = errors.Errorf("can't restart process from state [%v]", s)
		}
	case SignalInterrupt:
		err = d.proc.Signal(syscall.SIGINT)
	case SignalTtin:
		err = d.proc.Signal(syscall.SIGTTIN)
	case SignalKill:
		s := d.ProcessState()
		switch s {
		case ProcStatRunning, ProcStatStopping, ProcStatRestarting:
			d.changeToState(ProcStatKilling)
			d.lockOnce = 1
			err = d.proc.Signal(syscall.SIGKILL)
		default:
			err = errors.Errorf("can't kill process from state [%v]", s)
		}
	case SignalTtou:
		err = d.proc.Signal(syscall.SIGTTOU)
	case SignalStop:
		err = d.proc.Signal(syscall.SIGSTOP)
	case SignalQuit:
		err = d.proc.Signal(syscall.SIGQUIT)
	case SignalUp:
		s := d.ProcessState()
		switch s {
		case ProcStatStopped, ProcStatKilled, ProcStatExited:
			d.changeToState(ProcStatStarting)
			d.lockOnce = 0
			atomic.StoreUint32(&d.lock, d.lockOnce)
			d.runch <- struct{}{}
		default:
			err = errors.Errorf("can't start process from state [%v]", s)
		}
	case SignalUsr1:
		err = d.proc.Signal(syscall.SIGUSR1)
	case SignalUsr2:
		err = d.proc.Signal(syscall.SIGUSR2)
	case SignalWinch:
		err = d.proc.Signal(syscall.SIGWINCH)
	default:
		err = errors.Errorf("unknown signal")
	}

	r.respc <- err
}
