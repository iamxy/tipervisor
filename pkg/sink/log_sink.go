package sink

import "os"

// LogSink controls to start/stop log listening
type LogSink interface {
	Start(pout, perr *os.File)
	Stop()
}

// LogSinkFactory is the interface that wraps the create LogSink method
type LogSinkFactory interface {
	NewLogSink() LogSink
}
