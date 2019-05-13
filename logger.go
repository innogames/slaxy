package slaxy

// Logger defines simple methods for info and debugging messages
// that could be of interest for the outside
type Logger interface {
	Debug(string)
	Debugf(string, ...interface{})
	Info(string)
	Infof(string, ...interface{})
	Warn(string)
	Warnf(string, ...interface{})
	Error(string)
	Errorf(string, ...interface{})
}

// nullLogger is a logger that does nothing
type nullLogger struct {
}

// NewNullLogger returns a new instance of a logger that does nothing
func NewNullLogger() Logger {
	return nullLogger{}
}

// Debug logs debug messages
func (l nullLogger) Debug(string) {
}

// Debugf logs debug messages
func (l nullLogger) Debugf(string, ...interface{}) {
}

// Info logs info messages
func (l nullLogger) Info(string) {
}

// Infof logs info messages
func (l nullLogger) Infof(string, ...interface{}) {
}

// Warn logs info messages
func (l nullLogger) Warn(string) {
}

// Warnf logs info messages
func (l nullLogger) Warnf(string, ...interface{}) {
}

// Error logs info messages
func (l nullLogger) Error(string) {
}

// Errorf logs info messages
func (l nullLogger) Errorf(string, ...interface{}) {
}
