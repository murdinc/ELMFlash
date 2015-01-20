package main

import (
	"bufio"
	"encoding/hex"
	"errors"
	"fmt"
	serial "github.com/huin/goserial"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

// App constants
////////////////..........
const baud = 115200
const debug = false
const headless = false
const obdDevice = "STY3M"
const EOL = 0x3E
const testerAddr = 0xF5
const ecuAddr = 0x10
const errResp = 0x7F

var errCodes = []string{
	0x10: "10 - General Reject",
	0x11: "11 - Service Not Supported",
	0x12: "12 - Sub Function Not Supported - Invalid Format",
	0x21: "21 - Busy - repeat Request",
	0x22: "22 - Conditions Not Correct Or Request Sequence Error",
	0x23: "23 - Routine Not Complete Or Service In Progress",
	0x31: "31 - Request Out Of Range",
	0x33: "33 - Security Access Denied - security Access Requested",
	0x35: "35 - Invalid Key",
	0x36: "36 - Exceed Number Of Attempts",
	0x37: "37 - Required Time Delay Not Expired",
	0x40: "40 - Download Not Accepted",
	0x41: "41 - Improper Download Type",
	0x42: "42 - Can Not Download To Specified Address",
	0x43: "43 - Can Not Download Number Of Bytes Requested",
	0x50: "50 - Upload Not Accepted",
	0x51: "51 - Improper Upload Type",
	0x52: "52 - Can Not Upload From Specified Address",
	0x53: "53 - Can Not Upload Number Of Bytes Requested",
	0x71: "71 - Transfer Suspended",
	0x72: "72 - Transfer Aborted",
	0x74: "74 - Illegal Address In Block Transfer",
	0x75: "75 - Illegal Byte Count In Block Transfer",
	0x76: "76 - Illegal Block Trasnfer Type",
	0x77: "77 - Block Transfer Data Checksum Error",
	0x78: "78 - Request Correcty Rcvd - Rsp Pending",
	0x79: "79 - Incorrect Byte Count During Block Transfer",
	0x80: "80 - Service Not Supported In Active Diagnostic Mode",
	0xC1: "C1 - Start Comms +ve response",
	0xC2: "C2 - Stop Comms +ve response",
	0xC3: "C3- Access Timing Params +ve response",
}

// Main Function
////////////////..........
func main() {
	obd := New()

	cmdResp, err := obd.Cmd("ATDP")
	if err != nil {
		log("ATDP", err)
	} else {
		log("Protocol: "+cmdResp, nil)
		//fmt.Printf("Test Command response: %v\n", cmdResp)
	}

	secMode := []byte{0xA0}
	msgResp, err := obd.Msg(secMode)
	if err != nil {
		log("CMD A0", err)
	} else {
		fmt.Printf("Test Message response: %X\n", msgResp)
	}
}

// OBD Types
////////////////..........
type Packet struct {
	Header   []byte
	Message  []byte
	Checksum byte
	Prepared bool
	Error    error
}

// Connection represents an OBD-II serial connection
type Device struct {
	Packet     Packet
	serial     io.ReadWriteCloser
	location   string
	baud       int
	lastHeader []byte
}

// Device Functions
////////////////..........
func (d Device) Send(packet Packet) Packet {

	// Check for open connection
	if d.serial == nil {
		dbg("No serial connection!", nil)
		return Packet{}
	}

	var send string

	// Issue command to device
	if packet.Prepared == true {
		send = string(packet.Message)
		dbg("Sending]: ["+send, nil)
	} else {
		send = string(packet.Message)
		dbg("Sending]: ["+send, nil)
	}

	_, err := d.serial.Write(append(packet.Message, []byte("\r")...))
	if err != nil {
		dbg("Error sending packet to serial device!", nil)
	}

	// Wait for our reply
	return d.Recieve()
}

func (d Device) Recieve() Packet {

	// Read OBD-II response, loop until a response is generated
	reader := bufio.NewReader(d.serial)
	reply, err := reader.ReadBytes(EOL)
	reply = []byte(strings.Trim(string(reply[:]), "\r\n>"))
	dbg("Recieved]: ["+string(reply), nil)
	if (err != nil) || (string(reply) == "?") {
		errMsg := errors.New("Unknown command")
		return Packet{Error: errMsg}
	}
	response := Packet{Message: reply}
	return response
}

func New() *Device {
	device := new(Device)
	device.FindDevice()
	device.ConnectDevice()
	return device
}

func (d Device) Cmd(cmd string) (string, error) {
	command := Packet{Message: []byte(cmd)}
	resp := d.Send(command)

	strResp := string(resp.Message)

	// Check for OK
	if strings.Contains(string(strResp), "?") {
		resp.Error = errors.New("Error sending command: [" + cmd + "] response: [" + strResp + "]")
	}

	return strResp, resp.Error
}

func (d *Device) Msg(msg []byte) (Packet, error) {
	str := toString(msg)
	msg = []byte(str)
	message := Packet{Message: msg}

	// Handle header
	message.prepare()
	if len(d.lastHeader) > 0 && message.Header[0] == d.lastHeader[0] {
		dbg("Prepare]: [Header already correctly set ", nil)
	} else {
		headerMsg := Packet{Message: append([]byte("AT SH"), message.Header...)}
		d.Send(headerMsg)
		d.lastHeader = message.Header
	}

	// Format response into Packet type
	resp := d.Send(message)
	hex := toHex(resp.Message)
	resp.Header = hex[0:3]
	resp.Message = hex[3:(len(hex) - 1)]
	resp.Checksum = hex[(len(hex) - 1)]

	if resp.Message[0] == errResp {
		errCode := hex[(len(hex) - 2)]
		resp.Error = errors.New("Recieved error from ECU: " + errCodes[errCode])
	}

	return resp, resp.Error
}

func (p *Packet) prepare() {
	// Message
	p.Message = []byte(toHex(p.Message))

	// Header
	h1 := byte((len(p.Message)+3)<<4) + 0x04 // length +1 for the checksum
	header := []byte{h1, ecuAddr, testerAddr}
	p.Header = []byte(toString(header))

	// Checksum
	crc := byte(0x00)
	for i := 0; i < len(p.Message); i++ {
		crc = crc + p.Message[i]
	}

	p.Message = []byte(toString(p.Message))

	p.Prepared = true
	p.Checksum = crc
}

func toString(in []byte) string {
	return hex.EncodeToString(in)
}

func toHex(in []byte) []byte {
	out := make([]byte, hex.DecodedLen(len(in)))
	_, err := hex.Decode(out, in)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return out
}

func (d *Device) ConnectDevice() {

	if len(d.location) < 1 {
		d.FindDevice()
	}

	dbg("Setting up connection to device: "+d.location, nil)
	//fmt.Printf("Creating new device... [%v %v]\n", d.location, d.baud)
	config := &serial.Config{
		Name: d.location,
		Baud: d.baud,
	}

	var conn io.ReadWriteCloser

	// Attempt to open serial connection
	dbg("Opening serial connection to device: "+d.location, nil)
	conn, err := serial.OpenPort(config)
	if err != nil {
		dbg("Open Connection", err)
	}

	// Create OBD-II connection
	d.serial = conn

	// AT D - Sets All Defaults
	// AT E0 - Disable device echo
	// AT L0 - Disable line feed
	// AT S0 - Disable spaces
	// AT AT2 - Enable faster responses
	// AT SP 00 - Automatically select protocol
	// AT H1 - Turns on headers
	// AT L1 - Enables line feeds
	// AT CA F1 - CAN Automatic Formatting on
	// AT AL - Allow Long Messages

	// Run set of commands to properly setup our communication with the car
	commands := []string{"AT D", "AT E0", "AT S0", "AT SP 3", "AT H1", "AT L0", "AT AL", "AT SI"}
	for _, c := range commands {
		pkt := Packet{Message: []byte(c)}
		resp := d.Send(pkt)
		if resp.Error != nil {
			log("Setup Command Failure: "+c, nil)
		}
	}
}

func (d *Device) FindDevice() bool {
	contents, _ := ioutil.ReadDir("/dev")

	// Look for what is mostly likely the Arduino device
	for _, f := range contents {
		if strings.Contains(f.Name(), "STY3M") && strings.Contains(f.Name(), "tty") {
			d.location = "/dev/" + f.Name()
			d.baud = baud
			dbg("Found Device: "+d.location, nil)
			return true
		}
	}
	return false
}

func (d *Device) DisconnectDevice() {

}

// Debug Function
////////////////..........
func dbg(kind string, err error) {
	if debug {
		if err == nil {
			fmt.Printf("### [DEBUG LOG - %s]\n\n", kind)
		} else {
			fmt.Printf("### [DEBUG ERROR - %s]: %s\n\n", kind, err)
		}
	}
}

func log(kind string, err error) {
	if err == nil {
		fmt.Printf("[LOG - %s]\n", kind)
	} else {
		fmt.Printf("[ERROR - %s]: %s\n", kind, err)
	}
}
