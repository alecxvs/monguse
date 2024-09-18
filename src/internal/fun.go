package internal

import (
	"math/rand/v2"

	"github.com/gdamore/tcell/v2"
)

var AMOGUS_COLORS []int32 = []int32{
	0xD71E22,
	0x1D3CE9,
	0x1B913E,
	0xFF63D4,
	0xFF8D1C,
	0xFFFF67,
	0x4A565E,
	0xE9F7FF,
	0x783DD2,
	0x80582D,
	0x44FFF7,
	0x5BFE4B,
	0x6C2B3D,
	0xFFD6EC,
	0xFFFFBE,
	0x8397A7,
	0x9F9989,
	0xEC7578,
}

func RandomAmogusColor() tcell.Color {
	return tcell.NewHexColor(AMOGUS_COLORS[rand.IntN(18)])
}
