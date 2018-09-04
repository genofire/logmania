package syslog

import (
	"net"

	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"

	"dev.sum7.eu/genofire/logmania/input"
)

const inputType = "syslog"

var logger = log.WithField("input", inputType)

type Input struct {
	input.Input
	exportChannel chan *log.Entry
	serverSocket  *net.UDPConn
}

type InputConfig struct {
	Type    string `mapstructure:"type"`
	Address string `mapstructure:"address"`
}

func Init(configInterface interface{}, exportChannel chan *log.Entry) input.Input {
	var config InputConfig
	if err := mapstructure.Decode(configInterface, &config); err != nil {
		logger.Warnf("not able to decode data: %s", err)
		return nil
	}
	addr, err := net.ResolveUDPAddr(config.Type, config.Address)
	ln, err := net.ListenUDP(config.Type, addr)

	if err != nil {
		logger.Error("init ", err)
		return nil
	}
	input := &Input{
		serverSocket:  ln,
		exportChannel: exportChannel,
	}

	logger.Info("init")

	return input
}

const maxDataGramSize = 8192

func (in *Input) Listen() {
	logger.Info("listen")
	for {
		buf := make([]byte, maxDataGramSize)
		n, src, err := in.serverSocket.ReadFromUDP(buf)
		if err != nil {
			logger.Warn("failed to accept connection", err)
			continue
		}

		raw := make([]byte, n)
		copy(raw, buf)
		entry := toLogEntry(raw, src.IP.String())
		if entry != nil {
			in.exportChannel <- entry
		}
	}
}

func (in *Input) Close() {
	in.serverSocket.Close()
}

func init() {
	input.Add(inputType, Init)
}
