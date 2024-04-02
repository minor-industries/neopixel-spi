package default_driver

import (
	"device/sam"
	"github.com/minor-industries/uc-go/pkg/neopixel-spi/driver"
	"runtime/interrupt"
)

var defaultDriver *driver.NeoSpiDriver

func defaultDriverDREHandler(i interrupt.Interrupt) {
	defaultDriver.SpiInterruptHandler(i)
}

func defaultDriverTXCHandler(i interrupt.Interrupt) {
	defaultDriver.TxcInterruptHandler(i)
}

func Configure(cfg *driver.Cfg) *driver.NeoSpiDriver {
	defaultDriver = driver.NewNeoSpiDriver(cfg)

	// TODO: these IRQs shouldn't be hardcoded. Either computed or configured/overridden.
	interrupt.New(sam.IRQ_SERCOM5_0, defaultDriverDREHandler).Enable()
	interrupt.New(sam.IRQ_SERCOM5_1, defaultDriverTXCHandler).Enable()

	return defaultDriver
}
