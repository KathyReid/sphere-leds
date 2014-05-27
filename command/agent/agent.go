package agent

import "runtime"

//
// Pulls together the bridge, a cached state configuration and the bus.
//
type Agent struct {
	conf     *Config
	memstats *runtime.MemStats
	eventCh  chan statusEvent
	Leds     *SphereLeds
}

func createAgent(conf *Config) *Agent {
	return &Agent{conf: conf, bridge: createBridge(conf), memstats: &runtime.MemStats{}}
}

// TODO load the existing configuration on startup and start the bridge if needed
func (a *Agent) start() error {

	return nil
}

// stop all the things.
func (a *Agent) stop() error {

	return nil
}

func (a *Agent) updateLeds(connect *connectRequest) {
	a.bridge.start(connect.Url, connect.Token)
}

// save the state of the bridge then disconnect it
func (a *Agent) stopBridge(disconnect *disconnectRequest) {
	a.bridge.stop()
}

func (a *Agent) getStatus() statsEvent {

	runtime.ReadMemStats(a.memstats)

	return statsEvent{
		Alloc:      a.memstats.Alloc,
		HeapAlloc:  a.memstats.HeapAlloc,
		TotalAlloc: a.memstats.TotalAlloc,
	}
}
