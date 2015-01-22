package main

import (
	"fmt"
)

//var seed1 = []byte{0xD3, 0xFB, 0x8C} // Key: AB D9 26?
//var seed2 = []byte{0xAB, 0xED, 0xCC} // Key: 84 C4 EA?

var seed1 = 0xD3FB8C
var seed2 = 0xABEDCC
var seed3 = 0x5BD112

var seeds = []int{
	0xD3FB8C,
	0xABEDCC,
	0x73B3E0,
	0x03FF7E,
	0xE3E7CA,
	0x7B697E,
	0xB363C8,
	0x5B51D2,
	0xC38FD6,
	0x63873A,
	0xFB098E,
	0x3343D8,
	0x2B0D5C,
}

var keys = []int{
	0xABD9,
	0x84C4,
	0x492F,
	0x9A01,
	0x011B,
	0x901F,
	0x3935,
	0x6A0B,
	0x4976,
	0xCBB9,
	0x108F,
	0x1950,
	0x44F4,
}

// Main Function
//
//
func main() {
	for i, seed := range seeds {
		key := getKey(seed)
		fmt.Printf("Seed%v: [%X]	Key: [%X]	Expected Key: %X\n\n\n", i, seed, key, keys[i])

		key2 := (seed ^ (seed >> 8) ^ 0x9B) + 43314
		fmt.Printf("Key2-%v: [%X]\n\n", 1, key2)
	}
}

// getKey Function
//
//

func getKey(seed int) int {
	//key = abs((signed __int16)((_WORD)seed * 2 * seed - 253 * seed));
	b1 := byte1(seed)

	case1 := (b1 >> 0) & 3
	fmt.Printf("CASE1 Compare [%X]\n", case1)

	case2 := (b1 >> 2) & 3
	fmt.Printf("CASE2 Compare [%X]\n", case2)

	//muck1 := uint8(muck((seed>>8)&0xFF, int(uint8(seed)), (int(uint8(seed)>>0))&3))
	//muck2 := uint8(muck((seed>>16)&0xFF, int(uint8(seed)), (int(uint8(seed)>>2) & 3)))

	muck1 := muck((seed>>8)&0xFF, seed, (seed>>4)&3)
	muck2 := muck((seed>>16)&0xFF, seed, (seed>>6)&3)

	fmt.Printf("Muck1 [%X]	Muck2[%X]\n\n", muck1, muck2)

	result := muck1 | (muck2 << 8)

	return result
}

// byte1 Function
//
//
func byte1(in int) uint {
	b1 := uint(uint8(in))
	fmt.Printf("BYTE1 [%X]\n", b1)
	result := b1 & 1
	fmt.Printf("BYTE1 Compare [%X]\n", result)

	return b1
}

// muck Function
//
//
func muck(a1 int, a2 int, a3 int) int {
	var result int
	fmt.Printf("a3: %v\n", a3)
	if a3 != 0 {
		switch a3 {
		case 1:
			//result = Abs(uint8(int(a1)*(a2-15)*int(a1)*(a2-15)*int(a1)*(a2-15)) - 7)
			result = Abs((a1 * (a2 - 15) * a1 * (a2 - 15) * a1 * (a2 - 15)) - 7)

		case 2:
			//result = Abs(uint8((int(a1)+2*a2-57)*(int(a1)+2*a2-57)) - uint8(a1))
			result = Abs(((a1 + 2*a2 - 57) * (a1 + 2*a2 - 57)) - a1)
			//sub_10076650(   (unsigned __int8) (     ( (_BYTE)a1 + 2 * a2 - 57 ) * ( (_BYTE)a1 + 2 * a2 - 57 )      ) - a1);

		case 3:
			//result = (3*a2 - int(a1) + 12) * (3*a2 - int(a1) + 12)
			result = (3*a2 - a1 + 12) * (3*a2 - a1 + 12)
		}
		return int(result)
	} else {
		//result = a1 + 13*(3*a1+2*a2)*13*(3*a1+2*a2)*13*(3*a1+2*a2)
		result = a1 + 13*(3*a1+2*a2)*13*(3*a1+2*a2)*13*(3*a1+2*a2)
	}

	return result
}

// Abs returns the absolute value of x.
//
//
func Abs(x int) int {
	if x < 0 {
		return int(-x)
	}
	return int(x)
}
