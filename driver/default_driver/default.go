package default_driver

import (
	"device/sam"
	"runtime/interrupt"
	driver2 "uc-go/pkg/neopixel-spi/driver"
)

var defaultDriver *driver2.NeoSpiDriver

func defaultDriverDREHandler(i interrupt.Interrupt) {
	defaultDriver.SpiInterruptHandler(i)
}

func defaultDriverTXCHandler(i interrupt.Interrupt) {
	defaultDriver.TxcInterruptHandler(i)
}

func Configure(cfg *driver2.Cfg) *driver2.NeoSpiDriver {
	defaultDriver = driver2.NewNeoSpiDriver(cfg)

	// TODO: these IRQs shouldn't be hardcoded. Either computed or configured/overridden.
	interrupt.New(sam.IRQ_SERCOM5_0, defaultDriverDREHandler).Enable()
	interrupt.New(sam.IRQ_SERCOM5_1, defaultDriverTXCHandler).Enable()

	return defaultDriver
}
