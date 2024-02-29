package driver

import (
	"bytes"
	"device/sam"
	"machine"
	"runtime/interrupt"
	"sync/atomic"
	neopixel_spi "uc-go/pkg/neopixel-spi"
)

type NeoSpiDriver struct {
	Buf            []byte
	Spi            *machine.SPI
	Intr           interrupt.Interrupt
	InterruptCount uint64
	pos            int
}

var g = []byte{0x40, 0, 0}
var r = []byte{0, 0x40, 0}
var b = []byte{0, 0, 0x40}
var c = []byte{0, 0, 0}
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
	strip := bytes.Repeat(appendAll(r, r, b), 5)
	dmaStrip := make([]byte, len(strip)*3)
	neopixel_spi.ExpandBits(strip, dmaStrip)

	d.Buf = appendAll(
		dmaStrip,
		space,
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
	atomic.AddUint64(&d.InterruptCount, 1)

	d.pos++
	if d.pos >= len(d.Buf) {
		d.pos = 0
	}

	d.Spi.Bus.DATA.Set(uint32(d.Buf[d.pos]))
}
