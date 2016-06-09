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
	block           []byte
	intRoutineAdrs  []int          // slice of interrupt routine addresses for start locations
	intRoutineNames map[int]string // address of interrupt routine locations and name
	vectorAdr       map[int]string // address of interrupt vector locations and name
	memStarts       map[int]string // Starts of memory map Locations
	memStops        map[int]string // Ends of memory map Locations
	skip            map[int]int
}

var calibrations = map[string]string{
	"msp":   "MSP.BIN",
	"mp3":   "MP3.BIN",
	"mp3x2": "MP3x2.BIN",
	"p5":    "P5.BIN",
	"pre":   "PRE2.BIN",
}

func New(calName string) *DisAsm {
	controller := new(DisAsm)

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
	preBlock := make([]byte, preFileSize)
	calBlock := make([]byte, fileSize)

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

	block := append(preBlock, calBlock...)

	controller.block = block

	return controller
}

func (h *DisAsm) DisAsm() error {

	h.GetInterrupts()
	h.GetMemoryMap()

	log(fmt.Sprintf("Length: 0x%X", len(h.block)), nil)

	var opcodes Instructions
	subroutines := make(map[int][]Call)
	xrefs := make(map[int][]XRef)
	jumps := make(map[int][]Jump)
	other := make(map[int]bool)
	crawled := make(map[int]int)
	returns := 0
	errors := 0

	// Program Counter - Start Address: 0x172080
	pcs := []int{0x172080}
	pcs = append(pcs, h.intRoutineAdrs...)

	loops := 50

	for p := 0; p < len(pcs)+loops; p++ {

		pc := 0xFFFFFF

		if p < len(pcs) {
			pc = pcs[p]
		}

	Loop:
		for {

			// Sub and Jumps, or break if out of range
			if pc+10 > len(h.block) {
				if pc != 0xFFFFFF {
					crawled[pc] = 1
					pc &= 0x17FFFF
				} else {
					pc = 0xFFFFFF
					if other[pc] {
						crawled[pc] = 1
					}
				}

				// Conditional Jumps
				for adr, _ := range jumps {
					if crawled[adr] == 0 {
						pc = adr
						continue Loop
					}
				}

				// Subroutines
				for adr, _ := range subroutines {
					if crawled[adr] == 0 {
						pc = adr
						continue Loop
					}
				}

				// Other
				for adr, crl := range other {
					if crl == false && crawled[adr] == 0 {
						other[adr] = true
						pc = adr
						continue Loop
					}
				}

				break Loop
			}

			if crawled[pc] == 1 {
				pc = 0xFFFFFF
				continue Loop
			}

			// The Parser™
			b := h.block[pc : pc+10]
			instr, err := Parse(b, pc)
			crawled[pc] = 1
			for i := 1; i < instr.ByteLength; i++ {
				crawled[i+pc] = 1
			}

			if err != nil {
				errors++
				log(fmt.Sprintf("ERROR!! Address: 0x%X		Instruction %X", pc, b), err)
				crawled[pc] = 3
				pc = 0xFFFFFF
				continue Loop
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
				// If this is not a conditional jump, point the program counter at the address
				switch instr.Mnemonic {
				case "SJMP", "EJMP", "LJMP", "TIJMP":
					jumps[JumpAdd] = append(jumps[JumpAdd], JumpVal...)
					//log(instr.Mnemonic, nil)
					pc = JumpAdd
					continue Loop
				case "EBR", "BR":
					pc = 0xFFFFFF // TODO!!!!!!!!
					continue Loop
				default:
					jumps[JumpAdd] = append(jumps[JumpAdd], JumpVal...)
				}

			}

			// Subroutine Returns and Resets {
			if instr.Mnemonic == "RET" || instr.Mnemonic == "RST" {
				returns++
				pc = 0xFFFFFF
				continue Loop
			}

			// If we havent unconditionally jumped, move our program counter the length of this op
			pc += instr.ByteLength

		}
	}

	log(fmt.Sprintf("Found [%d] instructions", len(opcodes)), nil)
	log(fmt.Sprintf("Found [%d] XRefs", len(xrefs)), nil)
	log(fmt.Sprintf("Found [%d] Subroutines", len(subroutines)), nil)
	log(fmt.Sprintf("Found [%d] Returns", returns), nil)
	log(fmt.Sprintf("Found [%d] Jumps", len(jumps)), nil)

	sort.Sort(opcodes)

	// Print out the stuff before the Assembly
	for chkAdr := 0; chkAdr < opcodes[0].Address; chkAdr++ {

		chkAdr += h.doMemoryMap(chkAdr)

		if xrefs[chkAdr] != nil {
			referers := "  "
			for i, referrer := range xrefs[chkAdr] {
				referers += addSpaces(fmt.Sprintf("    [ XREF 0x%X - %s ] ", referrer.XRefFrom, referrer.Mnemonic), 35)
				if (i+1)%4 == 0 || i+1 == len(xrefs[chkAdr]) {
					referers += "\n  "
				}
			}

			log(fmt.Sprintf("======== XREF_ 0x%X %s \n%s", chkAdr, regName("", chkAdr), referers), nil)

			address := addSpaces(fmt.Sprintf("[0x%X]   X: ", chkAdr), 20)
			shortDesc := fmt.Sprintf("%.2X ", h.block[chkAdr])

			log(address+shortDesc, nil)

		} else if crawled[chkAdr] != 1 {
			address := addSpaces(fmt.Sprintf("[0x%X] %d ?: ", chkAdr, crawled[chkAdr]), 20)
			shortDesc := fmt.Sprintf("%.2X ", h.block[chkAdr])

			for i := 1; i < 32; i++ {
				if i%8 == 0 {
					shortDesc += " "
				}
				chkAdr++
				if chkAdr >= len(h.block) {
					break
				} else if h.memStarts[chkAdr] != "" || h.memStops[chkAdr] != "" {
					shortDesc += fmt.Sprintf("%.2X ", h.block[chkAdr])
					break
				} else if crawled[chkAdr] != 1 && xrefs[chkAdr] == nil {
					shortDesc += fmt.Sprintf("%.2X ", h.block[chkAdr])
				} else {
					chkAdr--
					break
				}
			}
			log(address+shortDesc, nil)
		}

	}

	// Print out the Assembly

	for index, instr := range opcodes {

		h.doMemoryMap(instr.Address)

		if subroutines[instr.Address] != nil {
			callers := ""
			for _, caller := range subroutines[instr.Address] {
				callers = callers + fmt.Sprintf("  ============================================================= [CALLED FROM 0x%X - %s] \n", caller.CallFrom, caller.Mnemonic)
			}
			log(fmt.Sprintf("\n======== SUBROUTINE_ 0x%X ==================================================================================\n%s", instr.Address, callers), nil)
		}

		if h.intRoutineNames[instr.Address] != "" {

			log(fmt.Sprintf("\n======== INTERRUPT ROUTINE_ %s ==================================================================================", h.intRoutineNames[instr.Address]), nil)
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
				log("\n== RETURN FROM SUBROUTINE ===============================================================================\n", nil)
			}

		}

		// Print our uncrawled addresses
		chkAdr := instr.Address + instr.ByteLength
	Check:
		for {
			//h.doMemoryMap(chkAdr)

			if chkAdr >= len(h.block) {
				h.doMemoryMap(instr.Address)
				break Check

			} else if h.vectorAdr[chkAdr] != "" { // Vector Addresses
				address := addSpaces(fmt.Sprintf("[0x%X]   V: ", chkAdr), 20)
				shortDesc := fmt.Sprintf("%.2X		(%s)", h.block[chkAdr:chkAdr+2], h.vectorAdr[chkAdr])
				log(address+shortDesc, nil)
				chkAdr++

			} else if crawled[chkAdr] != 1 { // Crawled but not parsed
				address := addSpaces(fmt.Sprintf("[0x%X] %d ?: ", chkAdr, crawled[chkAdr]), 20)
				shortDesc := fmt.Sprintf("%.2X ", h.block[chkAdr])
				for i := 1; i < 32; i++ {
					if i%8 == 0 {
						shortDesc += " "
					}
					chkAdr++
					if chkAdr >= len(h.block) {
						h.doMemoryMap(instr.Address)
						break
					} else if h.memStarts[chkAdr] != "" || h.memStops[chkAdr] != "" {
						shortDesc += fmt.Sprintf("%.2X ", h.block[chkAdr])
						break
					} else if crawled[chkAdr] != 1 && xrefs[chkAdr] == nil && h.vectorAdr[chkAdr] == "" {
						shortDesc += fmt.Sprintf("%.2X ", h.block[chkAdr])
					} else {
						chkAdr--
						break
					}
				}
				log(address+shortDesc, nil)
			} else { // Bomb out
				break Check
			}
			chkAdr++
		}

	}

	count := 0
	for i := 0x100000; i < len(h.block); i++ {

		if crawled[i] == 0 {
			count++
		}
	}

	log(fmt.Sprintf("UNCRAWLED ADDRESSES: %d", count), nil)
	log(fmt.Sprintf("CRAWLED ADDRESSES: %d", len(crawled)), nil)
	log(fmt.Sprintf("ERRORS: %d", errors), nil)
	log(fmt.Sprintf("Found [%d] instructions", len(opcodes)), nil)

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
