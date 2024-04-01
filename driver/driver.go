package driver

import (
	"device/sam"
	"image/color"
	"machine"
	"runtime/interrupt"
	"sync/atomic"
	neopixel_spi "uc-go/pkg/neopixel-spi"
)

type NeoSpiDriver struct {
	dmaBuf          []uint32
	spi             *machine.SPI
	pos             int
	spaceCount      int
	spacesRemaining int

	InterruptCount    uint64
	TXCInterruptCount uint64
	ledCount          int
}

func NewNeoSpiDriver(
	spi *machine.SPI,
	ledCount int,
	spaceCount int,
) *NeoSpiDriver {
	return &NeoSpiDriver{
		spi:        spi,
		ledCount:   ledCount,
		spaceCount: spaceCount,

		spacesRemaining: spaceCount,
	}
}

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

	d.dmaBuf = make([]uint32, neopixel_spi.Bufsize32(d.ledCount))
}

func (d *NeoSpiDriver) Animate(buf []color.RGBA) {
	neopixel_spi.ExpandBits32(buf, d.dmaBuf)
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
