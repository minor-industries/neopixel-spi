package driver

import (
	"bytes"
	"machine"
	"runtime/interrupt"
	"sync/atomic"
	neopixel_spi "uc-go/pkg/neopixel-spi"
)

type NeoSpiDriver struct {
	dmaBuf          []byte
	spi             *machine.SPI
	intr            interrupt.Interrupt
	InterruptCount  uint64
	pos             int
	spaceCount      int
	spacesRemaining int
	buf             []byte
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

func appendAll(a0 []byte, as ...[]byte) []byte {
	var result []byte
	result = append(result, a0...)
	for _, a := range as {
		result = append(result, a...)
	}
	return result
}

func (d *NeoSpiDriver) Init() {
	d.buf = bytes.Repeat(appendAll(r, r, g), 5)
	d.dmaBuf = make([]byte, len(d.buf)*3)
	neopixel_spi.ExpandBits(d.buf, d.dmaBuf)
}

func (d *NeoSpiDriver) SpiInterruptHandler(i interrupt.Interrupt) {
	atomic.AddUint64(&d.InterruptCount, 1)

	if d.spacesRemaining > 0 {
		goto space
	} else {
		if d.pos >= len(d.dmaBuf) {
			d.pos = 0
			d.spacesRemaining = d.spaceCount
			goto space
		} else {
			d.spi.Bus.DATA.Set(uint32(d.dmaBuf[d.pos]))
			d.pos++
		}
	}

space:
	d.spacesRemaining--
	d.spi.Bus.DATA.Set(uint32(0))
}
