package obd

import (
	"bufio"
	"encoding/hex"
	"errors"
	"fmt"
	serial "github.com/huin/goserial"
	"io"
	"os"
	"strings"
)

const testerAddr = 0xF5
const ecuAddr = 0x10

// DEBUG MODE
const debug = false

// BaudDefault is the default baud rate to connect to a ELM327 OBD-II device
//const BaudDefault = 9600

// BaudFast is a faster baud rate available for some ELM327 OBD-II devices
//const BaudFast = 38400

// Constants
const EOL = 0x3E

var requestSeed = []byte{0x27, 0x01}
var sendKey = []byte{0x27, 0x02}
var ErrConnClosed = errors.New("obd: connection to device is closed")
var ErrIdentify = errors.New("obd: device is not a valid OBD-II device")
var ErrSetup = errors.New("obd: cannot configure OBD-II device for use")

// Connection represents a ELM327 OBD-II serial connection
type Connection struct {
	serial io.ReadWriteCloser
}

// Header
type Header struct {
	Length      byte
	Source      byte
	Destination byte
}

// Request
type Request struct {
	Header   Header
	Message  []byte
	Checksum byte
}

// Response
type Response struct {
	Header   Header
	Message  []byte
	Checksum byte
	Valid    bool
}

// New creates a new ELM327 OBD-II serial connection
func New(device string, baud int) (Connection, error) {
	// Create serial connection configuration
	fmt.Printf("Creating new device... [%v %v]\n", device, baud)
	config := &serial.Config{
		Name: device,
		Baud: baud,
	}

	var conn io.ReadWriteCloser

	if debug == true {
		fmt.Print("Debug is on!\n")
		return Connection{nil}, nil
	}
	// Attempt to open serial connection
	fmt.Print("Opening connection...\n")
	conn, err := serial.OpenPort(config)
	if err != nil {
		fmt.Printf("Connection error!: %v\n", err)
		return Connection{conn}, err
	}

	// Create OBD-II connection
	obd := Connection{conn}

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
	commands := []string{"AT D", "AT E0", "AT S0", "AT SP 3", "AT H1", "AT L0", "AT AL", "AT SI"}
	for _, c := range commands {
		// Send command, verify command received
		buf, err := obd.command(c)
		if err != nil {
			fmt.Print(err)
			return obd, ErrSetup
		}

		// Check for OK
		if !strings.Contains(string(buf), "OK") {
			return obd, errors.New("Error: \"" + string(buf) + "\"")
		}
	}

	if debug {
		return obd, nil
	} else {
		// Return new OBD-II connection
		return obd, nil
	}
}

// Close destroys connection to a ELM327 OBD-II device
func (c Connection) Close() error {
	// Reset device
	if err := c.Reset(); err != nil {
		return err
	}

	// Close connection
	return c.serial.Close()
}

// Reset closes the connection
func (c Connection) Reset() error {
	// AT Z - resets the device
	_, err := c.command("AT Z")
	return err
}

// Identify returns the identity of the current OBD-II device
func (c Connection) Identify() (string, error) {
	// AT I - Identify device
	buf, err := c.command("AT I")
	return string(buf), err
}

// Voltage returns the current battery voltage as reported by OBD-II device
func (c Connection) Voltage() (string, error) {
	// AT RV - Return battery voltage
	volt, err := c.command("AT RV")
	if err != nil {
		return "0.0V", err
	}

	return string(volt), nil
}

// Protocol
func (c Connection) Protocol() (string, error) {
	//
	proto, err := c.command("AT DP")
	if err != nil {
		return "0.0V", err
	}

	return string(proto), nil
}

// Speed returns the current vehicle speed as reported by OBD-II device
func (c Connection) Speed() (string, error) {
	// 01 0D - Return current vehicle speed
	speed, err := c.command("01 0D")
	if err != nil {
		return "0.0", err
	}

	// Convert from hex to decimal
	//return strconv.ParseInt("0x"+string(speed[4:6]), 0, 32)
	return string(speed), nil
}

// StartSecurity will enable mode 27 01 security
func (c Connection) StartSecurity() error {
	awake := false
	for !awake {
		fmt.Printf("Sending initial bit...\n")
		initialCommand := []byte{0xA0}
		resp, err := c.command(toString(initialCommand))
		if err == nil {
			obj := c.ResponseHandler(resp)
			fmt.Printf("ResponseHandler: [%X]\n\n", obj)
			if obj.Valid && (obj.Message[0] != 0xE0) {
				fmt.Print("Turn Ignition off...\n\n")
			} else if obj.Valid && (obj.Message[0] == 0xE0) {
				fmt.Print("FEPS detected!\n\n")
				break
			} else {
				fmt.Print("Turn ignition on...\n\n")
			}
		} else {
			return err
		}

	}

	// FEPS detected, now lets do this A0 thing
	//a0Thing := []byte{0x31, 0xA0, 0x02, 0x00, 0x67, 0x01}
	a0Thing := []byte{0x31, 0xA0, 0x02, 0x00, 0x4C, 0x01}
	resp, err := c.command(toString(a0Thing))
	if err == nil {
		obj := c.ResponseHandler(resp)
		fmt.Printf("ResponseHandler: [%X]\n\n", obj)
	}

	// A0 thing done, now lets get a security seed
	// This seems to tell the ECU which key and seed ot use *facepalm*
	getSeed := []byte{0x27, 0x01}
	resp, err = c.command(toString(getSeed))
	if err == nil {
		obj := c.ResponseHandler(resp)
		fmt.Printf("ResponseHandler: [%X]\n\n", obj)
	}
	return nil
}

func (c Connection) ResponseHandler(resp []byte) Response {
	strResp := string(resp)

	switch strResp {
	case "NO DATA":
		return Response{Valid: false}
	default:
		hex := toHex(resp)
		header := Header{Length: hex[0], Destination: hex[1], Source: hex[2]}
		message := hex[3:(len(hex) - 1)]
		checksum := hex[(len(hex) - 1)]
		formatResp := Response{Header: header, Message: message, Checksum: checksum, Valid: true}
		return formatResp
	}

}

// DumpBIN will read an entire bin in mode 23 (no auth)
func (c Connection) DumpBIN() error {
	// Open a file for writing
	f, err := os.Create("./dump1")
	if err != nil {
		fmt.Printf("ERROR opening file: %v\n\n", err)
	}

	defer f.Close()

	for i := 0xFF2080; i < 0xFFFFFF; i = i + 4 {
		addr := i // + 0x108000
		dumpCommand := []byte{0x23, byte(addr >> 16), byte(addr >> 8), byte(i)}
		resp, err := c.command(toString(dumpCommand))
		if err == nil {

			hexProp := toHex(resp[12:20])
			fmt.Printf("Response: %v\n\n", string(resp))
			fmt.Printf("Trimmed Response:%X\n\n", hexProp)

			n, err := f.Write(hexProp)
			if err != nil {
				fmt.Printf("ERROR writing file: %v\n\n", err)
			}
			fmt.Printf("wrote %d bytes\n", n)

			f.Sync()

		} else {
			return err
		}
	}

	return nil

}

// toHex converts an array of individual bytes into proper hex bytes
func toHex(in []byte) []byte {
	out := make([]byte, hex.DecodedLen(len(in)))
	_, err := hex.Decode(out, in)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return out
}

// iso_checksum generates a checksum value based on the input bytes
func iso_checksum(data []byte) byte {
	crc := byte(0x00)
	for i := 0; i < len(data); i++ {
		crc = crc + data[i]
	}
	return crc
}

// toString encodes our hex to a string so we can output it to the ELM 327
func toString(in []byte) string {
	return hex.EncodeToString(in)
}

// toBytes decodes our string from the ELM 327 into a byte array we can work with
func toByte(in string) []byte {
	byt, _ := hex.DecodeString(in)
	return byt
}

func (c Connection) setHeader(cmd []byte) error {
	// Header automation
	h1 := byte((len(cmd)+3)<<4) + 0x04 // length +1 for the checksum
	fullCmd := append([]byte{h1, ecuAddr, testerAddr}, cmd...)
	fullHeader := []byte{h1, ecuAddr, testerAddr}
	chks := iso_checksum(fullCmd) // just for the pretty factor
	fmt.Printf("Sending ECU: %X %2.2X %2.2X\n", fullHeader, cmd, chks)

	_, err := c.command("AT SH " + toString(fullHeader))
	if err != nil {
		return err
	}
	return nil
}

// command issues a command and retrieves a response from an OBD-II device
func (c Connection) command(cmd string) ([]byte, error) {

	// See if this is for the ELM or ECU
	if !strings.Contains(cmd, "AT") {
		err := c.setHeader(toByte(cmd))
		if err != nil {
			return []byte{}, err
		}
	}

	if debug {
		return []byte("DEBUG\n"), nil
	} else {
		// Check for open connection
		if c.serial == nil {
			fmt.Print("No open connections!")
			return nil, ErrConnClosed
		}

		// Issue command to device
		if _, err := c.serial.Write([]byte(cmd + "\r")); err != nil {
			return nil, err
		}

		// Read OBD-II response, loop until a response is generated
		reader := bufio.NewReader(c.serial)
		reply, err := reader.ReadBytes(EOL)
		reply = []byte(strings.Trim(string(reply[:]), "\r\n>"))
		if err != nil {
			return []byte{}, err
		}

		return reply, nil
	}
}
