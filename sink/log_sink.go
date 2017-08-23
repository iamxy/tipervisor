package sink

import "os"

// LogSink controls to start/stop log listening
type LogSink interface {
	Start(pout, perr *os.File)
	Stop()
}
