package driver

import "machine"

type Cfg struct {
	SPI *machine.SPI

	SCK machine.Pin
	SDO machine.Pin
	SDI machine.Pin

	LedCount   int
	SpaceCount int // TODO: rethink and/or remove
	EightBit   bool
}
