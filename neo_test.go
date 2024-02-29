package neopixel_spi

import (
	"bytes"
	"fmt"
	"testing"
)

func Test_packBit(t *testing.T) {
	in := []byte{0x40, 0, 0}
	out := make([]byte, len(in)*3)
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
	//t.Skip()
	var g = []byte{0x40, 0, 0}
	strip := bytes.Repeat(g, 13)
	dmaStrip := make([]byte, len(strip)*3)
	ExpandBits(strip, dmaStrip)

	buf := appendAll(
		dmaStrip,
	)

	for i, b := range buf {
		fmt.Printf("%4d, %08b\n", i, b)
	}

}
