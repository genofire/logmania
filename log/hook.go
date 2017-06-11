package log

type Hook func(e *Entry)

var hooks = make([]Hook, 0)

func AddHook(hook Hook) {
	hooks = append(hooks, hook)
}

func save(e *Entry) {
	for _, hook := range hooks {
		hook(e)
	}
}
