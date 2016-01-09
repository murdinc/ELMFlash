package j3

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/tarm/serial"
)

//var led uint8 = 13

// App constants
////////////////..........
const baud = 115200
const debug = false
const EOL = 0x3E

// SDU PINS
/*
var led uint8 = 13   // LED
var crdclk uint8 = 5 // mcu input, arduino output
var crin uint8 = 6   // mcu input, arduino output
var crout uint8 = 7  // mcu output, arduino input
var crbusy uint8 = 8 // mcu output, arduino input
*/

type J3Port struct {
	serial   io.ReadWriteCloser
	location string
	baud     int
	Dummy    bool
}

type Packet struct {
	/*
		Clock   byte
		In      byte
		Out     byte
		Busy    byte
	*/
	Message []byte
	Error   error
}

type CommandByte struct {
	Packets []Packet
}

func Test() {
	j3 := New(false)

	// Set the pin mode
	// 010xxxxx - Configure pins as input(1) or output(0)
	// AUX|MOSI|CLK|MISO|CS
	// 0|1|0|A|O|C|I|C - abbreviations
	// 0|1|0|0|O|O|I|I - direction
	// 0|1|0|0|0|0|1|1 - set in this call (0x43)
	resp := j3.Send(Packet{Message: []byte{0x43}})
	if (resp.Message[0]|byte(0x40) != 0x40) && resp.Message[0] != 0x42 { // TODO: reset from: BBIO1BBIO
		log("Error setting pin directions!", nil)
		fmt.Printf("MESSAGE: %s", resp.Message)
		//os.Exit(1)
	}

	// Shift in SDU command 06H to select the code RAM address access instruction
	j3.SendCommand(0x06)

	// Shift the high byte (04H) of the address, followed by the low byte (80H) of the address into the CR_ADDR register
	j3.SendCommand(0x04)
	j3.SendCommand(0x80)

	// Shift in SDU command 10H to select the code RAM word read instruction.
	j3.SendCommand(0x80)

	// Wait for CRBUSY# to be deasserted, indicating that the data is ready.

	// Shift in SDU command 10H to read the next word.
	j3.SendCommand(0x10)

	// Shift in SDU command 28H to put the SDU into an idle state. This shift causes the low byte of data0 to be shifted out on CROUT.
	j3.SendCommand(0x28)

	// Shift in SDU command 00H to reset the SDU and terminate the data transfer. This shift causes the high byte of data1 to be shifted out on CROUT.
	j3.SendCommand(0x00)

	// Shift in SDU command 28H to put the SDU into an idle state. This shift causes the low byte of data1 to be shifted out on CROUT.
	j3.SendCommand(0x28)

	/*
		for {
			// Running at 31.446 Hz?

			pkt.Message = []byte{0xFF}
			resp = j3.Send(pkt)
			//fmt.Printf("GOT: %X", resp.Message)

			pkt.Message = []byte{0x80}
			resp = j3.Send(pkt)
			//fmt.Printf("GOT: %X", resp.Message)
		}
	*/
}

// 1xxxxxxx - Set on (1) or off (0):
// 1|POWER|PULLUP|AUX|MOSI|CLK|MISO|CS
// 1|0|0|0|O|C|I|B - abbreviations
// 1|0|0|0|O|O|I|I - direction
// 1|0|0|0|*|0|0|0 - set in this function
func (j *J3Port) SendCommand(command byte) {
	clock := false
	for i := 8; i > 0; i-- {
		if clock == false {
			clock = true
			i++
			j.Send(Packet{Message: []byte{0x80}})

		} else {
			clock = false
			bitno := uint(i) - 1
			bit := (command >> bitno) & 0x01

			// clock + = 0x88
			// clock - = 0x80

			bb := ((bit << 3) | 0x84)
			//fmt.Printf("bit: 0x%.2X\n", bit)
			//fmt.Printf("bb: 0x%.2X\n\n", bb)

			resp := j.Send(Packet{Message: []byte{bb}})

			fmt.Printf("\nSent: %X\n", bb)
			fmt.Printf("Received: %X\n", resp.Message)

			if resp.Message[0] != bb {
				fmt.Println("DIFFERENCE")
			}

		}

	}

}

func New(test bool) *J3Port {

	if test {
		j3 := new(J3Port)
		j3.Dummy = true
		return j3
	}

	j3 := new(J3Port)
	//device.FindDevice()
	j3.location = "/dev/tty.usbserial-AL004BZR"
	j3.baud = baud
	j3.ConnectDevice()
	return j3
}

func (j *J3Port) ConnectDevice() {

	if len(j.location) < 1 {
		//d.FindDevice()
	}

	dbg("Setting up connection to device: "+j.location, nil)
	//fmt.Printf("Creating new device... [%v %v]\n", d.location, d.baud)
	config := &serial.Config{
		Name: j.location,
		Baud: j.baud,
	}

	var conn io.ReadWriteCloser

	// Attempt to open serial connection
	dbg("Opening serial connection to device: "+j.location, nil)
	conn, err := serial.OpenPort(config)
	if err != nil {
		log("ConnectDevice - [FAIL", err)
		os.Exit(1)
	}

	// Add serial connection
	j.serial = conn

	reset := []byte{0x0F, 0x35}
	j.Send(Packet{Message: reset})

	// enter bitbang mode
	for i := 0; i < 26; i++ {
		pkt := Packet{Message: []byte{0x00}}
		resp := j.Send(pkt)
		if resp.Error != nil {
			dbg("Setup Command Failure: 0x00", nil)
			log("Unable to enter bitbang mode", nil)
			os.Exit(1)
		}
		if bytes.Equal(resp.Message, []byte("BBIO1")) {
			log("Successfully Entered bitbang mode!", nil)
			return
		}
	}
}

func (j J3Port) Send(packet Packet) Packet {

	// Check for open connection
	if j.serial == nil {
		dbg("No serial connection!", nil)
		return Packet{}
	}

	dbg(fmt.Sprintf("Sending]: [%X", packet.Message), nil)

	_, err := j.serial.Write(append(packet.Message, []byte("\r")...))
	if err != nil {
		dbg("Error sending packet to serial device!", nil)
	}

	// Wait for our reply
	return j.Receive()
}

func (j J3Port) Receive() Packet {

	reader := bufio.NewReader(j.serial)

	reply := make([]byte, 255)
	n, err := reader.Read(reply)
	reply = reply[0 : n-1]

	reply = []byte(strings.Trim(string(reply[:]), "\r\n>"))
	dbg("Received]: [ "+string(reply), nil)

	if err != nil {
		errMsg := errors.New("Error reading reply")
		return Packet{Error: errMsg}
	}
	resp := Packet{Message: reply}

	/*
		strResp := string(resp.Message)
			// Check for ERROR
			if (strings.Contains(strResp, "?")) || (strings.Contains(strResp, "ERROR")) || (strings.Contains(strResp, "NO DATA")) {
				resp.Error = errors.New("Response: [" + strResp + "]")
				resp.Message = nil
				resp.ErrCode = 0xFF
			}
	*/
	return resp
}

// Debug Function
////////////////..........
func dbg(kind string, err error) {
	if debug {
		if err == nil {
			fmt.Printf("[ %s ]\n", kind)
		} else {
			fmt.Printf("### [DEBUG ERROR - %s]: %s\n\n", kind, err)
		}
	}
}

func log(kind string, err error) {
	if err == nil {
		fmt.Printf("====> %s\n", kind)
	} else {
		fmt.Printf("[ERROR - %s]: %s\n", kind, err)
	}
}
