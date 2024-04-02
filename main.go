package main

import (
	"device/sam"
	"fmt"
	"image/color"
	"machine"
	"time"
	"tinygo.org/x/drivers/irremote"
	"uc-go/app/bikelights/cfg"
	"uc-go/pkg/leds/animations/bounce"
	"uc-go/pkg/leds/strip"
	"uc-go/pkg/neopixel-spi/driver"
	"uc-go/pkg/neopixel-spi/driver/default_driver"
	"uc-go/pkg/util"
)

// TODO:
// Try other IR receivers, IR performance

func main() {
	strip1 := strip.NewStrip(cfg.DefaultConfig)
	b := bounce.Bounce(&bounce.App{Strip: strip1})
	buf := make([]color.RGBA, cfg.DefaultConfig.NumLeds)

	d := default_driver.Configure(&driver.Cfg{
		SPI:        &machine.SPI{Bus: sam.SERCOM5_SPIM, SERCOM: 5},
		SCK:        machine.PA22, // 5.1 (sercom alt)
		SDO:        machine.PA23, // 5.0 (sercom alt)
		SDI:        machine.PA20, // 5.2 (sercom alt)
		LedCount:   cfg.DefaultConfig.NumLeds,
		SpaceCount: 2000,
		EightBit:   false,
	})

	if err := d.Init(); err != nil {
		forever(err)
	}

	irPin := machine.D5

	ir := irremote.NewReceiver(irPin)
	ir.Configure()

	ch := make(chan irremote.Data, 10)
	ir.SetCommandHandler(func(data irremote.Data) {
		ch <- data
	})

	irCount := 0

	frameNo := 0

	t0 := time.Now()
	dt := 30 * time.Millisecond
	ticker := time.NewTicker(dt)
	for {
		select {
		case <-ticker.C:
			// do the animation
			frameNo++
			b.Tick(0.0, dt.Seconds())
			t := time.Now().Sub(t0).Seconds()
			animate(strip1, buf, d, frameNo, t)

		case data := <-ch:
			if data.Flags&irremote.DataFlagIsRepeat != 0 {
				continue
			}
			switch data.Command {
			case 16:
				irCount++
			case 17:
				irCount--
			}
			fmt.Println(irCount)
		}
	}
}

func animate(
	strip1 *strip.Strip,
	buf []color.RGBA,
	d *driver.NeoSpiDriver,
	frameNo int,
	t float64,
) {
	const ledMaxLevel = 0.5
	const scale = 1.0

	convert := func(x float32) uint8 {
		val := x * scale
		return uint8(util.Clamp(0, val, 1.0) * ledMaxLevel * 255.0)
	}

	strip1.Each(func(i int, led *strip.Led) {
		buf[i].R = convert(led.R)
		buf[i].G = convert(led.G)
		buf[i].B = convert(led.B)
	})

	d.Animate(buf)
}

func forever(err error) {
	for {
		fmt.Println(err)
		<-time.After(time.Second)
	}
}
