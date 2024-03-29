package di

type Logger interface {
	Info(a ...interface{})
	Success(a ...interface{})
	Warning(a ...interface{})
}

type nilLogger struct{}

func (log nilLogger) Info(a ...interface{})    {}
func (log nilLogger) Success(a ...interface{}) {}
func (log nilLogger) Warning(a ...interface{}) {}
