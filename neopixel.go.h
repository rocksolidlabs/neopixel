// This header file is a wrapper for the rpi_ws281x library:
// https://github.com/jgarff/rpi_ws281x
 
#include <stdint.h>
#include <string.h>
#include <ws2811.h>

void ws2811_set_led(ws2811_t *ws2811, int chan, int index, uint32_t value) {
	ws2811->channel[chan].leds[index] = value;
}

void ws2811_clear_channel(ws2811_t *ws2811, int chan) {
	ws2811_channel_t *channel = &ws2811->channel[chan];
	memset(channel->leds, 0, sizeof(ws2811_led_t) * channel->count);
}

void ws2811_clear_all(ws2811_t *ws2811) {
	for (int chan = 0; chan < RPI_PWM_CHANNELS; chan++) {
		ws2811_clear_channel(ws2811, chan);
	}
}

void ws2811_set_bitmap(ws2811_t *ws2811, int chan, void* a, int len) {
	memcpy(ws2811->channel[chan].leds, a, len);
}
