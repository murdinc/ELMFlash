package disasm

import (
	"fmt"
	"os"
	"sort"
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
	"p5":  "P5.BIN",
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
	count := 0

	var opcodes []Instruction
	subroutines := make(map[int][]Call)
	xrefs := make(map[int][]XRef)
	jumps := make(map[int][]Jump)

	for i := 0x100000; i < len(block); i = i + opSize {
		//for i := 0x000000; i < len(block); i = i + opSize {

		// Registers and Ram
		if i > 0x0 && i < 0xFFF {
			continue
		}

		// Unknown
		if i > 0x1BFF && i < 0x100000 {
			continue
		}

		// Maps, it seems
		if i > 0x108000 && i < 0x11FFFF {
			continue
		}

		// More Maps?
		/*
			if i > 0x139D54 {
				continue
			}
		*/

		// Checksum
		if i == 0x103FDE {
			log(fmt.Sprintf("Checksum: [0x%X] Address: 0x%X", block[i:i+2], i), nil)
			opSize = 2
			continue
		}

		// Unknown, but def not opcode
		if i >= 0x1074D4 && i <= 0x1076A3 {
			opSize = 1
			continue
		}

		// Copy block and the stuff around it
		if i >= 0x107FFE && i <= 0x108103 {
			opSize = 1
			continue
		}

		// The Parser™
		b := block[i : i+10]
		instr, err := Parse(b, i)

		// Ignore ops set to ignore in our ops list
		if instr.Ignore == false && err == nil {
			count++
		}

		// Append our instruction to our opcodes list
		opcodes = append(opcodes, instr)

		// Append our Call addresses to the subroutines list
		for CallAdd, CallVal := range instr.Calls {
			subroutines[CallAdd] = append(subroutines[CallAdd], CallVal...)
		}

		// Append our XRefs to our XRefs list
		for XRefAdd, XRefVal := range instr.XRefs {
			xrefs[XRefAdd] = append(xrefs[XRefAdd], XRefVal...)
		}

		// Append our Jumps to our Jumps list
		for JumpAdd, JumpVal := range instr.Jumps {
			jumps[JumpAdd] = append(jumps[JumpAdd], JumpVal...)
		}

		//}

		opSize = instr.ByteLength

	}
	log(fmt.Sprintf("Found [%d] instructions", count), nil)
	log(fmt.Sprintf("Found [%d] XRefs", len(xrefs)), nil)
	log(fmt.Sprintf("Found [%d] Subroutines", len(subroutines)), nil)
	log(fmt.Sprintf("Found [%d] Jumps", len(jumps)), nil)

	// Print out the Assembly
	for _, instr := range opcodes {

		if subroutines[instr.Address] != nil {
			log("\n==================================================================================================================================================================", nil)
			callers := ""
			for _, caller := range subroutines[instr.Address] {
				callers = callers + fmt.Sprintf("[CALLED FROM 0x%X - %s] ", caller.CallFrom, caller.Mnemonic)
			}
			log(fmt.Sprintf("SUB_0x%X %s", instr.Address, callers), nil)

		}

		if jumps[instr.Address] != nil {
			log("\n==================================================================================================================================================================", nil)
			jumpers := ""
			for _, jumper := range jumps[instr.Address] {
				jumpers = jumpers + fmt.Sprintf("[JUMP FROM 0x%X - %s] ", jumper.JumpFrom, jumper.Mnemonic)
			}
			log(fmt.Sprintf("JUMP_0x%X %s", instr.Address, jumpers), nil)

		}

		if instr.Ignore == false {

			address := addSpaces(fmt.Sprintf("Address: [0x%X]", instr.Address), 20)
			length := addSpaces(fmt.Sprintf(" Length: [%d]", instr.ByteLength), 14)
			mode := addSpaces(fmt.Sprintf(" Mode: [%s]", instr.AddressingMode), 26)
			mnemonic := addSpaces(fmt.Sprintf(" Mnemonic: [%s]", instr.Mnemonic), 23)
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
				l1 += addSpaces(fmt.Sprintf("%s", instr.Vars[varStr].Type), 25)
				l2 += addSpaces(fmt.Sprintf("%s", varStr), 25)
				l3 += addSpaces(fmt.Sprintf("%s", instr.Vars[varStr].Value), 25)
			}

			// Pseudo Code
			l3 = addSpaces(l3, 150)
			l3 += fmt.Sprintf("%s", instr.PseudoCode)

			log(address+mnemonic+length+operandCount+mode+raw+"\n", nil)
			log(shortDesc, nil)

			if instr.VarCount > 0 {
				log(addSpacesL(l1, 15), nil)
				log(addSpacesL(l2, 15), nil)
			}
			log(addSpacesL(l3, 15), nil)

			if instr.Mnemonic == "RET" {
				log("\n== RETURN FROM SUBROUTINE ========================================================================================================================================", nil)
			}
		}

	}

	// Sort and print the list of address references collected
	log("\n== ADDRESS REFERENCES ========================================================================================================================================", nil)
	var keys []int
	for k := range xrefs {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for _, k := range keys {

		log(fmt.Sprintf("0x%X", k), nil)

		for _, ref := range xrefs[k] {
			log(fmt.Sprintf("       [%s] XREF  [%s] AT 0x%X", ref.String, ref.Mnemonic, ref.XRefFrom), nil)
		}
	}

	// Print the list of subroutines
	log("\n== SUBROUTINE REFERENCES ========================================================================================================================================", nil)
	var sKeys []int
	for s := range subroutines {
		sKeys = append(sKeys, s)
	}

	sort.Ints(sKeys)

	for _, s := range sKeys {

		log(fmt.Sprintf("SUB_0x%X", s), nil)

		for _, sub := range subroutines[s] {
			log(fmt.Sprintf("       [%s] AT 0x%X", sub.Mnemonic, sub.CallFrom), nil)
		}
	}

	// Print the list of jumps
	log("\n== JUMP REFERENCES ========================================================================================================================================", nil)
	var jKeys []int
	for j := range jumps {
		jKeys = append(jKeys, j)
	}

	sort.Ints(jKeys)

	for _, j := range jKeys {

		log(fmt.Sprintf("JUMP_0x%X", j), nil)

		for _, sub := range jumps[j] {
			log(fmt.Sprintf("       [%s] AT 0x%X", sub.Mnemonic, sub.JumpFrom), nil)
		}
	}

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
