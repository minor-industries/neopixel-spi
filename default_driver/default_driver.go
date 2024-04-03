package default_driver

import (
	"device/sam"
	"github.com/minor-industries/neopixel-spi"
	"runtime/interrupt"
)

var defaultDriver *neopixel_spi.NeoSpiDriver

func defaultDriverDREHandler(i interrupt.Interrupt) {
	defaultDriver.SpiInterruptHandler(i)
}

func defaultDriverTXCHandler(i interrupt.Interrupt) {
	defaultDriver.TxcInterruptHandler(i)
}

func Configure(cfg *neopixel_spi.Cfg) *neopixel_spi.NeoSpiDriver {
	defaultDriver = neopixel_spi.NewNeoSpiDriver(cfg)

	// TODO: these IRQs shouldn't be hardcoded. Either computed or configured/overridden.
	interrupt.New(sam.IRQ_SERCOM5_0, defaultDriverDREHandler).Enable()
	interrupt.New(sam.IRQ_SERCOM5_1, defaultDriverTXCHandler).Enable()

	return defaultDriver
}
