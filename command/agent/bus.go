package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"time"

	mqtt "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
)

const (
	updateTopic = "$sphere/leds/update"
	statusTopic = "$sphere/leds/status"
	stateTopic  = "$sphere/leds/state"
)

/*
 Just manages all the data going into out of this service.
*/
type Bus struct {
	conf   *Config
	agent  *Agent
	client *mqtt.MqttClient
	ticker *time.Ticker
}

type updateRequest struct {
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

	return &Bus{conf: conf, agent: agent}
}

func (b *Bus) listen() {
	log.Printf("[INFO] connecting to the bus")

	opts := mqtt.NewClientOptions().SetBroker(b.conf.LocalUrl).SetClientId("mqtt-bridgeify")

	// shut up
	opts.SetTraceLevel(mqtt.Off)

	b.client = mqtt.NewClient(opts)

	_, err := b.client.Start()
	if err != nil {
		log.Fatalf("error starting connection: %s", err)
	} else {
		fmt.Printf("Connected as %s\n", b.conf.LocalUrl)
	}

	topicFilter, _ := mqtt.NewTopicFilter(updateTopic, 0)
	if _, err := b.client.StartSubscription(b.handleUpdate, topicFilter); err != nil {
		log.Fatalf("error starting subscription: %s", err)
	}

	ev := &statusEvent{Status: "started"}

	b.client.PublishMessage(stateTopic, b.encodeRequest(ev))

	b.setupBackgroundJob()

}

func (b *Bus) handleUpdate(client *mqtt.MqttClient, msg mqtt.Message) {
	log.Printf("[INFO] handleUpdate")
	req := &updateRequest{}
	err := b.decodeRequest(&msg, req)
	if err != nil {
		log.Printf("[ERR] Unable to decode connect request %s", err)
	}
	b.agent.updateLeds(req)

}

func (b *Bus) setupBackgroundJob() {
	b.ticker = time.NewTicker(10 * time.Second)

	for {
		select {
		case <-b.ticker.C:
			// emit the status
			status := b.agent.getStatus()
			log.Printf("[DEBUG] status %+v", status)
			b.client.PublishMessage(statusTopic, b.encodeRequest(status))
		}
	}

}

func (b *Bus) encodeRequest(data interface{}) *mqtt.Message {
	buf := bytes.NewBuffer(nil)
	json.NewEncoder(buf).Encode(data)
	return mqtt.NewMessage(buf.Bytes())
}

func (b *Bus) decodeRequest(msg *mqtt.Message, data interface{}) error {
	return json.NewDecoder(bytes.NewBuffer(msg.Payload())).Decode(data)
}
