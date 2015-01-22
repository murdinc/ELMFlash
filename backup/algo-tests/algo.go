package main

import (
	"fmt"
	//"strings"
	//"io/ioutil"
)

var seed1 = []byte{0xD3, 0xFB, 0x8C} // Key: AB D9 26?
var seed2 = []byte{0xAB, 0xED, 0xCC} // Key: 84 C4 EA?
var seed3 = []byte{0x03, 0xD9, 0x77} // Key: BE EA 61
var challenge = []byte{0, 0, 0, 'M', 'a', 'z', 'd', 'A'}

const initialValue int64 = 0xC541A9

const v1 int64 = 0x109028
const v2 int64 = 0xFFEF6FD7

func main() {

	fmt.Print("The first two seeds are known for the Protege platform, the third seed and key is known for the Mazda3 platform and is what this algo is known to be working for.\n\n")

	fmt.Print("At the very least, the third Seed/Key should match what is expected. \n\n")

	/*
		key1 := getKey(seed1)
		fmt.Printf("	Seed1: [%X] Key1: [%X] Expected Key: AB D9 26?\n", seed1, key1)

		key2 := getKey(seed2)
		fmt.Printf("	Seed2: [%X] Key2: [%X] Expected Key: 84 C4 EA?\n", seed2, key2)
	*/
	key3 := getKey(seed3)
	fmt.Printf("	Seed3: [%X] Key3: [%X] Expected Key: BE EA 61\n", seed3, key3)

	fmt.Print("\n\n")

}

func getKey(seed []byte) []byte {

	challenge[0] = seed[0]
	challenge[1] = seed[1]
	challenge[2] = seed[2]

	fmt.Printf("\nChallenge: [%X] \n", challenge)

	buffer := initialValue

	for i := 0; i < len(challenge); i++ {
		b := challenge[i] & 0xff

		for j := 0; j < 8; j++ {
			tempBuffer := int64(0)
			if (b & 1) != (byte(buffer) & 1) {
				buffer |= 0x1000000
				tempBuffer = v1
			}

			b = TripleShift(int64(b), 1)

			tempBuffer ^= buffer >> 1
			tempBuffer &= v1
			tempBuffer |= v2 & (buffer >> 1)

			buffer = tempBuffer & 0xFFFFFF
		}
	}

	fmt.Printf("BUFFER: [%X]\n\n", buffer)

	key := []byte{
		(TripleShift(int64(buffer), 4) & 0xFF),
		(TripleShift(buffer, 20) & 0x0F) + (TripleShift(buffer, 8) & 0xF0),
		((byte(buffer) << 4) & 0xFF) + (TripleShift(buffer, 16) & 0x0F),
	}

	return key
}

func TripleShift(n int64, s uint) byte {
	//fmt.Printf("Triple Shift n: [%X] s: [%X]\n", n, s)
	if n >= 0 {
		return byte(n) >> s
	}
	fmt.Print("Less than zero\n")

	return (byte(n) >> s) + (2 << s)
}
