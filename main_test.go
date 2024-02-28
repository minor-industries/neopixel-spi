package neopixel_spi

import (
	"fmt"
	"testing"
)

func Test_packBit(t *testing.T) {
	out := ExpandBits([]byte{0x05})
	for i, b := range out {
		fmt.Printf("%4d, %08b\n", i, b)
	}
}
