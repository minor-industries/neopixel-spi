package main

import (
	"device/sam"
	"fmt"
	"image/color"
	"machine"
	"time"
	"tinygo.org/x/drivers/irremote"
	"uc-go/app/neopixel-spi/driver"
	"uc-go/app/neopixel-spi/driver/default_driver"
)

// TODO:
// Try other IR receivers, IR performance

var (
	r = color.RGBA{0x40, 0, 0, 0}
	g = color.RGBA{0, 0x40, 0, 0}
	b = color.RGBA{0, 0, 0x40, 0}
)

func main() {
	orig := []color.RGBA{b, g, b, r, r, r, r, r, r, r, r, r, r, r, r, r, r, r}
	buf := make([]color.RGBA, len(orig))

	d := default_driver.Configure(&driver.Cfg{
		SPI:        &machine.SPI{Bus: sam.SERCOM5_SPIM, SERCOM: 5},
		SCK:        machine.PA22, // 5.1 (sercom alt)
		SDO:        machine.PA23, // 5.0 (sercom alt)
		SDI:        machine.PA20, // 5.2 (sercom alt)
		LedCount:   len(buf),
		SpaceCount: 2000,
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
	ticker := time.NewTicker(1000 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			// do the animation
			frameNo++
			for i := range orig {
				i2 := (i + frameNo) % len(buf)
				buf[i2] = orig[i]
			}
			d.Animate(buf)
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

func forever(err error) {
	for {
		fmt.Println(err)
		<-time.After(time.Second)
	}
}
