package agent

import (
	"bytes"
	"encoding/json"
	"github.com/ninjasphere/go-ninja/bus"
	"strings"
	"time"
)

const (
	updateTopic = "$hardware/status/+"
	statusTopic = "$sphere/leds/status"
)

/*
 Just manages all the data going into out of this service.
*/
type Bus struct {
	conf   *Config
	agent  *Agent
	bus    bus.Bus
	ticker *time.Ticker
}

type updateRequest struct {
	Topic      string
	Brightness int    `json:"brightness"`
	On         bool   `json:"on"`
	Color      string `json:"color"`
	Flash      bool   `json:"flash"`
}

type statusEvent struct {
	Status string `json:"status"`
}

type statsEvent struct {

	// memory related information
	Alloc      uint64 `json:"alloc"`
	HeapAlloc  uint64 `json:"heapAlloc"`
	TotalAlloc uint64 `json:"totalAlloc"`
}

func createBus(conf *Config, agent *Agent) *Bus {
	var addr string
	if strings.HasPrefix(conf.LocalUrl, "tcp://") {
		addr = conf.LocalUrl[len("tcp://"):]
	} else {
		addr = conf.LocalUrl
	}
	logger.Infof("connecting to the bus %s...", addr)
	theBus := bus.MustConnect(addr, "sphere-leds")
	logger.Infof("connected to the bus.")
	return &Bus{conf: conf, bus: theBus, agent: agent}
}

func (b *Bus) listen() {
	logger.Infof("subscribing to topic %s", updateTopic)
	_, err := b.bus.Subscribe(updateTopic, b.handleUpdate)
	if err != nil {
		panic(err)
	}

	b.setupBackgroundJob()

}

func (b *Bus) handleUpdate(topic string, payload []byte) {
	logger.Debugf("handleUpdate %s", string(payload))
	req := &updateRequest{}
	err := b.decodeRequest(payload, req)
	if err != nil {
		logger.Errorf("Unable to decode connect request %s", err)
	}
	req.Topic = topic
	b.agent.updateLeds(req)

}

func (b *Bus) setupBackgroundJob() {
	b.ticker = time.NewTicker(10 * time.Second)

	for {
		select {
		case <-b.ticker.C:
			// emit the status
			status := b.agent.getStatus()
			logger.Debugf("[DEBUG] status %+v", status)
			b.bus.Publish(statusTopic, b.encodeRequest(status))
		}
	}

}

func (b *Bus) encodeRequest(data interface{}) []byte {
	buf := bytes.NewBuffer(nil)
	json.NewEncoder(buf).Encode(data)
	return buf.Bytes()
}

func (b *Bus) decodeRequest(payloadBytes []byte, data interface{}) error {
	payload := string(payloadBytes)
	if strings.HasPrefix(payload, "[") && strings.HasSuffix(payload, "]") {
		// go rpc can't handle enclosing parameters
		payload = payload[1 : len(payload)-1]
	}
	return json.NewDecoder(bytes.NewBuffer([]byte(payload))).Decode(data)
}
