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

func appendAll(a0 []byte, as ...[]byte) []byte {
	var result []byte
	result = append(result, a0...)
	for _, a := range as {
		result = append(result, a...)
	}
	return result
}

func Test_It(t *testing.T) {
	strip := []color.RGBA{r, r, r, r, r, r, r, r, r, r, r, r, r}
	dmaStrip := make([]byte, len(strip)*3*3)
	ExpandBits(strip, dmaStrip)

	buf := appendAll(
		dmaStrip,
	)

	for i, b := range buf {
		fmt.Printf("%4d, %08b\n", i, b)
	}

}
