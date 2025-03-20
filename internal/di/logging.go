package di

type Logger interface {
	Debug(a ...any)
	Info(a ...any)
	Success(a ...any)
	Warning(a ...any)
}

type nilLogger struct{}

func (log nilLogger) Debug(a ...any)   {}
func (log nilLogger) Info(a ...any)    {}
func (log nilLogger) Success(a ...any) {}
func (log nilLogger) Warning(a ...any) {}
