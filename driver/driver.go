package driver

import (
	"device/sam"
	"image/color"
	"machine"
	"runtime/interrupt"
	"sync/atomic"
	"time"
	neopixel_spi "uc-go/pkg/neopixel-spi"
)

type NeoSpiDriver struct {
	dmaBuf          []uint32
	spi             *machine.SPI
	pos             int
	spaceCount      int
	spacesRemaining int
	buf             []color.RGBA
	t0              time.Time
	frameNo         int
	orig            []color.RGBA

	InterruptCount    uint64
	TXCInterruptCount uint64
	dreHandler        func(i interrupt.Interrupt)
	txcHandler        func(i interrupt.Interrupt)
}

func NewNeoSpiDriver(
	spi *machine.SPI,
	spaceCount int,
) *NeoSpiDriver {
	return &NeoSpiDriver{
		spi:        spi,
		spaceCount: spaceCount,

		spacesRemaining: spaceCount,
		t0:              time.Now(),
	}
}

var r = color.RGBA{0x40, 0, 0, 0}
var g = color.RGBA{0, 0x40, 0, 0}
var b = color.RGBA{0, 0, 0x40, 0}

func (d *NeoSpiDriver) Init() {
	// Disable SPI port.
	d.spi.Bus.CTRLA.ClearBits(sam.SERCOM_SPIM_CTRLA_ENABLE)
	for d.spi.Bus.SYNCBUSY.HasBits(sam.SERCOM_SPIM_SYNCBUSY_ENABLE) {
	}

	// set 32 bit mode
	d.spi.Bus.CTRLC.Set(sam.SERCOM_SPIM_CTRLC_DATA32B)

	// Enable SPI port.
	d.spi.Bus.CTRLA.SetBits(sam.SERCOM_SPIM_CTRLA_ENABLE)
	for d.spi.Bus.SYNCBUSY.HasBits(sam.SERCOM_SPIM_SYNCBUSY_ENABLE) {
	}

	d.spi.Bus.INTENSET.Set(sam.SERCOM_SPIM_INTENSET_DRE)
	d.spi.Bus.INTENSET.Set(sam.SERCOM_SPIM_INTENSET_TXC)

	d.orig = []color.RGBA{b, r, r, r, r, r, r, r, r, r, r, r, r, r, r, r, r, r}
	d.buf = make([]color.RGBA, len(d.orig))
	d.dmaBuf = make([]uint32, neopixel_spi.Bufsize32(len(d.buf)))
	d.Animate()
}

func (d *NeoSpiDriver) Animate() {
	d.frameNo++
	for i := range d.orig {
		i2 := (i + d.frameNo) % len(d.buf)
		d.buf[i2] = d.orig[i]
	}

	neopixel_spi.ExpandBits32(d.buf, d.dmaBuf)
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

func (d *NeoSpiDriver) TxcInterruptHandler(i interrupt.Interrupt) {
	atomic.AddUint64(&d.TXCInterruptCount, 1)
	d.spi.Bus.INTFLAG.Set(sam.SERCOM_SPIM_INTFLAG_TXC)
}
