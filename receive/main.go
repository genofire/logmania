package receive

import (
	"github.com/genofire/logmania/lib"
	log "github.com/sirupsen/logrus"
)

var Register = make(map[string]ReceiverInit)

type Receiver interface {
	Listen()
	Close()
}

type ReceiverInit func(*lib.ReceiveConfig, chan *log.Entry) Receiver

func AddReceiver(name string, n ReceiverInit) {
	Register[name] = n
}
