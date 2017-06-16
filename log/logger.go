package log

// interface of a logger
type Logger interface {
	Hook(*Entry)
	Close()
}

var loggers = make(map[string]Logger)

// bind logger to handle saving/output of a Log entry
func AddLogger(name string, logger Logger) {
	if logger != nil {
		loggers[name] = logger
	}
}
func RemoveLogger(name string) {
	loggers[name].Close()
	delete(loggers, name)
}

func save(e *Entry) {
	for _, logger := range loggers {
		logger.Hook(e)
	}
	if e.Level == PanicLevel {
		for _, logger := range loggers {
			logger.Close()
		}
		panic("panic see last log in logmania")
	}
}
