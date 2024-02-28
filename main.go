package main

import (
	"bytes"
	"device/sam"
	"fmt"
	"machine"
	"time"
	neopixel_spi "uc-go/pkg/neopixel-spi"
)

var (
	ledPin = machine.PA23
)

func main() {
	spi := machine.SPI{Bus: sam.SERCOM5_SPIM, SERCOM: 5}

	//machine.SPI0
	//.0 SDO MOSI
	//.1 CLK
	//.3 MISO

	err := spi.Configure(machine.SPIConfig{
		Frequency: 2_400_000,
		SCK:       machine.PA22, // 5.1 (sercom alt)
		SDO:       machine.PA23, // 5.0 (sercom alt)
		SDI:       machine.PA20, // 5.2 (sercom alt)
		LSBFirst:  true,
		Mode:      0,
	})
	if err != nil {
		forever(err)
	}

	g := neopixel_spi.ExpandBits([]byte{0x40, 0, 0})
	r := neopixel_spi.ExpandBits([]byte{0, 0x40, 0})
	b := neopixel_spi.ExpandBits([]byte{0, 0, 0x40})
	c := neopixel_spi.ExpandBits([]byte{0, 0, 0})

	space := bytes.Repeat([]byte{0}, 1000)

	var buf []byte
	buf = append(buf, space...)

	buf = append(buf, g...)
	buf = append(buf, g...)
	buf = append(buf, g...)

	buf = append(buf, c...)

	buf = append(buf, b...)
	buf = append(buf, b...)

	buf = append(buf, r...)

	buf = append(buf, c...)

	buf = append(buf, space...)

	for {
		err := spi.Tx(buf, nil)
		if err != nil {
			forever(err)
		}
	}
}

func forever(err error) {
	for {
		fmt.Println(err)
		<-time.After(time.Second)
	}
}
