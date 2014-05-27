package leds

import (
	// "bufio"
	"flag"
	"fmt"
	// "io/ioutil"
	"os"
)

var colors = map[string][]int{
	"black":   {0, 0, 0},
	"red":     {1, 0, 0},
	"green":   {0, 1, 0},
	"blue":    {0, 0, 1},
	"cyan":    {0, 1, 1},
	"magenta": {1, 0, 1},
	"yellow":  {1, 1, 0},
	"white":   {1, 1, 1},
}

var led_types = map[string]int{
	"power":          0,
	"wired_internet": 1,
	"wireless":       2,
	"pairing":        3,
	"radio":          4,
}

var led_positions = [][]int{
	{15, 13, 14},
	{12, 10, 11},
	{9, 1, 8},
	{2, 4, 3},
	{5, 7, 6},
}

var leds = []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

var p = fmt.Println

func main() {
	colorPtr := flag.String("color", "black", "give me a color mofo")
	flag.Parse()
	for i := 0; i < 5; i++ {
		setColor(i, *colorPtr, false)
	}
	setLEDs()
}

func SetStatus(ledtype string, color string, flash bool) {
	if led_types[ledtype] == nil {
		p("bad led type: " + ledtype)
		//#TODO barf nicely?
	}
	if colors[color] == nil {
		p("bad color name: " + color)
		//#TODO barf nicely?
	}

	setColor(led_types[ledtype], colors[color], flash)
}

func _setColor(position int, color []int) {
	var indexes = led_positions[position]
	for i := 0; i < 3; i++ {
		leds[indexes[i]] = color[i]
	}
}

func setColor(position int, color string, flash bool) {
	_setColor(position, colors[color])
	//#TODO FLASHY HERE
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func setLEDs() {

	writetofile("/sys/kernel/debug/omap_mux/lcd_data15", "27")
	writetofile("/sys/kernel/debug/omap_mux/lcd_data14", "27")
	writetofile("/sys/kernel/debug/omap_mux/uart0_ctsn", "27")
	writetofile("/sys/kernel/debug/omap_mux/mii1_col", "27")

	if _, err := os.Stat("/sys/class/gpio/gpio11/direction"); os.IsNotExist(err) {
		writetofile("/sys/class/gpio/export", "11")
	}

	if _, err := os.Stat("/sys/class/gpio/gpio10/direction"); os.IsNotExist(err) {
		writetofile("/sys/class/gpio/export", "10")
	}

	if _, err := os.Stat("/sys/class/gpio/gpio40/direction"); os.IsNotExist(err) {
		writetofile("/sys/class/gpio/export", "40")
	}

	if _, err := os.Stat("/sys/class/gpio/gpio96/direction"); os.IsNotExist(err) {
		writetofile("/sys/class/gpio/export", "96")
	}

	writetofile("/sys/class/gpio/gpio11/direction", "low")
	writetofile("/sys/class/gpio/gpio10/direction", "low")
	writetofile("/sys/class/gpio/gpio40/direction", "low")
	writetofile("/sys/class/gpio/gpio96/direction", "low")

	for i := 0; i < len(leds); i++ {
		writetofile("/sys/class/gpio/gpio40/value", fmt.Sprintf("%d", leds[i]))
		writetofile("/sys/class/gpio/gpio96/value", "1")
		writetofile("/sys/class/gpio/gpio96/value", "0")
	}

	writetofile("/sys/class/gpio/gpio11/value", "1")
	writetofile("/sys/class/gpio/gpio11/value", "0")
}

func writetofile(fn string, val string) error {
	df, err := os.OpenFile(fn, os.O_WRONLY|os.O_SYNC, 0666)
	if err != nil {
		panic(err)
	}

	_, err = fmt.Fprintln(df, val)

	if err != nil {
		panic(err)
	}

	df.Close()
	return nil
}
