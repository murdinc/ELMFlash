package disasm

import (
	"fmt"
	//"github.com/murdinc/cli"
	"os"
	//"regexp"
)

// App constants
////////////////..........
const debug = false

func New() *DisAsm {
	controller := new(DisAsm)
	controller.location = "MSP.BIN"
	return controller
}

type DisAsm struct {
	location string
}

// Opcode struct
type OpCode struct {
	code        byte
	length      int
	mnemonic    string
	description string
	pseudo      string
	mode        string
}

func (h *DisAsm) Test() error {

	// Pull in the Calibration file
	f, err := os.Open("calibrations/MSP.BIN")
	if err != nil {
		log("BIN - Error opening file", err)
		return err
	}

	// Make a 1024 byte buffer
	block := make([]byte, 524288)

	// Read 1024 bytes
	n, err := f.Read(block)
	if err != nil {
		log("UploadBIN - Error reading calibration", err)
		return err
	}

	dbg(fmt.Sprintf("BIN - reading %d bytes.", n), nil)

	for i, b := range block {

		if contains(b, keys(mnemonics)) {
			log(fmt.Sprintf("Address:	%X	:	0x%.2X	length:	[%d]	mode: [%s]				[%s]", i, b, lengths[b], modes[b], mnemonics[b]), nil)
		}

	}
	return nil

}

func keys(m map[byte]string) (keys []byte) {
	for k := range m {
		keys = append(keys, k)
	}
	return keys
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
