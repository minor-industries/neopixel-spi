package bits

import (
	"image/color"
)

func Bufsize(sz int) int {
	return (sz*9 + 3) / 4
}

func ExpandBits(in []color.RGBA, out []uint32) {
	outIndex := 0
	bitPos := 0
	for _, c := range in {
		outIndex, bitPos = packByte(out, outIndex, bitPos, c.G)
		outIndex, bitPos = packByte(out, outIndex, bitPos, c.R)
		outIndex, bitPos = packByte(out, outIndex, bitPos, c.B)
	}
}

func packByte(
	out []uint32,
	outIndex int,
	bitPos int,
	b byte,
) (int, int) {
	for i := 7; i >= 0; i-- {
		bit := b&(1<<i) != 0
		outIndex, bitPos = packBit(out, outIndex, bitPos, bit)
	}
	return outIndex, bitPos
}

func packBit(
	out []uint32,
	outIndex int,
	bitPos int,
	bit bool,
) (int, int) {
	if bit {
		outIndex, bitPos = appendBit(out, outIndex, bitPos, true)
		outIndex, bitPos = appendBit(out, outIndex, bitPos, true)
		outIndex, bitPos = appendBit(out, outIndex, bitPos, false)
	} else {
		outIndex, bitPos = appendBit(out, outIndex, bitPos, true)
		outIndex, bitPos = appendBit(out, outIndex, bitPos, false)
		outIndex, bitPos = appendBit(out, outIndex, bitPos, false)
	}

	return outIndex, bitPos
}

func appendBit(
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
		out[outIndex] &= ^(1 << bitPos)
	}

	bitPos++
	return outIndex, bitPos
}
