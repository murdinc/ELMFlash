package main

import (
    "github.com/murdinc/ELMFlash/obd"
    "fmt"
    "strings"
    "io/ioutil"
)

//const baud = 9600
const baud = 38400


func main() {

    // Locate our ELM device
    device := findELM()

    // Connect to our device
    fmt.Print("Connecting...\n");
    elm , err := obd.New(device, baud)

    if err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {

        // Successful connection!
        fmt.Printf("Connected to device: %v\n", device)


        // Identify Device
        identity, err := elm.Identify()
        if err == nil {
            fmt.Printf("Identity: %v\n", identity)
        }

        // List Voltage
        voltage, err := elm.Voltage()
        if err == nil {
            fmt.Printf("Voltage: %v\n", voltage)
        } else {
            fmt.Printf("Error: %v\n", err)
        }


    }
}


func findELM() string {
    contents, _ := ioutil.ReadDir("/dev")

    // Look for what is mostly likely the Arduino device
    for _, f := range contents {
        if strings.Contains(f.Name(), "PL2303") || strings.Contains(f.Name(), "tty") {
            return "/dev/" + f.Name()
        }
    }

    // Have not been able to find a USB device that 'looks'
    // like an Arduino.
    return ""
}

