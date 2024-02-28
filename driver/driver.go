package driver

import (
	"bytes"
	neopixel_spi "uc-go/pkg/neopixel-spi"
)

type NeoSpiDriver struct {
	Buf []byte
}

func (d *NeoSpiDriver) Init() {
	g := neopixel_spi.ExpandBits([]byte{0x40, 0, 0})
	r := neopixel_spi.ExpandBits([]byte{0, 0x40, 0})
	b := neopixel_spi.ExpandBits([]byte{0, 0, 0x40})
	c := neopixel_spi.ExpandBits([]byte{0, 0, 0})

	space := bytes.Repeat([]byte{0}, 1000)

	d.Buf = nil
	d.Buf = append(d.Buf, space...)

	d.Buf = append(d.Buf, c...)

	d.Buf = append(d.Buf, g...)
	d.Buf = append(d.Buf, g...)
	d.Buf = append(d.Buf, g...)

	d.Buf = append(d.Buf, c...)

	d.Buf = append(d.Buf, b...)
	d.Buf = append(d.Buf, b...)

	d.Buf = append(d.Buf, c...)

	d.Buf = append(d.Buf, r...)

	d.Buf = append(d.Buf, c...)

	d.Buf = append(d.Buf, r...)
	d.Buf = append(d.Buf, g...)
	d.Buf = append(d.Buf, b...)

	d.Buf = append(d.Buf, space...)
}
