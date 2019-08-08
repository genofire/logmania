package all

import (
	"github.com/bdlm/log"

	"dev.sum7.eu/sum7/logmania/input"
)

type Input struct {
	input.Input
	list []input.Input
}

func Init(configInterface interface{}, exportChannel chan *log.Entry) input.Input {
	config := configInterface.(map[string]interface{})

	var list []input.Input
	for inputType, init := range input.Register {
		configForItem := config[inputType]
		if configForItem == nil {
			log.Warnf("the input type '%s' has no configuration", inputType)
			continue
		}
		in := init(configForItem, exportChannel)

		if in == nil {
			continue
		}
		list = append(list, in)
	}
	return &Input{
		list: list,
	}
}

func (in *Input) Listen() {
	for _, item := range in.list {
		go item.Listen()
	}
}

func (in *Input) Close() {
	for _, item := range in.list {
		item.Close()
	}
}
