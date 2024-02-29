package neopixel_spi

import (
	"fmt"
	"image/color"
	"testing"
)

func Test_packBit32(t *testing.T) {
	in := []color.RGBA{r, r, r, r}

	outSize := ((len(in) * 9) + 3) / 4

	out := make([]uint32, outSize)
	ExpandBits32(in, out)
	for i, b := range out {
		fmt.Printf("%4d, %032b\n", i, b)
	}
}
