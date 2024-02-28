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

	expanded := neopixel_spi.ExpandBits([]byte{0xFF, 0xFF, 0xFF})

	space := bytes.Repeat([]byte{0}, 1000)

	var buf []byte
	buf = append(buf, space...)
	buf = append(buf, expanded...)
	buf = append(buf, expanded...)
	buf = append(buf, space...)

	for range time.NewTicker(30 * time.Millisecond).C {
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
