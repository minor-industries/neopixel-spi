package main

import (
	"device/sam"
	"fmt"
	"machine"
	"runtime/interrupt"
	"runtime/volatile"
	"time"
	"uc-go/app/neopixel-spi/driver"
)

var (
	ledPin = machine.PA23
)

var spi *machine.SPI

var d *driver.NeoSpiDriver

func spiInterruptHandler(i interrupt.Interrupt) {
	d.SpiInterruptHandler(i)
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

	d = &driver.NeoSpiDriver{Spi: spi}
	d.Init()

	spi.Bus.INTENSET.Set(sam.SERCOM_SPIM_INTENSET_DRE)
	d.Intr = interrupt.New(sam.IRQ_SERCOM5_0, spiInterruptHandler)
	d.Intr.Enable()

	if err != nil {
		forever(err)
	}

	for {
		d.Loop()
		fmt.Println("hello", volatile.LoadUint8(&d.InterruptCount))
	}

}

func forever(err error) {
	for {
		fmt.Println(err)
		<-time.After(time.Second)
	}
}
