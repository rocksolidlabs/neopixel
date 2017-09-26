package main

import (
	"flag"
	"fmt"
	"os/user"
	"runtime"

	"github.com/rocksolidlabs/neopixel"
)

var gpioPin = flag.Int("gpio-pin", 12, "GPIO pin")
var width = flag.Int("width", 1, "LED matrix width")
var height = flag.Int("height", 93, "LED matrix height")
var brightness = flag.Int("brightness", 255, "Brightness (0-255)")
var red = flag.Int("red", 255, "Red")
var green = flag.Int("green", 255, "Gren")
var blue = flag.Int("blue", 255, "Blue")
var duration = flag.Int("duration", 20, "Duration between rings in miliseconds")

func main() {
	user, err := user.Current()
	if err != nil {
		fmt.Println(err)
	}
	if runtime.GOARCH == "arm" && user.Uid != "0" {
		fmt.Println("This program requires root privilege")
		fmt.Println("Please try \"sudo ledtest\"")
	}

	size := *width * *height
	conf := neopixel.DefaultConfig
	conf.Brightness = *brightness
	conf.LEDCount = size
	conf.Brightness = *gpioPin
	conf.StripeType = neopixel.StripGRB

	np, err := neopixel.NewNeoPixel(&conf)
	if err != nil {
		fmt.Println(err)
	}

	err = np.Init()
	if err != nil {
		fmt.Println(err)
	}

	//color := neopixel.Color(*red, *green, *blue)

	red := neopixel.Color(255, 0, 0)
	green := neopixel.Color(0, 255, 0)
	blue := neopixel.Color(0, 0, 255)

	fmt.Println("Theater Chase Animation")
	np.TheaterChase(blue, *duration, 100)

	fmt.Println("Rainbow Animation")
	np.Rainbow(*duration, 100)

	fmt.Println("Rainbow Cycle Animation")
	np.RainbowCycle(*duration, 100)

	fmt.Println("Color Wipe Animation")
	np.ColorWipe(blue, *duration)
	np.ColorWipe(green, *duration)
	np.ColorWipe(blue, *duration)
	np.ColorWipe(red, *duration)

	fmt.Println("Theater Chase Rainbow Animation")
	np.TheaterChaseRainbow(*duration, 100)

	np.Clear()
	np.Render()
	np.Wait()
	np.Fini()
}

func init() {
	flag.Parse()
}
