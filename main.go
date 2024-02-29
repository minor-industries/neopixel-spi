package main

import (
	"device/sam"
	"fmt"
	"machine"
	"runtime/interrupt"
	"sync/atomic"
	"time"
	"tinygo.org/x/drivers/irremote"
	"uc-go/app/neopixel-spi/driver"
)

// TODO:
// - Basic animations
// - Get counts of transmit complete interrupts (TXC)
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

func txcInterruptHandler(i interrupt.Interrupt) {
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

	// Disable SPI port.
	spi.Bus.CTRLA.ClearBits(sam.SERCOM_SPIM_CTRLA_ENABLE)
	for spi.Bus.SYNCBUSY.HasBits(sam.SERCOM_SPIM_SYNCBUSY_ENABLE) {
	}

	// set 32 bit mode
	spi.Bus.CTRLC.Set(sam.SERCOM_SPIM_CTRLC_DATA32B)

	// Enable SPI port.
	spi.Bus.CTRLA.SetBits(sam.SERCOM_SPIM_CTRLA_ENABLE)
	for spi.Bus.SYNCBUSY.HasBits(sam.SERCOM_SPIM_SYNCBUSY_ENABLE) {
	}

	t0 := time.Now()

	intr := interrupt.New(sam.IRQ_SERCOM5_0, spiInterruptHandler)
	txcIntr := interrupt.New(sam.IRQ_SERCOM5_1, txcInterruptHandler)
	_ = txcIntr

	d = driver.NewNeoSpiDriver(spi, 2000)
	d.Init()

	intr.Enable()
	txcIntr.Enable()

	spi.Bus.INTENSET.Set(sam.SERCOM_SPIM_INTENSET_DRE)
	spi.Bus.INTENSET.Set(sam.SERCOM_SPIM_INTENSET_TXC)

	irPin := machine.PA15

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
			switch data.Command {
			case 16:
				irCount++
			case 17:
				irCount--
			}
			dreCount := atomic.LoadUint64(&d.InterruptCount)
			dt := time.Now().Sub(t0).Seconds()

			fmt.Println(atomic.LoadUint64(&d.TXCInterruptCount), dreCount, float64(dreCount)/dt, irCount)
		}
	}
}

func forever(err error) {
	for {
		fmt.Println(err)
		<-time.After(time.Second)
	}
}
