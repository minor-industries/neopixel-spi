package main

import (
	"device/sam"
	"fmt"
	"github.com/minor-industries/neopixel-spi/driver"
	"github.com/minor-industries/neopixel-spi/driver/default_driver"
	"image/color"
	"machine"
	"math"
	"time"
)

const (
	numLeds = 150
)

func main() {
	buf := make([]color.RGBA, numLeds)

	d := default_driver.Configure(&driver.Cfg{
		SPI:        &machine.SPI{Bus: sam.SERCOM5_SPIM, SERCOM: 5},
		SCK:        machine.PA22, // 5.1 (sercom alt)
		SDO:        machine.PA23, // 5.0 (sercom alt)
		SDI:        machine.PA20, // 5.2 (sercom alt)
		LedCount:   len(buf),
		SpaceCount: 2000,
		EightBit:   false,
	})

	if err := d.Init(); err != nil {
		forever(err)
	}

	t0 := time.Now()
	dt := 30 * time.Millisecond
	ticker := time.NewTicker(dt)
	for range ticker.C {
		// do the animation
		t := time.Now().Sub(t0).Seconds()
		animate(buf, d, t)
	}
}

func animate(
	buf []color.RGBA,
	d *driver.NeoSpiDriver,
	t float64,
) {
	const ledMaxLevel = 0.5

	convert := func(x float64) uint8 {
		return uint8(clamp(0, x, 1.0) * ledMaxLevel * 255.0)
	}

	for i := range buf {
		x := float64(i) * 0.1
		v := 0.5 - 0.5*math.Cos(t+x)

		buf[i].R = convert(v)
		buf[i].G = 0.0
		buf[i].B = 0.0
	}

	d.Animate(buf)
}

func forever(err error) {
	for {
		fmt.Println(err)
		<-time.After(time.Second)
	}
}

func clamp(a, x, b float64) float64 {
	if x < a {
		return a
	}

	if x > b {
		return b
	}

	return x
}
