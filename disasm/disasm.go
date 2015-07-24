package disasm

import (
	"fmt"
	"os"
	"strings"
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

var calibrations = map[string]string{
	"msp": "MSP.BIN",
	"mp3": "MP3.BIN",
	"pre": "PRE.BIN",
}

func (h *DisAsm) DisAsm(calName string) error {

	// Pull in the stuff before the calibration file
	preCalFile := "./calibrations/" + calibrations["pre"]
	log(fmt.Sprintf("Disassemble - Pre-calibration File: %s", preCalFile), nil)

	p, err := os.Open(preCalFile)
	pi, err := p.Stat()
	preFileSize := pi.Size()

	// Pull in the Calibration file
	calFile := "./calibrations/" + calibrations[calName]
	log(fmt.Sprintf("Disassemble - Calibration File: %s", calFile), nil)

	f, err := os.Open(calFile)
	fi, err := f.Stat()
	fileSize := fi.Size()
	if err != nil {
		log("Disassemble - Error opening file", err)
		return err
	}

	log(fmt.Sprintf("Disassemble - [%s] is %d bytes long", calibrations["pre"], preFileSize), nil)
	log(fmt.Sprintf("Disassemble - [%s] is %d bytes long", calibrations[calName], fileSize), nil)

	// Make some buffers
	preBlock := make([]byte, 0x108000)
	calBlock := make([]byte, 0x78000)

	// Read in all the bytes
	n, err := p.Read(preBlock)
	if err != nil {
		log("Disassemble - Error reading calibration", err)
		return err
	}
	log(fmt.Sprintf("Disassemble - reading 0x%X bytes from pre-calibration file.", n), nil)

	n, err = f.Read(calBlock)
	if err != nil {
		log("Disassemble - Error reading calibration", err)
		return err
	}

	log(fmt.Sprintf("Disassemble - reading 0x%X bytes from calibration file.", n), nil)

	block := append(preBlock, calBlock...)

	// Doubletime
	//block = append(block, block[0x100000:0x180000]...)

	log(fmt.Sprintf("Length: 0x%X", len(block)), nil)

	opSize := 1
	count := 1
	for i := 0x100000; i < len(block); i = i + opSize {

		if i > 0x108000 && i < 0x11FFFF {
			continue
		}

		b := block[i:]
		instr, err := Parse(b)

		if err != nil {
			log("", err)
		} else if instr.Ignore == false {

			address := addSpaces(fmt.Sprintf("Address: [0x%X]", i), 20)
			length := addSpaces(fmt.Sprintf(" Length: [%d]", instr.ByteLength), 14)
			mode := addSpaces(fmt.Sprintf(" Mode: [%s]", instr.AddressingMode), 26)
			mnemonic := addSpaces(fmt.Sprintf("	Mnemonic: [%s]", instr.Mnemonic), 23)
			shortDesc := addSpaces(fmt.Sprintf("%s", instr.Description), 10)
			operandCount := addSpaces(fmt.Sprintf(" [%d] Operands", instr.VarCount), 23)
			raw := addSpaces(fmt.Sprintf(" Raw: 0x%.10X", instr.Raw), 20)

			count++
			log("---------", nil)

			var l1, l2, l3 string

			l1 += addSpaces("", 10)
			l2 += addSpaces("", 10)
			l3 += addSpaces(instr.Mnemonic, 10)

			if !instr.Checked {
				log("#### ERROR DISASEMBLING OPCODE ####", nil)
			}

			for _, varStr := range instr.VarStrings {
				l1 += addSpaces(fmt.Sprintf("%s", instr.Vars[varStr].Type), 15)
				l2 += addSpaces(fmt.Sprintf("%s", varStr), 15)
				l3 += addSpaces(fmt.Sprintf("0x%X", instr.Vars[varStr].Value), 15)
			}

			log(address+mnemonic+length+operandCount+mode+raw+"\n", nil)
			log(shortDesc, nil)

			if instr.VarCount > 0 {
				log(addSpacesL(l1, 15), nil)
				log(addSpacesL(l2, 15), nil)
			}
			log(addSpacesL(l3, 15), nil)

		}

		opSize = instr.ByteLength

	}
	log(fmt.Sprintf("Found [%d] instructions", count), nil)
	return nil

}

func addSpaces(s string, w int) string {
	if len(s) < w {
		s += strings.Repeat(" ", w-len(s))
	}
	return s
}

func addSpacesL(s string, w int) string {
	l := ""
	if len(s) < w {
		l += strings.Repeat(" ", w-len(s))
	}
	l += s
	return l
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
		fmt.Printf(" %s\n", kind)
	} else {
		fmt.Printf("[ERROR - %s]: %s\n", kind, err)
	}
}
