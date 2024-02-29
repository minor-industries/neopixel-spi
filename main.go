package main

import (
	"device/sam"
	"fmt"
	"machine"
	"runtime/interrupt"
	"time"
	"uc-go/app/neopixel-spi/driver"
)

// TODO:
// - Basic animations
// - Add IR
// - 32-bit extension

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
	if err != nil {
		forever(err)
	}

	intr := interrupt.New(sam.IRQ_SERCOM5_0, spiInterruptHandler)

	d = driver.NewNeoSpiDriver(spi, intr, 2000)
	d.Init()

	spi.Bus.INTENSET.Set(sam.SERCOM_SPIM_INTENSET_DRE)
	intr.Enable()

	for range time.NewTicker(100 * time.Millisecond).C {
		d.Animate()
	}
}

func forever(err error) {
	for {
		fmt.Println(err)
		<-time.After(time.Second)
	}
}
