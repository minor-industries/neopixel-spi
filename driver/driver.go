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
	buf             []byte
	spi             *machine.SPI
	intr            interrupt.Interrupt
	InterruptCount  uint64
	pos             int
	spaceCount      int
	spacesRemaining int
}

func NewNeoSpiDriver(spi *machine.SPI, intr interrupt.Interrupt, spaceCount int) *NeoSpiDriver {
	return &NeoSpiDriver{
		spi:             spi,
		intr:            intr,
		spaceCount:      spaceCount,
		spacesRemaining: spaceCount,
	}
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
	strip := bytes.Repeat(appendAll(r, b, b), 5)
	dmaStrip := make([]byte, len(strip)*3)
	neopixel_spi.ExpandBits(strip, dmaStrip)

	d.buf = appendAll(
		space,
		dmaStrip,
		space,
	)
}

func (d *NeoSpiDriver) Loop() {
	i := 0
	count := 100
	for count > 0 {
		for !d.spi.Bus.INTFLAG.HasBits(sam.SERCOM_SPIM_INTFLAG_DRE) {
		}
		val := d.buf[i]
		d.spi.Bus.DATA.Set(uint32(val))
		i++
		if i >= len(d.buf) {
			i = 0
			count--
		}
	}
}

func (d *NeoSpiDriver) SpiInterruptHandler(i interrupt.Interrupt) {
	atomic.AddUint64(&d.InterruptCount, 1)

	if d.spacesRemaining > 0 {
		d.spacesRemaining--
		d.spi.Bus.DATA.Set(uint32(0))
	}

	d.pos++
	if d.pos >= len(d.buf) {
		d.pos = 0
	}

	d.spi.Bus.DATA.Set(uint32(d.buf[d.pos]))
}
