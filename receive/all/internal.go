package all

import (
	log "github.com/sirupsen/logrus"

	"dev.sum7.eu/genofire/logmania/lib"
	"dev.sum7.eu/genofire/logmania/receive"
)

type Receiver struct {
	receive.Receiver
	list []receive.Receiver
}

func Init(config *lib.ReceiveConfig, exportChannel chan *log.Entry) receive.Receiver {
	var list []receive.Receiver
	for _, init := range receive.Register {
		receiver := init(config, exportChannel)

		if receiver == nil {
			continue
		}
		list = append(list, receiver)
	}
	return &Receiver{
		list: list,
	}
}

func (r *Receiver) Listen() {
	for _, item := range r.list {
		go item.Listen()
	}
}

func (r *Receiver) Close() {
	for _, item := range r.list {
		item.Close()
	}
}
