package main

import (
    "fmt"
    "encoding/hex"
)


func main() {
    // 75 F5 23 10 80 90 BC
    test := "7410F523108090"
    test2 := []byte(test)

    out, err := hex.DecodeString(test)

    if err != nil {
        fmt.Printf("ERROR: %v", err)
    }

    fmt.Printf("Decode Out: [%X]\n", out)
    fmt.Printf("Byte Out: [%v]", test2)
}

