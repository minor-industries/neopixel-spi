package neopixel_spi

func ExpandBits(in []byte, out []byte) {
	bytePos := 0
	bitPos := 0
	for _, b := range in {
		bytePos, bitPos = packByte(out, bytePos, bitPos, b)
	}
}

func packByte(
	out []byte,
	bytePos int,
	bitPos int,
	b byte,
) (int, int) {
	for i := 0; i < 8; i++ {
		bit := b&(1<<i) != 0
		bytePos, bitPos = packBit(out, bytePos, bitPos, bit)
	}
	return bytePos, bitPos
}

func packBit(
	out []byte,
	bytePos int,
	bitPos int,
	bit bool,
) (int, int) {
	if bit {
		bytePos, bitPos = appendBit(out, bytePos, bitPos, true)
		bytePos, bitPos = appendBit(out, bytePos, bitPos, true)
		bytePos, bitPos = appendBit(out, bytePos, bitPos, false)
	} else {
		bytePos, bitPos = appendBit(out, bytePos, bitPos, true)
		bytePos, bitPos = appendBit(out, bytePos, bitPos, false)
		bytePos, bitPos = appendBit(out, bytePos, bitPos, false)
	}

	return bytePos, bitPos
}

func appendBit(
	out []byte,
	bytePos int,
	bitPos int,
	val bool,
) (int, int) {
	if bitPos == 8 {
		bitPos = 0
		bytePos++
	}

	if val {
		out[bytePos] |= 1 << bitPos
	} else {
		// let's explicitly clear the bit so we can used a fixed buffer in the future
		out[bytePos] &= ^(1 << bitPos)
	}

	bitPos++
	return bytePos, bitPos
}
