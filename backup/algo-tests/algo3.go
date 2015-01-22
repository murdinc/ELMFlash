package main

import (
	"fmt"
)

var seed1 = 0xD3FB8C // Key: AB D9 26?
var seed2 = 0xABEDCC // Key: 84 C4 EA?

var intseed1 = int(seed1)
var intseed2 = int(seed2)

func main() {

	key1 := getKey(intseed1)
	fmt.Printf("	Seed1: [%X] Key1: [%X] Expected Key: AB D9 26?\n", intseed1, key1)

	key2 := getKey(intseed2)
	fmt.Printf("	Seed2: [%X] Key2: [%X] Expected Key: 84 C4 EA?\n", intseed2, key2)

	fmt.Print("\n\n")

}

func getKey(seed int) uint8 {

	return uint8((seed ^ (seed >> 8) ^ 0x9B) + 43314)

}
