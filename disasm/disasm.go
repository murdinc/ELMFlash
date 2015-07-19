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
}

func (h *DisAsm) DisAsm(calName string) error {

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

	log(fmt.Sprintf("Disassemble - [%s] is %d bytes long", calibrations[calName], fileSize), nil)

	// Make a buffer
	block := make([]byte, fileSize)

	// Read in all the bytes bytes
	n, err := f.Read(block)
	if err != nil {
		log("UploadBIN - Error reading calibration", err)
		return err
	}

	dbg(fmt.Sprintf("BIN - reading %d bytes.", n), nil)

	//for i, b := range block {
	opSize := 1
	count := 1
	for i := 0x008000; i < int(fileSize); i = i + opSize {

		b := block[i:]
		instr, err := Parse(b)

		if err != nil {
			log("", err)
		} else if instr.Ignore == false {
			count++
			log("---------", nil)

			address := addSpaces(fmt.Sprintf("Address: [0x%X]", i), 25)
			length := addSpaces(fmt.Sprintf(" Length: [%d]", instr.ByteLength), 16)
			mode := addSpaces(fmt.Sprintf(" Mode: [%s]", instr.AddressingMode), 28)
			//addrModeRef :=
			mnemonic := addSpaces(fmt.Sprintf(" Mnemonic: [%s] %s", instr.Mnemonic, instr.AddrModeRef), 40)
			operandCount := addSpaces(fmt.Sprintf(" Operand Count: [%d]", instr.OperandCount), 50)
			raw := addSpaces(fmt.Sprintf(" Raw: 0x%.10X", instr.Raw), 20)

			log(address+length+mode+mnemonic+operandCount+raw, nil)

		}

		opSize = instr.ByteLength

		/*
			b := block[i:]
			firstByte := b[0]
			var secondByte byte

			if i < int(fileSize-1) {
				secondByte = b[1]
			}

			if contains(firstByte, keys(special)) && (special[firstByte] != "multi") {
				// This is a "special" command
				log("---------------------------------------------------------------------------------------------------------------------------------------------------------------------------", nil)
				log(fmt.Sprintf("Special Opcode: [%X] - [%s]", firstByte, special[firstByte]), nil)
			}

			if firstByte != 0xFE && contains(firstByte, keys(mnemonics)) {

				log("---------------------------------------------------------------------------------------------------------------------------------------------------------------------------", nil)

				// everything except FE
				opSize = byteLengths[firstByte]
				if opSize == 0 {
					opSize = 1
					continue
				}

				// Check if this is a variable length opcode and its , and adjust accordingly
				if contains(firstByte, variableLengths) && (secondByte&1 == 1) {
					log(fmt.Sprintf("Variable first byte and even second byte: [0x%X] - [%d]", secondByte, secondByte), nil)
					opSize++
				}

				address := addSpaces(fmt.Sprintf("Address: [0x%X]", i), 25)
				length := addSpaces(fmt.Sprintf(" Length: [%d]", opSize), 16)
				mode := addSpaces(fmt.Sprintf(" Mode: [%s]", modes[firstByte]), 28)
				mnemonic := addSpaces(fmt.Sprintf(" Mnemonic: [%s]", mnemonics[firstByte]), 70)
				opLine := addSpaces(fmt.Sprintf(" Op Code: 0x%.10X", b[0:opSize]), 20)

				newData, _ := Parse(b)

				if newData.ByteLength == opSize {
					continue
				} else {
					log("MISMATCH",nil)
				}

				log(address+length+mode+mnemonic+opLine, nil)
				log(fmt.Sprintf("New Data: %d", newData.ByteLength), nil)
				count++
			} else if firstByte == 0xFE && contains(secondByte, keys(mnemonicsSigned)) {
				log("---------------------------------------------------------------------------------------------------------------------------------------------------------------------------", nil)
				// FE (signed multi and div) opcodes
				opSize = byteLengthsSigned[secondByte]
				if opSize == 0 {
					opSize = 1
					continue
				}

				address := addSpaces(fmt.Sprintf("Address: [0x%X]", i), 25)
				length := addSpaces(fmt.Sprintf(" Length: [%d]", opSize), 16)
				mode := addSpaces(fmt.Sprintf(" Mode: [%s]", modesSigned[secondByte]), 28)
				mnemonic := addSpaces(fmt.Sprintf(" Mnemonic: [Signed %s]", mnemonicsSigned[secondByte]), 70)
				opLine := addSpaces(fmt.Sprintf(" Op Code: 0x%.10X", b[0:opSize]), 20)

				log(address+length+mode+mnemonic+opLine, nil)
				//log(fmt.Sprintf("Address:   0x%X      length: [%d]    modes: [%s]     [Signed %s] 0x%.10X", i, opSize, modesSigned[secondByte], mnemonicsSigned[secondByte], b[0:opSize]), nil)

				count++
			} else {
				// All Else has Failed? Or we commented it out because its a NOP?
				opSize = 1
			}
		*/

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
		//fmt.Printf("====> %s\n", kind)
		fmt.Printf("	%s\n", kind)
	} else {
		fmt.Printf("[ERROR - %s]: %s\n", kind, err)
	}
}
