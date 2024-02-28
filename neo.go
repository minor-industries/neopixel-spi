package neopixel_spi

func ExpandBits(in []byte) (out []byte) {
	bitPos := 8
	for _, b := range in {
		out, bitPos = packByte(out, bitPos, b)
	}
	return out
}

func packByte(
	out []byte,
	bitPos int,
	b byte,
) ([]byte, int) {
	for i := 0; i < 8; i++ {
		bit := b&(1<<i) != 0
		out, bitPos = packBit(out, bitPos, bit)
	}
	return out, bitPos
}

func packBit(
	out []byte,
	bitPos int,
	bit bool,
) (
	[]byte,
	int,
) {
	if bit {
		out, bitPos = appendBit(out, bitPos, true)
		out, bitPos = appendBit(out, bitPos, true)
		out, bitPos = appendBit(out, bitPos, false)
	} else {
		out, bitPos = appendBit(out, bitPos, true)
		out, bitPos = appendBit(out, bitPos, false)
		out, bitPos = appendBit(out, bitPos, false)
	}

	return out, bitPos
}

func appendBit(out []byte, bitPos int, val bool) ([]byte, int) {
	if bitPos == 8 {
		out = append(out, 0)
		bitPos = 0
	}

	if val {
		out[len(out)-1] |= 1 << bitPos
	} else {
		// let's explicitly clear the bit so we can used a fixed buffer in the future
		out[len(out)-1] &= ^(1 << bitPos)
	}

	return out, bitPos + 1
}
