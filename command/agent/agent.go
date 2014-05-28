package agent

import (
	"errors"
	"log"
	"runtime"
	"strings"

	"github.com/ninjablocks/sphere-leds/leds"
)

//
// Pulls together the bridge, a cached state configuration and the bus.
//
type Agent struct {
	conf     *Config
	memstats *runtime.MemStats
	eventCh  chan statusEvent
	leds     *leds.LedArray
}

func createAgent(conf *Config) *Agent {
	return &Agent{conf: conf, memstats: &runtime.MemStats{}, leds: leds.CreateLedArray()}
}

// TODO load the existing configuration on startup and start the bridge if needed
func (a *Agent) start() error {

	return nil
}

// stop all the things.
func (a *Agent) stop() error {

	return nil
}

func (a *Agent) updateLeds(update *updateRequest) {

	log.Printf("update leds : %v", update)

	err, action := getLastToken(update.Topic)

	if err != nil {
		log.Printf("[ERROR] bad update request - %s", err)
	}

	if action == "reset" {
		a.leds.Reset()
		if leds.ValidBrightness(update.Brightness) {
			a.leds.SetPwmBrightness(update.Brightness)
		} else {
			log.Printf("[WARN] bad brightness %d", update.Brightness)
		}
		return
	}

	if leds.ValidLedName(action) && leds.ValidColor(update.Color) {
		a.leds.SetColor(leds.LedNameIndex(action), update.Color, update.Flash)
		a.leds.SetLEDs()
	} else {
		log.Printf("[WARN] bad SetColor params - %s %s %s", action, update.Color, update.Flash)
	}

}

func (a *Agent) getStatus() statsEvent {

	runtime.ReadMemStats(a.memstats)

	return statsEvent{
		Alloc:      a.memstats.Alloc,
		HeapAlloc:  a.memstats.HeapAlloc,
		TotalAlloc: a.memstats.TotalAlloc,
	}
}

func getLastToken(topic string) (error, string) {
	toks := strings.Split(topic, "/")
	if len(toks) > 1 {
		return nil, strings.Join(toks[len(toks)-1:], "")
	} else {
		return errors.New("bad topic"), ""
	}
}
