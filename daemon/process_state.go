package daemon

// ProcessState defines the process running state
type ProcessState int

// Enum values of the ProcessState type,
// borrowed from Supervisord
const (
	ProcStatStopped ProcessState = iota
	ProcStatStarting
	ProcStatRunning
	ProcStatRestarting
	ProcStatStopping
	ProcStatKilling
	ProcStatTerminating
	ProcStatExited
	ProcStatKilled
	ProcStatFatal
	ProcStatUnknown
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
	case ProcStatTerminating:
		ret = "TERMINATING"
	case ProcStatExited:
		ret = "EXITED"
	case ProcStatKilled:
		ret = "KILLED"
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
	// todo: emit an event here
}

// ProcessState returns the current process state
func (d *Daemon) ProcessState() ProcessState {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.state
}
