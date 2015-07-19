package iso9141

import (
	"bufio"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/cheggaaa/pb"
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

var calibrations = map[string]string{
	"msp": "MSP.BIN",
	"mp3": "MP3.BIN",
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
	DataAddr int
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

func (d *Device) EcuId() error {
	ecuIdCommand := []byte{0x10} // note: flipped most and least significant bytes
	idResp, err := d.Msg(ecuIdCommand)
	if err != nil {
		log("EcuId", err)
		return err
	} else {
		resp := fmt.Sprintf("ECU ID: %X", idResp.Message)
		log(resp, nil)
	}
	return nil
}

// TEST
func (d *Device) Test() error {
	return nil
}

// Reads an ECU memory in Security mode (faster)
func (d *Device) DownloadBIN(outfile string) error {
	// Make sure we have Security Access
	if d.SecurityMode == false {
		err := d.EnableSecurity()
		if err != nil {
			log("RunRoutine - Unable to enter secutiy mode!", err)
			return err
		}
	}

	// Open a file for writing
	ts := time.Now().Format(time.RFC3339)
	f, err := os.Create("./" + outfile + ts + ".BIN")

	defer f.Close()

	fmt.Println("./" + outfile + ts + ".BIN")
	if err != nil {
		log("DownloadBIN - Error opening file", err)
		return err
	}

	log("Starting Download...", nil)
	bar := pb.StartNew(480)

	for i := 0xFF0000; i < 0xFFFFFF; i = i + 0x0400 {
		bar.Increment()
		addr := i

		block, err := d.DownloadBlock(addr, 1024)
		if err != nil {
			return err
		}
		dbg(fmt.Sprintf("BLOCK: %X\n\n", block), nil)

		n, err := f.Write(block)
		if err != nil {
			log("DownloadBIN - Error writing to file", err)
			return err
		}
		dbg(fmt.Sprintf("DownloadBIN - wrote %d bytes", n), nil)

		f.Sync()
	}
	bar.FinishPrint("Download Finished!")
	return nil
}

// DumpBIN will read an entire bin in mode 23 (no auth)
func (d Device) DumpBIN(outfile string) error {
	// Open a file for writing
	ts := time.Now().Format(time.RFC3339)
	f, err := os.Create("./" + outfile + ts + ".BIN")

	defer f.Close()

	fmt.Println("./" + outfile + ts + ".BIN")
	if err != nil {
		log("DumpBIN - Error opening file", err)
		return err
	}

	log("Starting Dump...", nil)
	//bar := pb.StartNew(480)
	for i := 0x000000; i < 0x200000; i = i + 4 {
		addr := i
		dumpCommand := []byte{0x23, byte(addr >> 16), byte(addr >> 8), byte(i)}

		// Read some bytes
		log(fmt.Sprintf("Requesting Address: 0x%X \n", dumpCommand[1:]), nil)
		msgResp, err := d.Msg(dumpCommand)
		if err != nil {
			log(fmt.Sprintf("DumpBIN - Error accessing Address:  %X", dumpCommand[1:]), err)
			return err
		} else {
			contents := msgResp.Message[3:7]

			// Seek to new offset
			writeAddr64 := int64(i)
			writeOffset, err := f.Seek(writeAddr64, 0)

			n, err := f.Write(contents)
			if err != nil {
				log("DumpBIN - Error writing to file", err)
				return err
			}
			log(fmt.Sprintf("DumpBIN - Recieved: %X - wrote %d bytes to offset 0x%.2X", contents, n, writeOffset), nil)
		}

		f.Sync()
	}

	return nil
}

// Uploads a BIN file
func (d *Device) UploadBIN(calName string) error {

	// Pull in the Calibration file
	calFile := "./calibrations/" + calibrations[calName]
	log(fmt.Sprintf("Calibration File: %s", calFile), nil)

	f, err := os.Open(calFile)
	if err != nil {
		log("UploadBIN - Error opening file", err)
		return err
	}

	// Make sure we have Security Access
	if d.SecurityMode == false {
		err := d.EnableSecurity()
		if err != nil {
			log("RunRoutine - Unable to enter secutiy mode!", err)
			return err
		}
	}

	// Request by PID
	_, err = d.Msg([]byte{0x22, 0x11, 0x00})
	if err != nil {
		return err
	}

	// Read some bytes
	_, err = d.Msg([]byte{0x23, 0x10, 0x80, 0x80})
	if err != nil {
		return err
	}

	// Delete BIN on ECU
	err = d.RunRoutine([]byte{0x31, 0xA1}, []byte{0x32, 0xA1, 0x00}, []byte{0x22, 0x23})
	if err != nil {
		dbg("UploadBIN - Routine 31 A1", err)
	}

	// Make a 1024 byte buffer
	block := make([]byte, 1024)

	// Write
	count := 0
	log("Starting Upload...", nil)
	bar := pb.StartNew(480)
	for i := 0x007C00; count < 480; i = i - 0x0400 {
		writeOffset := 0x10FC00
		if count != 0 && count%32 == 0 {
			i = i + 0x10000
		}
		if count >= 96 {
			writeOffset = writeOffset + 0x80000
		}
		count++
		bar.Increment()

		// Seek to new offset
		readAddr64 := int64(i)
		readOffset, err := f.Seek(readAddr64, 0)

		// Read 1024 bytes
		n, err := f.Read(block)
		if err != nil {
			log("UploadBIN - Error reading calibration", err)
			return err
		}
		dbg(fmt.Sprintf("UploadBIN - reading %d bytes from offset: 0x%X", n, readOffset), nil)

		// Get the destination addresses
		writeAddr := writeOffset + int(readOffset)
		dbg(fmt.Sprintf("UploadBIN - writing %d bytes to offset 0x%X", n, writeAddr), nil)

		// Append the Checksum
		crc := uint16(0x0000)
		for i := 0; i < len(block); i++ {
			crc = crc + uint16(block[i])
		}

		chkh := byte(crc >> 8)
		chkl := byte(crc)

		dbg(fmt.Sprintf("UploadBIN - appending checksum: 0x%X 0x%X", chkh, chkl), nil)
		blockChk := append(block, chkh, chkl)

		// Upload the Calibration
		err = d.UploadBlock(writeAddr, 1024, blockChk)
		if err != nil {
			return err
		}
	}

	bar.FinishPrint("Upload Finished!")

	// Run Routine A3
	err = d.RunRoutine([]byte{0x31, 0xA3, 0x1F, 0x3F}, []byte{0x32, 0xA3, 0x00}, []byte{0x22, 0x23})
	if err != nil {
		dbg("UploadBIN - Routine A3 [FAIL] [", err)
	}

	return nil
}

func (d *Device) CommonIdDump(outfile string) error {
	// Open a file for writing
	ts := time.Now().Format(time.RFC3339)
	f, _ := os.Create("./" + outfile + ts + ".BIN")
	//f, err := os.Create(fmt.Sprintf("./%s-%s.%s.%s.%s.%s.%s", outfile, ts.Day(), ts.Month(), ts.Year(), ts.Hour(), ts.Minute(), ts.Second()))

	defer f.Close()

	for i := 0x0000; i < 0x0FFFF; i++ {
		//i1 := byte(i >> 16)
		i2 := byte(i >> 8)
		i3 := byte(i)
		//commonIdCommand := []byte{0x22, i3, i2, i1} // note: flipped most and least significant bytes
		commonIdCommand := []byte{0x22, i2, i3}
		log(fmt.Sprintf("Trying Command: %X \n", commonIdCommand), nil)
		msgResp, err := d.Msg(commonIdCommand)
		if err != nil {
			dbg(fmt.Sprintf("CMD %X", commonIdCommand), err)
		} else {
			resp := fmt.Sprintf("Common ID: %X Response: %X", commonIdCommand[1:], msgResp.Message[3:len(msgResp.Message)-1])
			fmt.Print(resp + "\n")
			n, err := f.WriteString(resp + "\n")
			if err != nil {
				log("CommonIdDump - Error writing to file", err)
				return err
			}
			dbg(fmt.Sprintf("CommonIdDump - %s - wrote %d bytes", resp, n), nil)
		}

		f.Sync()
	}
	return nil

}

func (d *Device) LocalIdDump(outfile string) error {
	if d.SecurityMode == false {
		err := d.EnableSecurity()
		if err != nil {
			log("DownloadBlock - Unable to enter secutiy mode!", err)
			return err
		}
	}

	// Open a file for writing
	ts := time.Now().Format(time.RFC3339)
	f, _ := os.Create("./" + outfile + ts + ".BIN")
	//f, err := os.Create(fmt.Sprintf("./%s-%s.%s.%s.%s.%s.%s", outfile, ts.Day(), ts.Month(), ts.Year(), ts.Hour(), ts.Minute(), ts.Second()))

	defer f.Close()

	for i := 0x00; i < 0x0FF; i++ {
		i1 := byte(i)
		localIdCommand := []byte{0x21, i1}
		log(fmt.Sprintf("Trying Command: %X \n", localIdCommand), nil)
		msgResp, err := d.Msg(localIdCommand)
		if err != nil {
			dbg(fmt.Sprintf("CMD %X", localIdCommand), err)
		} else {
			resp := fmt.Sprintf("Local ID: %X Response: %X", localIdCommand[1:], msgResp.Message[3:len(msgResp.Message)-1])
			fmt.Print(resp + "\n")
			n, err := f.WriteString(resp + "\n")
			if err != nil {
				log("LocalIdDump - Error writing to file", err)
				return err
			}
			dbg(fmt.Sprintf("LocalIdDump - %s - wrote %d bytes", resp, n), nil)
		}

		f.Sync()
	}
	return nil

}

func (d *Device) DownloadBlock(start, length int) ([]byte, error) {
	if d.SecurityMode == false {
		err := d.EnableSecurity()
		if err != nil {
			log("DownloadBlock - Unable to enter secutiy mode!", err)
			return nil, err
		}
	}

	l1 := byte(length >> 8)
	l2 := byte(length)

	s1 := byte(start >> 16)
	s2 := byte(start >> 8)
	s3 := byte(start)

	// [1] 35	= download by address command
	// [2] 82	= ?
	// [3-4]	= length (01 00,02 00,04 00 - only) 256, 512, 1024
	// [5-6-7]	= 00 00 00 address

	// Request Download Transfer
	dbg(fmt.Sprintf("Requesting Bytes: 0x%.6X - 0x%.6X", start, start+length-1), nil)
	downloadCommand := []byte{0x35, 0x82, l1, l2, s1, s2, s3}
	resp, err := d.Msg(downloadCommand)
	if err != nil {
		log("DownloadBlock [FAIL] [", err)
		return []byte{}, err
	}

	// Request Download Transfer Exit
	exitCommand := []byte{0x37, 0x82}
	_, err = d.Msg(exitCommand)
	if err != nil {
		log("DownloadBlock [FAIL] [", err)
		return []byte{}, err
	}

	// Trim the data to proper size
	resp.Data = resp.Data[:length]

	resp.DataAddr = start
	return resp.Data, nil
}

func (d *Device) UploadBlock(start, length int, block []byte) error {
	if d.SecurityMode == false {
		err := d.EnableSecurity()
		if err != nil {
			log("UploadBlock - Unable to enter secutiy mode!", err)
			return err
		}
	}

	l1 := byte(length >> 8)
	l2 := byte(length)

	s1 := byte(start >> 16)
	s2 := byte(start >> 8)
	s3 := byte(start)

	// [1] 34   = upload by address command
	// [2] 82   = ?
	// [3-4]    = length (01 00,02 00,04 00 - only) 256, 512, 1024
	// [5-6-7]  = 00 00 00 address

	// Request Upload Transfer
	dbg(fmt.Sprintf("Requesting to upload Bytes: 0x%.6X - 0x%.6X", start, start+length-1), nil)
	uploadCommand := []byte{0x34, 0x82, l1, l2, s1, s2, s3}
	_, err := d.Msg(uploadCommand)
	if err != nil {
		dbg("UploadBlock - Request Upload - 34 82 [FAIL] [", err)
		return err
	}

	// Set the timeout low
	resp := d.Send(Packet{Message: []byte("AT ST 01")})
	dbg(fmt.Sprintf("Timeout low: %X", resp.Message), nil)

	// Send the block
	for i := 0; i < length; i++ {
		if i%6 == 0 {
			end := i + 6
			if end >= length {
				end = len(block)
			}

			uploadBlock := append([]byte{0x36}, block[i:end]...)
			_, err = d.Msg(uploadBlock)
		}
	}

	// Set the timeout high
	resp = d.Send(Packet{Message: []byte("AT ST 00")})
	dbg(fmt.Sprintf("Timeout default: %X", resp.Message), nil)

	// Request Download/Upload Transfer Exit
	exitCommand := []byte{0x37, 0x82}
	_, err = d.Msg(exitCommand)
	if err != nil {
		dbg("UploadBlock - Request Transfer Exit - 37 82 [FAIL] [", err)
	}

	// Run Routine A2
	err = d.RunRoutine([]byte{0x31, 0xA2}, []byte{0x32, 0xA2, 0x00}, []byte{0x23})
	if err != nil {
		dbg("UploadBlock - Routine A2 [FAIL] [", err)
	}

	return nil
}

func (d *Device) RunRoutine(start, stop, success []byte) error {
	// Make sure we have Security Access
	if d.SecurityMode == false {
		err := d.EnableSecurity()
		if err != nil {
			log("RunRoutine - Unable to enter secutiy mode!", err)
			return err
		}
	}

	done := false

	// Start Routine
	for !done {
		resp, err := d.Msg(start)

		if resp.Error != nil {
			if contains(resp.ErrCode, success) {
				dbg("Start Routine [PASS]!", nil)
				break
			} else {
				dbg("Start Routine [FAIL] [", err)
			}
		}
	}

	// Stop Routine
	for !done {
		resp, _ := d.Msg(stop)

		errCode := resp.Message[(len(resp.Message) - 2)]

		if resp.Message[0] == errResp && errCode == 0x00 {
			msg := fmt.Sprintf("Stop Routine [PASS] - Response: %X", resp.Message)
			dbg(msg, nil)
			break
		} else {
			msg := fmt.Sprintf("Stop Routine [FAIL] - Response: %X", resp.Message)
			dbg(msg, nil)
		}
	}

	return nil
}

func (d *Device) EnableSecurity() error {

	// Pick a random security key, because why not?
	rand.Seed(time.Now().UTC().UnixNano())
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
			dbg("FEPS [PASS]", nil)
			break
		}
	}

	// Setup Security Algorithm
	a0 := []byte{0x31, 0xA0, 0x02, 0x00, algo.ID, 0x01}
	_, err := d.Msg(a0)
	if err != nil {
		dbg("EnableSecurity - Set Algo [FAIL] [", err)
	} else {
		dbg("EnableSecurity - Set Algo [PASS]", nil)
	}

	// Request Security Seed
	getSeed := []byte{0x27, 0x01}
	_, err = d.Msg(getSeed)
	if err != nil {
		dbg("EnableSecurity - Request seed FAIL", err)
	} else {
		dbg("EnableSecurity - Request Seed [PASS]", nil)
	}

	// Submit Security Key
	submitKey := append([]byte{0x27, 0x02}, algo.Key...)
	_, err = d.Msg(submitKey)
	if err != nil {
		dbg("EnableSecurity - Submit Key] [FAIL", err)
	} else {
		dbg("EnableSecurity - Submit Key [PASS]", nil)
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
	dbg("Received]: ["+string(reply), nil)

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
		log("ConnectDevice - [FAIL", err)
		os.Exit(1)
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
	commands := []string{"AT D", "AT E0", "AT S0", "AT SP 3", "AT H1", "AT L0", "AT AL", "AT SI", "AT CAF0", "AT AT1"}
	for _, c := range commands {
		pkt := Packet{Message: []byte(c)}
		resp := d.Send(pkt)
		if resp.Error != nil {
			dbg("Setup Command Failure: "+c, nil)
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

func contains(n byte, h []byte) bool {
	for _, c := range h {
		if c == n {
			return true
		}
	}
	return false
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
