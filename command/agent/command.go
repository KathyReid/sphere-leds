package agent

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/hashicorp/logutils"
	"github.com/mitchellh/cli"
	nlog "github.com/ninjasphere/go-ninja/logger"
)

var logger = nlog.GetLogger("sphere-leds")

/*
 Holds all the context for a running agent
*/
type Command struct {
	Ui         cli.Ui
	ShutdownCh <-chan struct{}
	args       []string
	logFilter  *logutils.LevelFilter
	agent      *Agent
	bus        *Bus
}

type Config struct {
	Token       string
	CloudUrl    string
	LocalUrl    string
	Debug       bool
	StatusTimer int
}

func (c *Config) IsDebug() bool {
	return c.Debug
}

func (c *Command) readConfig() *Config {

	var cmdConfig Config

	cmdFlags := flag.NewFlagSet("agent", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()) }
	cmdFlags.StringVar(&cmdConfig.LocalUrl, "localurl", "tcp://localhost:1883", "cloud url to connect to")

	if err := cmdFlags.Parse(c.args); err != nil {
		return nil
	}

	return &cmdConfig
}

func (c *Command) handleSignals(config *Config) int {
	signalCh := make(chan os.Signal, 4)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	var sig os.Signal
	select {
	case s := <-signalCh:
		sig = s
	case <-c.ShutdownCh:
		sig = os.Interrupt
	}
	c.Ui.Output(fmt.Sprintf("Caught signal: %v", sig))

	return 0
}

func (c *Command) Run(args []string) int {
	c.Ui = &cli.PrefixedUi{
		OutputPrefix: "==> ",
		InfoPrefix:   "    ",
		ErrorPrefix:  "==> ",
		Ui:           c.Ui,
	}

	c.args = args
	config := c.readConfig()
	if config == nil {
		return 1
	}
	c.args = args

	c.Ui.Output("Sphere LEDs agent running!")
	c.Ui.Info("Getting on the bus: " + config.LocalUrl)

	c.agent = createAgent(config)

	if err := c.agent.start(); err != nil {
		c.Ui.Error(fmt.Sprintf("error starting agent %s", err))
	}

	c.bus = createBus(config, c.agent)

	c.bus.listen()

	return c.handleSignals(config)
}

func (c *Command) Synopsis() string {
	return "Runs a Sphere LEDs agent"
}

func (c *Command) Help() string {
	helpText := `
Usage: sphere-leds agent [options]

  Starts the Sphere LEDs agent and runs until an interrupt is received.

Options:

  -localurl=tcp://localhost:1883      URL for the local broker.
`
	return helpText
}
