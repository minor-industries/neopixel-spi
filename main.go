package main

import (
	"device/sam"
	"fmt"
	"image/color"
	"machine"
	"runtime/interrupt"
	"time"
	"tinygo.org/x/drivers/irremote"
	"uc-go/app/neopixel-spi/driver"
)

// TODO:
// Try other IR receivers, IR performance

var (
	ledPin = machine.PA23
)

var spi *machine.SPI

var d *driver.NeoSpiDriver

func dreHandler(i interrupt.Interrupt) {
	d.SpiInterruptHandler(i)
}

func txcHandler(i interrupt.Interrupt) {
	d.TxcInterruptHandler(i)
}

func main() {
	spi = &machine.SPI{Bus: sam.SERCOM5_SPIM, SERCOM: 5}

	err := spi.Configure(machine.SPIConfig{
		Frequency: 2_400_000,
		SCK:       machine.PA22, // 5.1 (sercom alt)
		SDO:       machine.PA23, // 5.0 (sercom alt)
		SDI:       machine.PA20, // 5.2 (sercom alt)
		LSBFirst:  true,
		Mode:      0,
	})
	if err != nil {
		forever(err)
	}

	var r = color.RGBA{0x40, 0, 0, 0}
	var g = color.RGBA{0, 0x40, 0, 0}
	var b = color.RGBA{0, 0, 0x40, 0}

	orig := []color.RGBA{b, g, r, r, r, r, r, r, r, r, r, r, r, r, r, r, r, r}
	buf := make([]color.RGBA, len(orig))

	d = driver.NewNeoSpiDriver(
		spi,
		len(buf),
		2000,
	)

	interrupt.New(sam.IRQ_SERCOM5_0, dreHandler).Enable()
	interrupt.New(sam.IRQ_SERCOM5_1, txcHandler).Enable()

	d.Init()

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
