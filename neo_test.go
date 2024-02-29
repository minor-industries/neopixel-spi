package neopixel_spi

import (
	"fmt"
	"image/color"
	"testing"
)

var r = color.RGBA{0x40, 0, 0, 0}

func Test_packBit(t *testing.T) {
	in := []color.RGBA{r}
	out := make([]byte, len(in)*3*3)
	ExpandBits(in, out)
	for i, b := range out {
		fmt.Printf("%4d, %08b\n", i, b)
	}
}
