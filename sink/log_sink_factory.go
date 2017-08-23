package sink

type LogSinkFactory interface {
	NewLogSink() LogSink
}
