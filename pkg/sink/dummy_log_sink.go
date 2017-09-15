package sink

import (
	"os"
)

// DummyLogSink implements dummy log sink for testing use
type DummyLogSink struct {
}

// Start gets log sink to work
func (s *DummyLogSink) Start(pout, perr *os.File) {
	// do nothing
}

// Stop terminates log sink
func (s *DummyLogSink) Stop() {
	// do nothing
}

// DummyLogSinkFactory implements dummy log sink factory
type DummyLogSinkFactory struct {
}

// NewLogSink creates a new dummy log sink
func (f *DummyLogSinkFactory) NewLogSink() LogSink {
	return &DummyLogSink{}
}

// NewDummyLogSinkFactory returns a dummy log sink factory
func NewDummyLogSinkFactory() LogSinkFactory {
	return &DummyLogSinkFactory{}
}
