package input

import (
	"github.com/bdlm/log"
)

var Register = make(map[string]Init)

type Input interface {
	Listen()
	Close()
}

type Init func(interface{}, chan *log.Entry) Input

func Add(name string, init Init) {
	Register[name] = init
}
