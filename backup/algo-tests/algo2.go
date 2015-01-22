package main

import (
	"fmt"
)

//var seed1 = []byte{0xD3, 0xFB, 0x8C} // Key: AB D9 26?
//var seed2 = []byte{0xAB, 0xED, 0xCC} // Key: 84 C4 EA?

var seed1 = 0xD3FB8C
var seed2 = 0xABEDCC
var seed3 = 0x5BD112

var seed4 = 0x12D15B

func main() {

	muck := Muck(seed1, 2)
	fmt.Printf("SEED: [%X]MUCK: [%X] \n", seed1, muck)

	muck = Muck(seed2, 2)
	fmt.Printf("SEED: [%X]MUCK: [%X] \n", seed2, muck)

	muck = Muck(seed3, 2)
	fmt.Printf("SEED: [%X]MUCK: [%X] \n", seed3, muck)

	key1 := getKey(seed1)
	fmt.Printf("	Seed1: [%X] Key1: [%X] Expected Key: AB D9\n", seed1, key1)

	key2 := getKey(seed2)
	fmt.Printf("	Seed2: [%X] Key2: [%X] Expected Key: 84 C4\n", seed2, key2)

	key3 := getKey(seed3)
	fmt.Printf("	Seed3: [%X] Key3: [%b] Expected Key: 2A CB\n", seed3, key3)

	key4 := getKey(seed4)
	fmt.Printf("	Seed4: [%X] Key4: [%X] Expected Key: 2A CB\n", seed4, key4)

}

func getKey(seed int) int16 {
	//key = abs((signed __int16)((_WORD)seed * 2 * seed - 253 * seed));
	key := Abs(int16((seed * 2) * (seed - 253) * seed))

	return key
}

// Abs returns the absolute value of x.
func Abs(x int16) int16 {
	if x < 0 {
		return -x
	}
	return x
}

// The weird much func
func Muck(seed int, length int) int {
	var result int
	if length > 1 {
		for i := seed + length; ; i-- {
			fmt.Printf("i: [%X]\n", i)
			result = seed
			if seed >= i {
				break
			}
			v3 := byte(seed)
			seed = i
			i = int(v3)
		}
	}
	return result
}
