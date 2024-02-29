package driver

import (
	"bytes"
	"device/sam"
	"machine"
	"runtime/interrupt"
	"runtime/volatile"
	neopixel_spi "uc-go/pkg/neopixel-spi"
)

type NeoSpiDriver struct {
	Buf            []byte
	Spi            *machine.SPI
	Intr           interrupt.Interrupt
	InterruptCount uint8
}

var g = neopixel_spi.ExpandBits([]byte{0x40, 0, 0})
var r = neopixel_spi.ExpandBits([]byte{0, 0x40, 0})
var b = neopixel_spi.ExpandBits([]byte{0, 0, 0x40})
var c = neopixel_spi.ExpandBits([]byte{0, 0, 0})
var space = bytes.Repeat([]byte{0}, 1000)

func appendAll(a0 []byte, as ...[]byte) []byte {
	var result []byte
	result = append(result, a0...)
	for _, a := range as {
		result = append(result, a...)
	}
	return result
}

func (d *NeoSpiDriver) Init() {
	d.Buf = appendAll(
		space,
		r, g, b, r,
		bytes.Repeat(c, 13),
		space,
	)
}

func (d *NeoSpiDriver) Loop() {
	i := 0
	count := 100
	for count > 0 {
		for !d.Spi.Bus.INTFLAG.HasBits(sam.SERCOM_SPIM_INTFLAG_DRE) {
		}
		val := d.Buf[i]
		d.Spi.Bus.DATA.Set(uint32(val))
		i++
		if i >= len(d.Buf) {
			i = 0
			count--
		}
	}
}

func (d *NeoSpiDriver) SpiInterruptHandler(i interrupt.Interrupt) {
	d.Intr.Disable()
	volatile.StoreUint8(&d.InterruptCount, 1)
	//d.Spi.Bus.INTFLAG.Set(sam.SERCOM_SPIM_INTFLAG_DRE) // this doesn't do anything
}
