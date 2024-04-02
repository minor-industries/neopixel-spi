package driver

import (
	"device/sam"
	neopixel_spi "github.com/minor-industries/uc-go/pkg/neopixel-spi"
	"github.com/pkg/errors"
	"image/color"
	"machine"
	"runtime/interrupt"
	"sync/atomic"
)

type NeoSpiDriver struct {
	dmaBuf32        []uint32
	dmaBuf8         []uint8
	spi             *machine.SPI
	pos             int
	spaceCount      int
	spacesRemaining int

	InterruptCount    uint64
	TXCInterruptCount uint64
	ledCount          int
	spiConfig         *machine.SPIConfig
	eightBit          bool
}

func NewNeoSpiDriver(cfg *Cfg) *NeoSpiDriver {
	return &NeoSpiDriver{
		eightBit:   cfg.EightBit,
		spi:        cfg.SPI,
		ledCount:   cfg.LedCount,
		spaceCount: cfg.SpaceCount,

		spacesRemaining: cfg.SpaceCount,
		spiConfig: &machine.SPIConfig{
			Frequency: 2_400_000,
			SCK:       cfg.SCK,
			SDO:       cfg.SDO,
			SDI:       cfg.SDI,
			LSBFirst:  true,
			Mode:      0,
		},
	}
}

func (d *NeoSpiDriver) Init() error {
	if err := d.spi.Configure(*d.spiConfig); err != nil {
		return errors.Wrap(err, "configure spi")
	}

	// Disable SPI port.
	d.spi.Bus.CTRLA.ClearBits(sam.SERCOM_SPIM_CTRLA_ENABLE)
	for d.spi.Bus.SYNCBUSY.HasBits(sam.SERCOM_SPIM_SYNCBUSY_ENABLE) {
	}

	if !d.eightBit {
		// set 32 bit mode
		d.spi.Bus.CTRLC.Set(sam.SERCOM_SPIM_CTRLC_DATA32B)
	}

	// Enable SPI port.
	d.spi.Bus.CTRLA.SetBits(sam.SERCOM_SPIM_CTRLA_ENABLE)
	for d.spi.Bus.SYNCBUSY.HasBits(sam.SERCOM_SPIM_SYNCBUSY_ENABLE) {
	}

	d.spi.Bus.INTENSET.Set(sam.SERCOM_SPIM_INTENSET_DRE)
	d.spi.Bus.INTENSET.Set(sam.SERCOM_SPIM_INTENSET_TXC)

	if d.eightBit {
		d.dmaBuf8 = make([]uint8, neopixel_spi.Bufsize(d.ledCount))
	} else {
		d.dmaBuf32 = make([]uint32, neopixel_spi.Bufsize32(d.ledCount))
	}

	return nil
}

func (d *NeoSpiDriver) Animate(buf []color.RGBA) {
	if d.eightBit {
		neopixel_spi.ExpandBits(buf, d.dmaBuf8)
	} else {
		neopixel_spi.ExpandBits32(buf, d.dmaBuf32)
	}
}

func (d *NeoSpiDriver) SpiInterruptHandler(i interrupt.Interrupt) {
	if d.eightBit {
		d.handle8()
	} else {
		d.handle32()
	}
}

func (d *NeoSpiDriver) handle8() {
	atomic.AddUint64(&d.InterruptCount, 1)

	if d.spacesRemaining > 0 {
		goto space
	} else {
		if d.pos >= len(d.dmaBuf8) {
			d.pos = 0
			d.spacesRemaining = d.spaceCount
			goto space
		} else {
			d.spi.Bus.DATA.Set(uint32(d.dmaBuf8[d.pos]))
			d.pos++
		}
	}

space:
	d.spacesRemaining--
	d.spi.Bus.DATA.Set(uint32(0))
}

func (d *NeoSpiDriver) handle32() {
	atomic.AddUint64(&d.InterruptCount, 1)

	if d.spacesRemaining > 0 {
		goto space
	} else {
		if d.pos >= len(d.dmaBuf32) {
			d.pos = 0
			d.spacesRemaining = d.spaceCount
			goto space
		} else {
			d.spi.Bus.DATA.Set(uint32(d.dmaBuf32[d.pos]))
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
