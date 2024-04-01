package main

import (
	"device/sam"
	"fmt"
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

	d = driver.NewNeoSpiDriver(
		spi,
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

	ticker := time.NewTicker(1000 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			d.Animate()
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
