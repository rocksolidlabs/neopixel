// Interface to ws2811 chip (neopixel driver). Make sure that you have
// ws2811.h and pwm.h in a GCC include path (e.g. /usr/local/include) and
// libws2811.a in a GCC library path (e.g. /usr/local/lib).
// See https://github.com/jgarff/rpi_ws281x for instructions

package neopixel

// #cgo CFLAGS: -std=c99
// #cgo LDFLAGS: -lws2811
// #include "neopixel.go.h"
import "C"

import (
	"errors"
	"fmt"
	"unsafe"
)

const (
	// DefaultDmaNum is the default DMA number. Usually, this is 5 ob the Raspberry Pi
	DefaultDmaNum = 5
	// TargetFreq is the target frequency. It is usually 800kHz (800000), and an go as low as 400000
	TargetFreq = 800000
	StripRGB   = 0x100800 // StripRGB is the RGB Mode
	StripRBG   = 0x100008 // StripRBG is the RBG Mode
	StripGRB   = 0x081000 // StripGRB is the GRB Mode
	StripGBR   = 0x080010 // StripGBR is the GBR Mode
	StripBRG   = 0x001008 // StripBRG is the BRG Mode
	StripBGR   = 0x000810 // StripBGR is the BGR Mode
)

// Config is the arguments for creating an instance of NeoPixel.
type Config struct {
	Frequency  int
	DmaNum     int
	GpioPin    int
	LedCount   int
	Brightness int
	StripeType int
	Invert     bool
}

// neopixel represent the ws2811 device
type NeoPixel struct {
	Dev *C.ws2811_t
	Config *Config
}

// DefaultConfig defines sensible default options for MakeWS2811
var DefaultConfig = Config{
	Frequency:  800000,
	DmaNum:     5,
	GPIOPin:    17,
	LEDCount:   93,
	Brightness: 255,
	StripeType: StripGRB,
	Invert:     false,
	Ring:		true,
}

// Rings in Mokungit 93 LED WS2812 5050
// http://mokungit.com/product-item/ws2812b-pixel-rgb-ring
var rings = [][]int{
	{92},
	{91,90,89,88,87,86,85,84},
	{83,82,81,80,79,78,77,76,75,74,73,72},
	{71,70,69,68,67,66,65,64,63,62,61,60,59,58,57,56},
	{55,54,53,52,51,50,49,48,47,46,45,44,43,42,41,40,39,38,37,36,35,34,33,32},
	{31,30,29,28,27,26,25,24,23,22,21,20,19,18,17,16,15,14,13,12,11,10,9,8,7,6,5,4,3,2,1,0},
}

var gamma8 = []uint32{
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1, 1, 2, 2, 2, 2, 2, 2, 2,
	2, 3, 3, 3, 3, 3, 3, 3, 4, 4, 4, 4, 4, 5, 5, 5,
	5, 6, 6, 6, 6, 7, 7, 7, 7, 8, 8, 8, 9, 9, 9, 10,
	10, 10, 11, 11, 11, 12, 12, 13, 13, 13, 14, 14, 15, 15, 16, 16,
	17, 17, 18, 18, 19, 19, 20, 20, 21, 21, 22, 22, 23, 24, 24, 25,
	25, 26, 27, 27, 28, 29, 29, 30, 31, 32, 32, 33, 34, 35, 35, 36,
	37, 38, 39, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 50,
	51, 52, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 66, 67, 68,
	69, 70, 72, 73, 74, 75, 77, 78, 79, 81, 82, 83, 85, 86, 87, 89,
	90, 92, 93, 95, 96, 98, 99, 101, 102, 104, 105, 107, 109, 110, 112, 114,
	115, 117, 119, 120, 122, 124, 126, 127, 129, 131, 133, 135, 137, 138, 140, 142,
	144, 146, 148, 150, 152, 154, 156, 158, 160, 162, 164, 167, 169, 171, 173, 175,
	177, 180, 182, 184, 186, 189, 191, 193, 196, 198, 200, 203, 205, 208, 210, 213,
	215, 218, 220, 223, 225, 228, 231, 233, 236, 239, 241, 244, 247, 249, 252, 255}

func gammaCorrected(color uint32) uint32 {
	r := (color >> 16) & 0xff
	g := (color >> 8) & 0xff
	b := (color >> 0) & 0xff
	return gamma8[r]<<16 + gamma8[g]<<8 + gamma8[b]
}

// NewNeoPixel create an instance of NeoPixel.
func NewNeoPixel(conf *Config) (np *NeoPixel, err error) {
	err = nil
	np = &NeoPixel{}
	np.Dev = (*C.ws2811_t)(C.malloc(C.sizeof_ws2811_t))
	if np == nil {
		err = errors.New("Unable to allocate memory")
		return
	}
	C.memset(unsafe.Pointer(np.Dev), 0, C.sizeof_ws2811_t)
	
	np.Config = conf

	np.Dev.freq = C.uint32_t(conf.Frequency)
	np.Dev.dmanum = C.int(conf.DmaNum)
	np.Dev.channel[0].gpionum = C.int(conf.GPIOPin)
	np.Dev.channel[0].count = C.int(conf.LEDCount)
	np.Dev.channel[0].brightness = C.uint8_t(conf.Brightness)
	np.Dev.channel[0].strip_type = C.int(conf.StripeType)
	if opt.Invert {
		np.Dev.channel[0].invert = C.int(1)
	} else {
		np.Dev.channel[0].invert = C.int(0)
	}
	return
}

// Init initialize the device. It should be called only once before any other method.
func (np *NeoPixel) Init() error {
	res := int(C.ws2811_init(np.Dev))
	if res == 0 {
		return nil
	}
	return fmt.Errorf("Error ws2811.init.%d", res)
}

// Render sends a complete frame to the LED Matrix
func (np *NeoPixel) Render() error {
	res := int(C.ws2811_render(np.Dev))
	if res == 0 {
		return nil
	}
	return fmt.Errorf("Error neopixel.render.%d", res)
}

// Wait waits for render to finish. The time needed for render is given by:
// time = 1/frequency * 8 * 3 * LedCount + 0.05
// (8 is the color depth and 3 is the number of colors (LEDs) per pixel).
// See https://cdn-shop.adafruit.com/datasheets/WS2811.pdf for more details.
func (np *NeoPixel) Wait() error {
	res := int(C.ws2811_wait(np.Dev))
	if res == 0 {
		return nil
	}
	return fmt.Errorf("Error ws2811.wait.%d", res)
}

// Fini shuts down the device.
func (np *NeoPixel) Fini() {
	C.ws2811_fini(np.Dev)
}

// SetLed defines the color of a given pixel.
func (np *NeoPixel) SetLed(index int, value uint32) {
	C.ws2811_set_led(np.Dev, 0, C.int(index), C.uint32_t(gammaCorrected(value)))
}

// SetBitmap defines the color of a all pixels.
func (np *NeoPixel) SetBitmap(a []uint32) {
	t := make([]uint32, len(a))
	for i, color := range a {
		t[i] = gammaCorrected(color)
	}
	C.ws2811_set_bitmap(np.Dev, 0, unsafe.Pointer(&t[0]), C.int(len(t)*4))
}

// Clear sets all pixels to black.
func (np *NeoPixel) Clear() {
	C.ws2811_clear_channel(np.Dev, 0)
}
