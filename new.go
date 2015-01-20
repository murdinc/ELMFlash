package main

import (
	"bufio"
	"encoding/hex"
	"errors"
	"fmt"
	serial "github.com/huin/goserial"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"
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
	0x76: "76 - Illegal Block Transfer Type",
	0x77: "77 - Block Transfer Data Checksum Error",
	0x78: "78 - Request Correcty Rcvd - Rsp Pending",
	0x79: "79 - Incorrect Byte Count During Block Transfer",
	0x80: "80 - Service Not Supported In Active Diagnostic Mode",
	0xC1: "C1 - Start Comms +ve response",
	0xC2: "C2 - Stop Comms +ve response",
	0xC3: "C3 - Access Timing Params +ve response",
	0xFF: "FF - No Data",
}

var Algos = []SecAlgo{
	SecAlgo{ID: 0x4C, Seed: []byte{0xAB, 0xED, 0xCC}, Key: []byte{0x84, 0xC4}},
	SecAlgo{ID: 0x67, Seed: []byte{0xD3, 0xFB, 0x8C}, Key: []byte{0xAB, 0xD9}},
	// TODO
}

type SecAlgo struct {
	ID   byte
	Seed []byte
	Key  []byte
}

// Main Function
////////////////..........
func main() {

	rand.Seed(time.Now().UTC().UnixNano())

	obd := New()

	obd.EnableSecurity()

	cmdResp, err := obd.Cmd("ATDP")
	if err != nil {
		log("ATDP", err)
	} else {
		log("Protocol: - ["+cmdResp+"]", nil)
	}

	obd.DownloadBlock()

	/*
		secMode := []byte{0xA0}
		msgResp, err := obd.Msg(secMode)
		if err != nil {
			log("CMD A0", err)
		} else {
			fmt.Printf("Test Message response: %X\n", msgResp)
		}
	*/
}

// OBD Types
////////////////..........
type Packet struct {
	Header   []byte
	Message  []byte
	Checksum byte
	Prepared bool
	Error    error
	ErrCode  byte
	Multi    []Packet
	Data     []byte
}

// Connection represents an OBD-II serial connection
type Device struct {
	Packet       Packet
	serial       io.ReadWriteCloser
	location     string
	baud         int
	lastHeader   []byte
	SecurityMode bool
}

// Device Functions
////////////////..........

func (d *Device) DownloadBlock() error {
	if d.SecurityMode == false {
		err := d.EnableSecurity()
		if err != nil {
			log("Unable to enter secutiy mode!", err)
		}
	}

	// [1] 35 = download
	// [2] 82 =?
	// [3]04 - length? (1,2,3,4 only) 256, 512, ? , 1024
	// [4] 00 = ?
	// [5-6] 00,00 address?
	// [7] 00=? // padding?
	downloadCommand := []byte{0x35, 0x82, 0x04, 0x00, 0x10, 0xE8, 0x00}
	resp, err := d.Msg(downloadCommand)
	if err != nil {
		log("DownloadBIN [FAIL] [", err)
	} else {
		fmt.Printf("DOWNLOAD DATA: %X LEN: %v\n\n", resp.Data, resp.DataLen())
	}

	return nil
}

func (d *Device) EnableSecurity() error {

	// Pick a random security key, because why not?
	algo := Algos[rand.Intn(len(Algos))]

	awake := false
	for !awake {
		initialCommand := []byte{0xA0}
		resp, _ := d.Msg(initialCommand)
		if resp.Error != nil {
			if resp.ErrCode != 0xFF && (resp.Message[0] != 0xE0) {
				log("Turn Ignition off...", nil)
			} else {
				log("Turn Ignition on...", nil)
			}
		} else if resp.Error == nil && resp.Message[0] == 0xE0 {
			log("FEPS [PASS]", nil)
			break
		}
	}

	// Setup Security Algorithm
	a0 := []byte{0x31, 0xA0, 0x02, 0x00, algo.ID, 0x01}
	_, err := d.Msg(a0)
	if err != nil {
		log("EnableSecurity - Set Algo [FAIL] [", err)
	} else {
		log("EnableSecurity - Set Algo [PASS]", nil)
	}

	// Request Security Seed
	getSeed := []byte{0x27, 0x01}
	_, err = d.Msg(getSeed)
	if err != nil {
		log("EnableSecurity - Request seed FAIL", err)
	} else {
		log("EnableSecurity - Request Seed [PASS]", nil)
	}

	// Submit Security Key
	submitKey := append([]byte{0x27, 0x02}, algo.Key...)
	_, err = d.Msg(submitKey)
	if err != nil {
		log("EnableSecurity - Submit Key] [FAIL", err)
	} else {
		log("EnableSecurity - Submit Key [PASS]", nil)
		d.SecurityMode = true
	}

	return nil
}

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
	return d.Receive()
}

func (d Device) Receive() Packet {

	// Read OBD-II response, loop until a response is generated
	reader := bufio.NewReader(d.serial)
	reply, err := reader.ReadBytes(EOL)
	reply = []byte(strings.Trim(string(reply[:]), "\r\n>"))
	dbg("Recieved]: ["+string(reply), nil)

	reply = []byte(strings.TrimSuffix(string(reply[:]), "<DATA ERROR"))

	if (err != nil) || (string(reply) == "?") {
		errMsg := errors.New("Unknown command")
		return Packet{Error: errMsg}
	}
	resp := Packet{Message: reply}

	strResp := string(resp.Message)
	// Check for ERROR
	if (strings.Contains(strResp, "?")) || (strings.Contains(strResp, "ERROR")) || (strings.Contains(strResp, "NO DATA")) {
		resp.Error = errors.New("Response: [" + strResp + "]")
		resp.Message = nil
		resp.ErrCode = 0xFF
	}
	return resp
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

	// Check for ERROR
	if (strings.Contains(strResp, "?")) || (strings.Contains(strResp, "ERROR")) {
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
		dbg("Prepare]: [Header already correctly set", nil)
	} else {
		headerMsg := Packet{Message: append([]byte("AT SH"), message.Header...)}
		d.Send(headerMsg)
		d.lastHeader = message.Header
	}

	// Send the message
	resp := d.Send(message)

	// Check if we have a ELM error already
	if resp.Error != nil {
		return resp, resp.Error
	}

	// Organize
	hex := toHex(resp.Message)
	resp.Header = hex[0:3]
	length := int(hex[0]>>4) + 1
	resp.unPack(hex[length:])
	resp.Message = hex[3:length]
	resp.Checksum = hex[(len(hex) - 1)]

	// Detect errors
	errCode := hex[(len(hex) - 2)]
	if resp.Message[0] == errResp && errCode != 0x00 {
		resp.Error = errors.New("Recieved error from ECU: " + errCodes[errCode])
		resp.ErrCode = errCode
	}

	return resp, resp.Error
}

func (p *Packet) unPack(in []byte) {
	var unpacked []Packet
	var data []byte

	for start := 0; start < len(in)-1; {
		// Find a single packed packet
		length := int(in[start]>>4) + 1
		end := start + length
		single := in[start:end]

		// Chop chop
		var packet Packet
		packet.Header = single[0:3]
		packet.Message = single[3:(len(single) - 1)]
		packet.Checksum = single[(len(single) - 1)]

		data = append(data, single[4:(len(single)-1)]...)
		unpacked = append(unpacked, packet)
		start = end
	}

	p.Data = data
	p.Multi = unpacked
}

func (p *Packet) DataLen() int {
	return len(p.Data)
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
		os.Exit(2)
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
	commands := []string{"AT D", "AT E0", "AT S0", "AT SP 3", "AT H1", "AT L0", "AT AL", "AT SI", "AT CAF0"}
	for _, c := range commands {
		pkt := Packet{Message: []byte(c)}
		resp := d.Send(pkt)
		if resp.Error != nil {
			log("Setup Command Failure: "+c, nil)
			log("Try turning the ignition to position 0 and then position 1 again.", nil)
			os.Exit(1)
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
		fmt.Printf("====> %s\n", kind)
	} else {
		fmt.Printf("[ERROR - %s]: %s\n", kind, err)
	}
}
