package driver

import (
	"device/sam"
	neopixel_spi "github.com/minor-industries/neopixel-spi/internal/bits"
	"github.com/pkg/errors"
	"image/color"
	"machine"
	"runtime/interrupt"
	"sync/atomic"
)

type NeoSpiDriver struct {
	InterruptCount    uint64 // needs 64-bit alignment
	TXCInterruptCount uint64 // needs 64-bit alignment

	dmaBuf32        []uint32
	spi             *machine.SPI
	pos             int
	spaceCount      int
	spacesRemaining int

	ledCount  int
	spiConfig *machine.SPIConfig
}

func NewNeoSpiDriver(cfg *Cfg) *NeoSpiDriver {
	return &NeoSpiDriver{
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

	// Disable SPI port
	d.spi.Bus.CTRLA.ClearBits(sam.SERCOM_SPIM_CTRLA_ENABLE)
	for d.spi.Bus.SYNCBUSY.HasBits(sam.SERCOM_SPIM_SYNCBUSY_ENABLE) {
	}

	// set 32 bit mode
	d.spi.Bus.CTRLC.Set(sam.SERCOM_SPIM_CTRLC_DATA32B)

	// Enable SPI port
	d.spi.Bus.CTRLA.SetBits(sam.SERCOM_SPIM_CTRLA_ENABLE)
	for d.spi.Bus.SYNCBUSY.HasBits(sam.SERCOM_SPIM_SYNCBUSY_ENABLE) {
	}

	d.spi.Bus.INTENSET.Set(sam.SERCOM_SPIM_INTENSET_DRE)
	d.spi.Bus.INTENSET.Set(sam.SERCOM_SPIM_INTENSET_TXC)

	d.dmaBuf32 = make([]uint32, neopixel_spi.Bufsize(d.ledCount))

	return nil
}

func (d *NeoSpiDriver) Animate(buf []color.RGBA) {
	neopixel_spi.ExpandBits(buf, d.dmaBuf32)
}

func (d *NeoSpiDriver) SpiInterruptHandler(i interrupt.Interrupt) {
	atomic.AddUint64(&d.InterruptCount, 1)

	if d.spacesRemaining > 0 {
		goto space
	} else {
		if d.pos >= len(d.dmaBuf32) {
			d.pos = 0
			d.spacesRemaining = d.spaceCount
			goto space
		} else {
			d.spi.Bus.DATA.Set(d.dmaBuf32[d.pos])
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
