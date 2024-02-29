package main

import (
	"device/sam"
	"fmt"
	"machine"
	"runtime/interrupt"
	"sync/atomic"
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
	intr := interrupt.New(sam.IRQ_SERCOM5_0, spiInterruptHandler)

	d = driver.NewNeoSpiDriver(spi, intr, 1e6)
	d.Init()

	spi.Bus.INTENSET.Set(sam.SERCOM_SPIM_INTENSET_DRE)
	intr.Enable()

	if err != nil {
		forever(err)
	}

	for range time.NewTicker(time.Second).C {
		//fmt.Println("hello", atomic.LoadUint64(&d.InterruptCount))
	}

	for {
		d.Loop()
		fmt.Println("hello", atomic.LoadUint64(&d.InterruptCount))
	}

}

func forever(err error) {
	for {
		fmt.Println(err)
		<-time.After(time.Second)
	}
}
