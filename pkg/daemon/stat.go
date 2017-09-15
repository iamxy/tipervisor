package daemon

import (
	"sync"
	"time"
)

// RunStat keeps the statistics of the running process
type RunStat struct {
	sync.RWMutex
	LastStartTime      time.Time
	LastEndTime        time.Time
	LastUpTime         time.Duration
	LastUserTime       time.Duration
	LastSysTime        time.Duration
	LastTerminateState ProcessState
	LastExitErr        error
	StartTime          time.Time
	RunCount           uint32
	StoppedCount       uint32
	ExitedCount        uint32
	KilledCount        uint32
	Pid                int
}

// GetRunningStat return a RunStat Object containing the statistics of the daemon runtime
func (d *Daemon) GetRunningStat() *RunStat {
	d.runStat.Lock()
	defer d.runStat.Unlock()
	return &RunStat{
		LastStartTime:      d.runStat.LastStartTime,
		LastEndTime:        d.runStat.LastEndTime,
		LastUpTime:         d.runStat.LastUpTime,
		LastUserTime:       d.runStat.LastUserTime,
		LastSysTime:        d.runStat.LastSysTime,
		LastTerminateState: d.runStat.LastTerminateState,
		LastExitErr:        d.runStat.LastExitErr,
		StartTime:          d.runStat.StartTime,
		RunCount:           d.runStat.RunCount,
		StoppedCount:       d.runStat.StoppedCount,
		ExitedCount:        d.runStat.ExitedCount,
		KilledCount:        d.runStat.KilledCount,
		Pid:                d.runStat.Pid,
	}
}
