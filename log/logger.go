package log

type loggerFunc func(*Entry)

var Save loggerFunc

func init() {
	Save = func(*Entry) {}
}
