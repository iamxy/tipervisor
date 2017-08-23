package daemon

// ProcessState defines the process running state
type ProcessState int

// Enum values of the ProcessState type,
// borrowed from Supervisord
const (
	ProcStatStopped    ProcessState = iota
	ProcStatStarting                = 10
	ProcStatRunning                 = 20
	ProcStatRestarting              = 30
	ProcStatStopping                = 40
	ProcStatKilling                 = 50
	ProcStatKilled                  = 60
	ProcStatExited                  = 100
	ProcStatFatal                   = 200
	ProcStatUnknown                 = 1000
)

func (s ProcessState) String() string {
	var ret string

	switch s {
	case ProcStatStopped:
		ret = "STOPPED"
	case ProcStatStarting:
		ret = "STARTING"
	case ProcStatRunning:
		ret = "RUNNING"
	case ProcStatRestarting:
		ret = "RESTARTING"
	case ProcStatStopping:
		ret = "STOPPING"
	case ProcStatKilling:
		ret = "KILLING"
	case ProcStatExited:
		ret = "EXITED"
	case ProcStatFatal:
		ret = "FATAL"
	case ProcStatUnknown:
		ret = "UNKNOWN"
	default:
		ret = "UNKNOWN"
	}

	return ret
}

func (d *Daemon) changeToState(s ProcessState) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.state = s
	// todo: send an event here
}

// ProcessState returns the current process state
func (d *Daemon) ProcessState() ProcessState {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.state
}
