package leds

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

var Colors = map[string][]int{
	"black":   {0, 0, 0},
	"red":     {1, 0, 0},
	"green":   {0, 1, 0},
	"blue":    {0, 0, 1},
	"cyan":    {0, 1, 1},
	"magenta": {1, 0, 1},
	"yellow":  {1, 1, 0},
	"white":   {1, 1, 1},
}

var LedNames = []string{
	"power",
	"wired_internet",
	"wireless",
	"pairing",
	"radio",
}

var LedPositions = [][]int{
	{15, 13, 14},
	{12, 10, 11},
	{9, 1, 8},
	{2, 4, 3},
	{5, 7, 6},
}

// holds the state for an array of leds on our board.
type LedArray struct {
	Leds      []int
	LedStates []LedState
	ticker    *time.Ticker
	lock      *sync.Mutex
}

type LedState struct {
	Flash bool
	Color string
	On    bool
}

func CreateLedArray() *LedArray {
	ledArr := &LedArray{
		Leds:      []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		LedStates: make([]LedState, 5),
		lock:      &sync.Mutex{},
	}
	initLEDs()
	go ledArr.setupBackgroundJob()

	return ledArr
}

func (l *LedArray) setupBackgroundJob() {
	l.ticker = time.NewTicker(1 * time.Second)

	for {
		select {
		case <-l.ticker.C:

			l.lock.Lock()

			//log.Println("[DEBUG] flash")
			for n := range l.LedStates {
				if l.LedStates[n].Flash {

					if l.LedStates[n].On {
						l.setColorInt(n, Colors["black"])
						l.LedStates[n].On = false
					} else {
						l.setColorInt(n, Colors[l.LedStates[n].Color])
						l.LedStates[n].On = true
					}
				}
			}
			// we should use a cached copy of the array then do a diff and conditionally
			// set leds.
			l.SetLEDs()

			l.lock.Unlock()

		}
	}
}

func (l *LedArray) setColorInt(position int, color []int) {
	var indexes = LedPositions[position]
	for i := 0; i < 3; i++ {
		l.Leds[indexes[i]] = color[i]
	}
}

func (l *LedArray) SetPwmBrightness(brightness int) {
	writetofile("/sys/class/backlight/pwm-backlight/brightness", fmt.Sprintf("%d", brightness))
}

func (l *LedArray) SetColor(position int, color string, flash bool) {
	defer l.lock.Unlock()
	l.lock.Lock()
	// update the state
	l.LedStates[position].Flash = flash
	l.LedStates[position].Color = color
	l.LedStates[position].On = true
	l.setColorInt(position, Colors[color])
	// apply it
	l.SetLEDs()
}

func (l *LedArray) Reset() {
	defer l.lock.Unlock()
	l.lock.Lock()
	for pos := range LedNames {
		// update the states
		l.LedStates[pos].Flash = false
		l.LedStates[pos].Color = "black"
		l.LedStates[pos].On = true
		l.SetColor(pos, "black", false)
	}
	// apply it
	l.SetLEDs()
}

func ValidBrightness(brightness int) bool {
	return brightness >= 0 && brightness <= 100
}

func ValidColor(color string) bool {
	return Colors[color] != nil
}

func ValidLedName(name string) bool {
	for n := range LedNames {
		if LedNames[n] == name {
			return true
		}
	}
	return false
}

func LedNameIndex(name string) int {
	for n := range LedNames {
		if LedNames[n] == name {
			return n
		}
	}
	panic("LedName didn't exist.")
}

func initLEDs() {
	log.Printf("Initializing LEDs")
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

}

func (l *LedArray) SetLEDs() {
	//	log.Printf("[DEBUG] Updating leds: %v", l.Leds)
	//	log.Printf("[DEBUG] Updating flashstate: %v", l.LedStates)

	for i := range l.Leds {
		writetofile("/sys/class/gpio/gpio40/value", fmt.Sprintf("%d", l.Leds[i]))
		writetofile("/sys/class/gpio/gpio96/value", "1")
		writetofile("/sys/class/gpio/gpio96/value", "0")
	}

	writetofile("/sys/class/gpio/gpio11/value", "1")
	writetofile("/sys/class/gpio/gpio11/value", "0")
}

func writetofile(fn string, val string) error {

	df, err := os.OpenFile(fn, os.O_WRONLY|os.O_SYNC, 0666)

	if err != nil {
		log.Printf("[ERROR] failed to open file %s - %s", fn, err)
		return err
	}

	defer df.Close()

	if _, err = fmt.Fprintln(df, val); err != nil {
		log.Printf("[ERROR] failed write to %s - %s", fn, err)
		return err
	}

	return nil
}
