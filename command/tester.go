package command

import (
	"flag"
	"fmt"

	"github.com/mitchellh/cli"
	"github.com/ninjablocks/sphere-leds/leds"
)

type TesterCommand struct {
	Ui         cli.Ui
	ShutdownCh <-chan struct{}
	args       []string
	ledArray   *leds.LedArray
}

func (c *TesterCommand) Run(args []string) int {

	c.Ui = &cli.PrefixedUi{
		OutputPrefix: "==> ",
		InfoPrefix:   "    ",
		ErrorPrefix:  "==> ",
		Ui:           c.Ui,
	}

	var color string
	var ledName string

	cmdFlags := flag.NewFlagSet("agent", flag.ContinueOnError)
	cmdFlags.Usage = func() { c.Ui.Output(c.Help()) }
	cmdFlags.StringVar(&color, "color", "", "set the led to this color")
	cmdFlags.StringVar(&ledName, "ledname", "", "name of the led to set")

	if err := cmdFlags.Parse(args); err != nil {
		return 1
	}

	if color == "" {
		c.Ui.Error("color flag required")
		return 1
	}

	if leds.Colors[color] == nil {
		c.Ui.Error("Invalid color supplied")
		return 1
	}

	if !leds.ValidLedName(ledName) {
		c.Ui.Error("Invalid led name supplied")
		return 1
	}

	c.ledArray = leds.CreateLedArray()
	c.Ui.Output("Sphere LEDs tester")
	c.Ui.Info(fmt.Sprintf("Setting power %s to: %s", ledName, color))
	c.ledArray.SetColor(leds.LedNameIndex(ledName), color, true)
	c.ledArray.SetLEDs()

	return 0
}

func (c *TesterCommand) Synopsis() string {
	return "Runs a Sphere LEDs tester"
}

func (c *TesterCommand) Help() string {
	helpText := `
Usage: sphere-leds tester [options]

  Starts the Sphere LEDs tester.

Options:

  -color=red                          Set the led to this color.
  -ledname=power                      Name of the led to set.
`
	return helpText
}
