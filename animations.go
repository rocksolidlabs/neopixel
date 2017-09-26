package neopixel

import (
	"time"
)

// Takes RGB as ints (0-255) and returns a uint32 to be used in other color functions.
func Color(red, green, blue int) uint32 {
	white := 0
	return uint32((white << 24) | (red << 16) | (green << 8) | blue)
}

// Wipe color across display a pixel at a time.
func (np *NeoPixel) ColorWipe(color uint32, wait int) {
	for i := 0; i < np.Config.LEDCount; i++ {
		np.SetLED(i, color)
		np.Render()
		np.Wait()
		time.Sleep(time.Duration(wait) * time.Millisecond)
	}
}

// Movie theater light style chaser animation.
func (np *NeoPixel) TheaterChase(color uint32, wait, iterations int) {
	for j := 0; j < iterations; j++ {
		for q := 0; q < 3; q++ {
			for i := 0; i < np.Config.LEDCount; i++ {
				if i%3 == 0 {
					np.SetLED(i+q, color)
				}
			}
			np.Render()
			np.Wait()
			time.Sleep(time.Duration(wait) * time.Millisecond)
			for i := 0; i < np.Config.LEDCount; i++ {
				if i%3 == 0 {
					np.SetLED(i+q, 0)
				}
			}
			np.Render()
			np.Wait()
		}
		time.Sleep(time.Duration(wait) * time.Millisecond)
	}
}

// Generate rainbow colors across 0-255 positions.
func wheel(position int) uint32 {
	if position < 85 {
		return Color(position*3, 255-position*3, 0)
	} else if position < 170 {
		position = position - 85
		return Color(255-position, 0, position*3)
	} else {
		position = position - 170
		return Color(0, position*3, 255-position*3)
	}
	return 0
}

// Draw a rainbow that fades across all pixels at once.
func (np *NeoPixel) Rainbow(wait, iterations int) {
	for j := 0; j < 256*iterations; j++ {
		for i := 0; i < np.Config.LEDCount; i++ {
			np.SetLED(i, wheel((i+j)&255))
		}
		np.Render()
		np.Wait()
		time.Sleep(time.Duration(wait/1000) * time.Millisecond)
	}
}

//Draw rainbow that uniformly distributes itself across all pixels.
func (np *NeoPixel) RainbowCycle(wait, iterations int) {
	for j := 0; j < 256*iterations; j++ {
		for i := 0; i < np.Config.LEDCount; i++ {
			np.SetLED(i, wheel(((i*256/np.Config.LEDCount)+j)&255))
		}
		np.Render()
		np.Wait()
		time.Sleep(time.Duration(wait/1000) * time.Millisecond)
	}
}

// Rainbow movie theater light style chaser animation.
func (np *NeoPixel) TheaterChaseRainbow(wait, iterations int) {
	for j := 0; j < iterations; j++ {
		for q := 0; q < 3; q++ {
			for i := 0; i < np.Config.LEDCount; i++ {
				if i%3 == 0 {
					np.SetLED(i+q, wheel((i + j%255)))
				}
			}
			np.Render()
			np.Wait()
			time.Sleep(time.Duration(wait) * time.Millisecond)
			for i := 0; i < np.Config.LEDCount; i++ {
				if i%3 == 0 {
					np.SetLED(i+q, 0)
				}
			}
			np.Render()
			np.Wait()
		}
		time.Sleep(time.Duration(wait) * time.Millisecond)
	}
}
