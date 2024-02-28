package main

import (
	"bytes"
	"device/sam"
	"fmt"
	"machine"
	"runtime/interrupt"
	"runtime/volatile"
	"time"
	neopixel_spi "uc-go/pkg/neopixel-spi"
)

var (
	ledPin = machine.PA23
)

var interruptCount uint8 = 0
var spi *machine.SPI
var intr interrupt.Interrupt

func spiInterruptHandler(i interrupt.Interrupt) {
	intr.Disable()
	volatile.StoreUint8(&interruptCount, 1)
	spi.Bus.INTFLAG.Set(sam.SERCOM_SPIM_INTFLAG_DRE)
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

	spi.Bus.INTENSET.Set(sam.SERCOM_SPIM_INTENSET_DRE)
	intr = interrupt.New(sam.IRQ_SERCOM5_0, spiInterruptHandler)
	intr.Enable()

	if err != nil {
		forever(err)
	}

	g := neopixel_spi.ExpandBits([]byte{0x40, 0, 0})
	r := neopixel_spi.ExpandBits([]byte{0, 0x40, 0})
	b := neopixel_spi.ExpandBits([]byte{0, 0, 0x40})
	c := neopixel_spi.ExpandBits([]byte{0, 0, 0})

	space := bytes.Repeat([]byte{0}, 1000)

	var buf []byte
	buf = append(buf, space...)

	buf = append(buf, c...)

	buf = append(buf, g...)
	buf = append(buf, g...)
	buf = append(buf, g...)

	buf = append(buf, c...)

	buf = append(buf, b...)
	buf = append(buf, b...)

	buf = append(buf, c...)

	buf = append(buf, r...)

	buf = append(buf, c...)

	buf = append(buf, space...)

	for {
		loop(buf)
		fmt.Println("hello", volatile.LoadUint8(&interruptCount))
	}

}

func loop(buf []byte) {
	i := 0
	count := 100
	for count > 0 {
		for !spi.Bus.INTFLAG.HasBits(sam.SERCOM_SPIM_INTFLAG_DRE) {
		}
		val := buf[i]
		spi.Bus.DATA.Set(uint32(val))
		i++
		if i >= len(buf) {
			i = 0
			count--
		}
	}
}

func forever(err error) {
	for {
		fmt.Println(err)
		<-time.After(time.Second)
	}
}
