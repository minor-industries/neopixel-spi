package driver

import (
	"image/color"
	"machine"
	"runtime/interrupt"
	"sync/atomic"
	"time"
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
	buf             []color.RGBA
	t0              time.Time
	frameNo         int
	orig            []color.RGBA
}

func NewNeoSpiDriver(spi *machine.SPI, intr interrupt.Interrupt, spaceCount int) *NeoSpiDriver {
	return &NeoSpiDriver{
		spi:             spi,
		intr:            intr,
		spaceCount:      spaceCount,
		spacesRemaining: spaceCount,
		t0:              time.Now(),
	}
}

var r = color.RGBA{0x40, 0, 0, 0}
var g = color.RGBA{0, 0x40, 0, 0}
var b = color.RGBA{0, 0, 0x40, 0}

func (d *NeoSpiDriver) Init() {
	d.orig = []color.RGBA{b, r, r, r, r, r, r, r, r, r, r, r, r, r, r, r, r, r}
	d.buf = make([]color.RGBA, len(d.orig))
	d.dmaBuf = make([]byte, len(d.buf)*9) // TODO: hide the details of this *9
	d.Animate()
}

func (d *NeoSpiDriver) Animate() {
	d.frameNo++
	for i := range d.orig {
		i2 := (i + d.frameNo) % len(d.buf)
		d.buf[i2] = d.orig[i]
	}

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
