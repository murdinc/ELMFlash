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

type DisAsm struct {
	location string
	block    []byte
}

var calibrations = map[string]string{
	"msp": "MSP.BIN",
	"mp3": "MP3.BIN",
	"p5":  "P5.BIN",
	"pre": "PRE.BIN",
}

func New(calName string) *DisAsm {
	controller := new(DisAsm)
	controller.location = "MSP.BIN"

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
		os.Exit(1)
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
		os.Exit(1)
	}
	log(fmt.Sprintf("Disassemble - reading 0x%X bytes from pre-calibration file.", n), nil)

	n, err = f.Read(calBlock)
	if err != nil {
		log("Disassemble - Error reading calibration", err)
		os.Exit(1)
	}

	log(fmt.Sprintf("Disassemble - reading 0x%X bytes from calibration file.", n), nil)

	//block := append(preBlock, calBlock...)

	block := append(preBlock, calBlock...)

	controller.block = block

	return controller
}

func (h *DisAsm) DisAsm() error {

	// Double Trouble
	h.block = append(h.block, h.block[0x100000:0x180000]...)

	log(fmt.Sprintf("Length: 0x%X", len(h.block)), nil)

	var opcodes Instructions
	subroutines := make(map[int][]Call)
	xrefs := make(map[int][]XRef)
	jumps := make(map[int][]Jump)
	crawled := make(map[int]bool)

	// Program Counter - Start Address: 0x172080
	pc := 0x172080
	log(fmt.Sprintf("Start Address: [0x%X]", pc), nil)

Loop:
	for {

		// Sub and Jumps, or break if out of range
		if pc+10 > len(h.block) {
			crawled[pc] = true

			// Conditional Jumps
			for adr, _ := range jumps {
				if crawled[adr] == false {
					pc = adr
					log(fmt.Sprintf("CRAWLING JUMP ADDRESS 0x%X", adr), nil)
					continue Loop
				}
			}

			// Subroutines
			for adr, _ := range subroutines {
				if crawled[adr] == false {
					pc = adr
					log(fmt.Sprintf("CRAWLING SUB ADDRESS 0x%X", adr), nil)
					continue Loop
				}
			}

			log("here", nil)
			log(fmt.Sprintf("No more suff to crawl? : 0x%X", pc), nil)
			break Loop
		}

		if crawled[pc] == true {
			pc = 0xFFFFFF
			continue Loop
		}

		// The Parser™
		b := h.block[pc : pc+10]
		instr, err := Parse(b, pc)
		crawled[pc] = true
		if err != nil {
			log("Parser", err)
			log(fmt.Sprintf("Address: 0x%X		Instruction %X", pc, b), err)

			//break Loop
		}

		// Append our instruction to our opcodes list
		opcodes = append(opcodes, instr)

		// Append our XRefs to our XRefs list
		for XRefAdd, XRefVal := range instr.XRefs {
			xrefs[XRefAdd] = append(xrefs[XRefAdd], XRefVal...)
		}

		// Append our Call addresses to the subroutines list
		for CallAdd, CallVal := range instr.Calls {
			subroutines[CallAdd] = append(subroutines[CallAdd], CallVal...)
		}

		// Append our Jumps to our Jumps list
		for JumpAdd, JumpVal := range instr.Jumps {
			jumps[JumpAdd] = append(jumps[JumpAdd], JumpVal...)

			// If this is not a conditional jump, point the program counter at the address
			switch instr.Mnemonic {
			case "SJMP", "EJMP", "LJMP", "TIJMP":
				//log(instr.Mnemonic, nil)
				pc = JumpAdd
				continue Loop
			}
		}

		// Subroutine Returns {
		if instr.Mnemonic == "RET" {
			log("						##### RETURN", nil)
			pc = 0xFFFFFF
			continue Loop
		}

		// If we havent unconditionally jumped, move our program counter the length of this op
		pc += instr.ByteLength

	}

	//log(fmt.Sprintf("Found [%d] instructions", count), nil)
	log(fmt.Sprintf("Found [%d] XRefs", len(xrefs)), nil)
	log(fmt.Sprintf("Found [%d] Subroutines", len(subroutines)), nil)
	log(fmt.Sprintf("Found [%d] Jumps", len(jumps)), nil)

	sort.Sort(opcodes)

	// Print out the Assembly
	for index, instr := range opcodes {

		if subroutines[instr.Address] != nil {
			callers := ""
			for _, caller := range subroutines[instr.Address] {
				callers = callers + fmt.Sprintf("  ============================================================= [CALLED FROM 0x%X - %s] \n", caller.CallFrom, caller.Mnemonic)
			}
			log(fmt.Sprintf("\n======== SUB_ 0x%X ==================================================================================\n%s", instr.Address, callers), nil)

		}

		if jumps[instr.Address] != nil {
			jumpers := ""
			for _, jumper := range jumps[instr.Address] {
				jumpers = jumpers + fmt.Sprintf("  ============================================================= [JUMP FROM 0x%X - %s] \n", jumper.JumpFrom, jumper.Mnemonic)
			}
			log(fmt.Sprintf("\n======== JUMP_ 0x%X \n%s", instr.Address, jumpers), nil)

		}

		if instr.Ignore == false {

			if instr.Mnemonic == "CMPB" || instr.Mnemonic == "CMP" {
				switch opcodes[index+1].Mnemonic {

				//case "JNST":
				case "JNH":
					instr.PseudoCode = strings.Replace(instr.PseudoCode, "==", "<=", 1)
				case "JGT":
					instr.PseudoCode = strings.Replace(instr.PseudoCode, "==", ">", 1)
				//case "JNC":
				//case "JNVT":
				//case "JNV":
				case "JGE":
					instr.PseudoCode = strings.Replace(instr.PseudoCode, ">=", "!=", 1)
				case "JNE":
					instr.PseudoCode = strings.Replace(instr.PseudoCode, "==", "!=", 1)
				//case "JST":
				case "JH":
					instr.PseudoCode = strings.Replace(instr.PseudoCode, "==", ">", 1)
				case "JLE":
					instr.PseudoCode = strings.Replace(instr.PseudoCode, "==", "<=", 1)
				//case "JC":
				//case "JVT":
				//case "JV":
				case "JLT":
					instr.PseudoCode = strings.Replace(instr.PseudoCode, "==", "<", 1)

				}

			}

			address := addSpaces(fmt.Sprintf("[0x%X]", instr.Address), 20)
			//shortDesc := addSpaces(fmt.Sprintf("%s %s", instr.Description, instr.Mnemonic), 45)
			shortDesc := addSpaces(fmt.Sprintf("%s %s 0x%X 0x%X", instr.Description, instr.Mnemonic, instr.Op, instr.Raw), 45)

			var l1 string

			if !instr.Checked {
				log("#### ERROR DISASEMBLING OPCODE ####", nil)
			}

			// Pseudo Code
			l1 = addSpaces(l1, 15)
			l1 += fmt.Sprintf("%s", instr.PseudoCode)

			log(address+shortDesc+l1, nil)

			if instr.Mnemonic == "RET" {
				log("\n== RETURN FROM SUBROUTINE ===============================================================================", nil)
			}
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
