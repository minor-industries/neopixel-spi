package bits

import (
	"fmt"
	"image/color"
	"testing"
)

func Test_packBit32(t *testing.T) {
	in := []color.RGBA{r, r, r, r, r, r, r, r, r}

	outSize := Bufsize(len(in))

	out := make([]uint32, outSize)
	ExpandBits(in, out)
	for i, b := range out {
		fmt.Printf("%4d, %032b\n", i, b)
	}
}
