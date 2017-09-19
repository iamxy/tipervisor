package daemon

import (
	"context"
	"fmt"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/pingcap/tipervisor/pkg/sink"
	"github.com/stretchr/testify/assert"
)

func NewDaemonConfig(name string) *Config {
	wd, _ := os.Getwd()
	return &Config{
		Name:      name,
		Cmd:       "sleep",
		Args:      []string{"3600"},
		Cwd:       wd,
		StatusDir: os.TempDir(),
	}
}

func TestSuperviseRunning(t *testing.T) {
	cfg := NewDaemonConfig("test_supervise_running")
	lsf := sink.NewDummyLogSinkFactory()
	d, err := New(cfg, lsf)
	assert.NoError(t, err)
	// start to supervise
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	d.Supervise(ctx)
	assert.Equal(t, ProcStatRunning, d.ProcessState())
	pid := d.proc.Pid()
	assert.Equal(t, pid, d.GetRunningStat().Pid)
	assert.True(t, isRunning(pid))
	// force kill the process
	assert.NoError(t, syscall.Kill(pid, syscall.SIGKILL))
	// check
	time.Sleep(5 * time.Second)
	assert.Equal(t, ProcStatRunning, d.ProcessState())
	stat := d.GetRunningStat()
	assert.NotZero(t, stat.LastUpTime)
	assert.Equal(t, ProcStatExited, stat.LastTerminateState)
	assert.Equal(t, uint32(2), stat.RunCount)
	assert.Equal(t, uint32(1), stat.ExitedCount)
}

func TestManualKill(t *testing.T) {
	cfg := NewDaemonConfig("test_manual_kill")
	lsf := sink.NewDummyLogSinkFactory()
	d, err := New(cfg, lsf)
	assert.NoError(t, err)
	// start to supervise
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	d.Supervise(ctx)
	assert.Equal(t, ProcStatRunning, d.ProcessState())
	pid := d.proc.Pid()
	assert.Equal(t, pid, d.GetRunningStat().Pid)
	assert.True(t, isRunning(pid))
	// send KILL signal
	assert.NoError(t, d.Signal(SignalKill))
	// check
	time.Sleep(5 * time.Second)
	assert.Equal(t, ProcStatKilled, d.ProcessState())
	stat := d.GetRunningStat()
	assert.NotZero(t, stat.LastUpTime)
	assert.Equal(t, ProcStatKilled, stat.LastTerminateState)
	assert.Equal(t, uint32(1), stat.RunCount)
	assert.Equal(t, uint32(1), stat.KilledCount)
	assert.False(t, isRunning(pid))
}

func TestManualStopAndStart(t *testing.T) {
	cfg := NewDaemonConfig("test_manual_stop_and_start")
	lsf := sink.NewDummyLogSinkFactory()
	d, err := New(cfg, lsf)
	assert.NoError(t, err)
	// start to supervise
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	d.Supervise(ctx)
	assert.Equal(t, ProcStatRunning, d.ProcessState())
	pid := d.proc.Pid()
	assert.Equal(t, pid, d.GetRunningStat().Pid)
	assert.True(t, isRunning(pid))
	// send DOWN signal
	assert.NoError(t, d.Signal(SignalDown))
	// check
	time.Sleep(5 * time.Second)
	assert.Equal(t, ProcStatStopped, d.ProcessState())
	stat := d.GetRunningStat()
	assert.NotZero(t, stat.LastUpTime)
	assert.Equal(t, ProcStatStopped, stat.LastTerminateState)
	assert.Equal(t, uint32(1), stat.RunCount)
	assert.Equal(t, uint32(1), stat.StoppedCount)
	assert.False(t, isRunning(pid))
	// send UP signal
	assert.NoError(t, d.Signal(SignalUp))
	// check
	time.Sleep(1 * time.Second)
	assert.Equal(t, ProcStatRunning, d.ProcessState())
	pid = d.proc.Pid()
	stat = d.GetRunningStat()
	assert.Equal(t, uint32(2), stat.RunCount)
	assert.Equal(t, uint32(1), stat.StoppedCount)
	assert.True(t, isRunning(pid))
}

func TestManualRestart(t *testing.T) {
	cfg := NewDaemonConfig("test_manual_restart")
	lsf := sink.NewDummyLogSinkFactory()
	d, err := New(cfg, lsf)
	assert.NoError(t, err)
	// start to supervise
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	d.Supervise(ctx)
	assert.Equal(t, ProcStatRunning, d.ProcessState())
	pid := d.proc.Pid()
	assert.Equal(t, pid, d.GetRunningStat().Pid)
	assert.True(t, isRunning(pid))
	// send RESTART signal
	assert.NoError(t, d.Signal(SignalRestart))
	// check
	time.Sleep(5 * time.Second)
	assert.Equal(t, ProcStatRunning, d.ProcessState())
	stat := d.GetRunningStat()
	assert.NotZero(t, stat.LastUpTime)
	assert.Equal(t, ProcStatStopped, stat.LastTerminateState)
	assert.Equal(t, uint32(2), stat.RunCount)
	assert.Equal(t, uint32(1), stat.StoppedCount)
	assert.False(t, isRunning(pid))
}

func TestKillAfterRestart(t *testing.T) {
	cfg := NewDaemonConfig("test_kill_after_restart")
	lsf := sink.NewDummyLogSinkFactory()
	d, err := New(cfg, lsf)
	assert.NoError(t, err)
	// start to supervise
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		for {
			fmt.Printf("################################################## %v\n", d.ProcessState())
			time.Sleep(100 * time.Millisecond)
		}
	}()
	d.Supervise(ctx)
	assert.Equal(t, ProcStatRunning, d.ProcessState())
	pid := d.proc.Pid()
	assert.Equal(t, pid, d.GetRunningStat().Pid)
	assert.True(t, isRunning(pid))
	// send RESTART signal
	assert.NoError(t, d.Signal(SignalRestart))
	time.Sleep(1 * time.Second)
	assert.Equal(t, ProcStatRestarting, d.ProcessState())
	// send KILL signal immediately
	// expects to return a "no such process" error
	_ = d.Signal(SignalKill)
	time.Sleep(1 * time.Second)
	assert.Equal(t, ProcStatKilling, d.ProcessState())
	// check
	time.Sleep(3 * time.Second)
	assert.Equal(t, ProcStatKilled, d.ProcessState())
	stat := d.GetRunningStat()
	assert.NotZero(t, stat.LastUpTime)
	assert.Equal(t, ProcStatKilled, stat.LastTerminateState)
	assert.Equal(t, uint32(1), stat.RunCount)
	assert.Equal(t, uint32(1), stat.KilledCount)
	assert.False(t, isRunning(pid))
}
