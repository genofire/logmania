package log

// interface of a logger
type Logger interface {
	Hook(*Entry)
	Close()
}

var loggers = make([]Logger, 0)

// bind logger to handle saving/output of a Log entry
func AddLogger(logger Logger) {
	loggers = append(loggers, logger)
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
