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
		EightBit:   true,
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
	//numLEDs := len(buf)
	//
	////i := 0
	//
	//for j := 0; j < numLEDs; j++ {
	//	buf[j] = wheel((int32(j)*256/int32(numLEDs) + int32(frameNo)) & 255)
	//}
	//
	//strip1.Each(func(i int, led *strip.Led) {
	//	buf[i].R = uint8(5 * led.R)
	//	buf[i].G = uint8(5 * led.G)
	//	buf[i].B = uint8(5 * led.B)
	//})

	for i := range buf {
		buf[i].R = uint8(i)
		buf[i].G = 0
		buf[i].B = 0
	}

	d.Animate(buf)
}

func wheel(pos int32) color.RGBA {
	if pos < 85 {
		return color.RGBA{uint8(pos * 3), uint8(255 - pos*3), 0, 255}
	} else if pos < 170 {
		pos -= 85
		return color.RGBA{uint8(255 - pos*3), 0, uint8(pos * 3), 255}
	} else {
		pos -= 170
		return color.RGBA{0, uint8(pos * 3), uint8(255 - pos*3), 255}
	}
}

func forever(err error) {
	for {
		fmt.Println(err)
		<-time.After(time.Second)
	}
}
