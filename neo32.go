package neopixel_spi

import (
	"image/color"
)

func ExpandBits32(in []color.RGBA, out []uint32) {
	outIndex := 0
	bitPos := 0
	for _, c := range in {
		outIndex, bitPos = packByte32(out, outIndex, bitPos, c.G)
		outIndex, bitPos = packByte32(out, outIndex, bitPos, c.R)
		outIndex, bitPos = packByte32(out, outIndex, bitPos, c.B)
	}
}

func packByte32(
	out []uint32,
	outIndex int,
	bitPos int,
	b byte,
) (int, int) {
	for i := 0; i < 8; i++ {
		bit := b&(1<<i) != 0
		outIndex, bitPos = packBit32(out, outIndex, bitPos, bit)
	}
	return outIndex, bitPos
}

func packBit32(
	out []uint32,
	outIndex int,
	bitPos int,
	bit bool,
) (int, int) {
	if bit {
		outIndex, bitPos = appendBit32(out, outIndex, bitPos, true)
		outIndex, bitPos = appendBit32(out, outIndex, bitPos, true)
		outIndex, bitPos = appendBit32(out, outIndex, bitPos, false)
	} else {
		outIndex, bitPos = appendBit32(out, outIndex, bitPos, true)
		outIndex, bitPos = appendBit32(out, outIndex, bitPos, false)
		outIndex, bitPos = appendBit32(out, outIndex, bitPos, false)
	}

	return outIndex, bitPos
}

func appendBit32(
	out []uint32,
	outIndex int,
	bitPos int,
	val bool,
) (int, int) {
	if bitPos == 32 {
		bitPos = 0
		outIndex++
	}

	if val {
		out[outIndex] |= 1 << bitPos
	} else {
		// let's explicitly clear the bit so we can used a fixed buffer in the future
		out[outIndex] &= ^(1 << bitPos)
	}

	bitPos++
	return outIndex, bitPos
}
