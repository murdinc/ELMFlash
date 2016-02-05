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
	//h.block = append(h.block, h.block[0x100000:0x180000]...)

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
	// BIG TODO sort this?
	pcs := []int{
		0x1068AF,
		0x106683,
		0x1068E3,
		0x13FAD7,
		0x13FE80,
		0x100000,
		0x1000D4,
		0x100675,
		0x100709,
		0x1060EA,
		0x1061C4,
		0x106326,
		0x10638F,
		0x106405,
		0x1064A0,
		0x10651A,
		0x1754B6,
		0x17553B,
		0x17F4C4,
		0x17F4B9,
		0x175642,
		0x17517A,
		0x172765,
		0x10086A,
		0x1008D9,
		0x1061AF,
		0x100996,
		0x1061DB,
		0x106201,
		0x106220,
		0x1063B2,
		0x10641E,
		0x10646D,
		0x1064B1,
		0x10654D,
		0x17F4C5,
		0x1645D7,
		0x167D7F,
		0x16A852,
		0x16FF5E,
		0x170000,
		0x171116,
		0x1007CC,
		0x100A06,
		0x1061EE,
		0x106221,
		0x106433,
		0x1065A6,
		0x13E704,
		0x13E72C,
		0x13E769,
		0x13E7EF,
		0x13E8A9,
		0x13E96B,
		0x13E9C7,
		0x13EA73,
		0x13EAD9,
		0x13EB49,
		0x13EBF7,
		0x13EC34,
		0x13ECB4,
		0x13EE07,
		0x13EE86,
		0x100ABC,
		0x100B4E,
		0x100B96,
		0x100CB0,
		0x100EA2,
		0x100F5B,
		0x100FE6,
		0x101178,
		0x1012C4,
		0x1013BF,
		0x1014DC,
		0x101622,
		0x10166C,
		0x10182B,
		0x101A4E,
		0x101B49,
		0x106222,
		0x1065EA,
		0x158D36,
		0x15BC70,
		0x15C7E8,
		0x15CA66,
		0x15CDA7,
		0x15DAAE,
		0x15DFB0,
		0x163722,
		0x16576F,
		0x16C22F,
		0x1715AD,
		0x1007DF,
		0x100CF3,
		0x17615A,
		0x16E3E6,
		0x15E527,
		0x15CF51,
		0x15CAA8,
		0x1761CB,
		0x10120C,
		0x1016D4,
		0x10185A,
		0x101B97,
		0x106223,
		0x10664C,
		0x106705,
		0x13E77E,
		0x13E807,
		0x13E8F4,
		0x13E994,
		0x13E9C8,
		0x13E9DA,
		0x13EB81,
		0x13ECF6,
		0x13EE4D,
		0x13EE94,
		0x13EF6B,
		0x13F099,
		0x13F135,
		0x101BE5,
		0x106224,
		0x10612C,
		0x106225,
		0x101729,
		0x1069CF,
		0x106144,
		0x10665E,
		0x106755,
		0x13E81F,
		0x13ED39,
		0x13F168,
		0x13F329,
		0x13F3F3,
		0x176A47,
		0x176E2E,
		0x176FC7,
		0x177153,
		0x177629,
		0x17763C,
		0x13FA31,
		0x14033E,
		0x140737,
		0x140763,
		0x140797,
		0x140A9F,
		0x14104A,
		0x1068FD,
		0x1413E3,
		0x14D0F2,
		0x14E17B,
		0x14EFAA,
		0x14FB6D,
		0x15126F,
		0x141207,
		0x1412EA,
		0x1415B0,
		0x1415F8,
		0x14178D,
		0x141973,
		0x10666E,
		0x14AFBB,
		0x16DF20,
		0x1756A9,
		0x1475EB,
		0x1067AB,
		0x172000,
		0x14DD54,
		0x152165,
		0x15210B,
		0x151738,
		0x150000,
		0x140B66,
		0x15CAB2,
		0x15C822,
		0x15C529,
		0x159639,
		0x1574FA,
		0x156121,
		0x155CDC,
		0x155586,
		0x154A16,
		0x153051,
		0x10615C,
		0x10676B,
		0x16E438,
		0x17215D,
		0x106848,
		0x13E83E,
		0x15DABB,
		0x15CFF7,
		0x15EFF1,
		0x14B8B5,
		0x149574,
		0x13F184,
		0x13F38C,
		0x13F45A,
		0x10676C,
		0x13FA65,
		0x1403F3,
		0x140857,
		0x13E860,
		0x1413D6,
		0x141A68,
		0x13F52D,
		0x1426DE,
		0x147FFB,
		0x147B09,
		0x14B930,
		0x156479,
		0x13FACA,
		0x13FC0E,
		0x13FCA6,
		0x13FDDE,
		0x13F599,
		0x10685F,
		0x14B89E,
		0x142686,
		0x142F9A,
		0x14325C,
		0x143675,
		0x143C32,
		0x144496,
		0x14457F,
		0x1446B5,
		0x1448DB,
		0x14516A,
		0x146A75,
		0x146EC0,
		0x147868,
		0x147D52,
		0x147DB0,
		0x1425E6,
		0x142BE4,
		0x13F224,
		0x174E96,
		0x1751AE,
		0x1756A0,
		0x175DE1,
		0x17F4BA,
		0x1000A1,
		0x110000,
		0x100220,
		0x100384,
		0x1003E8,
		0x100419,
		0x10055A,
		0x102080,
		0x112080,
		0x13F9DE,
		0x172080,
		0x17F5D4,
		0x17F6D2,
		0x17F759,
		0x17F83B,
		0x17F944,
		0x17FA67,
		0x17FAE0,
		0x17F559,
		0x17FB3A,
		0x17FB7D,
		0x17FC37,
		0x17FD77,
		0x17FF63,
		0x17FFF1,
		0x17FFF2,
		0x17FFF3,
		0x17FFF4,
		0x17FFF5,
		0x17FFF9,
		0x106D4C,
		0x1720BD,
	}

	for p := 0; p < len(pcs); p++ {
		pc := pcs[p]

		//log(fmt.Sprintf("Start Address: [0x%X]", pc), nil)

	Loop:
		for {
			// Sub and Jumps, or break if out of range
			if pc+10 > len(h.block) {
				pc = 0xFFFFFF
				if other[pc] {
					crawled[pc] = 1
				}

				// Conditional Jumps
				for adr, _ := range jumps {
					if crawled[adr] == 0 {
						pc = adr
						//log(fmt.Sprintf("CRAWLING JUMP ADDRESS 0x%X", adr), nil)
						continue Loop
					}
				}

				// Subroutines
				for adr, _ := range subroutines {
					if crawled[adr] == 0 {
						pc = adr
						//log(fmt.Sprintf("CRAWLING SUB ADDRESS 0x%X", adr), nil)
						continue Loop
					}
				}

				// Other
				for adr, crl := range other {
					if crl == false {
						other[adr] = true
						pc = adr
						//log(fmt.Sprintf("CRAWLING OTHER ADDRESS 0x%X", adr), nil)
						continue Loop
					}
				}

				//log(fmt.Sprintf("No more suff to crawl? : 0x%X", pc), nil)
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
				/*
					log("Parser", err)
					log(fmt.Sprintf("ERROR!! Address: 0x%X		Instruction %X", pc, b), err)
				*/
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
				returns++
				other[pc+instr.ByteLength] = false
				pc = 0xFFFFFF
				continue Loop
			}

			// If we havent unconditionally jumped, move our program counter the length of this op
			pc += instr.ByteLength

		}
	}

	//log(fmt.Sprintf("Found [%d] instructions", count), nil)
	log(fmt.Sprintf("Found [%d] XRefs", len(xrefs)), nil)
	log(fmt.Sprintf("Found [%d] Subroutines", len(subroutines)), nil)
	log(fmt.Sprintf("Found [%d] Returns", returns), nil)
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
			shortDesc := addSpaces(fmt.Sprintf("%s %s", instr.Description, instr.Mnemonic), 45)
			//shortDesc := addSpaces(fmt.Sprintf("%s %s 0x%X 0x%X", instr.Description, instr.Mnemonic, instr.Op, instr.Raw), 45)

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
			if chkAdr >= 0x180000 {
				break Check
			} else if crawled[chkAdr] != 1 {
				address := addSpaces(fmt.Sprintf("[0x%X] %d ?: ", chkAdr, crawled[chkAdr]), 20)
				shortDesc := fmt.Sprintf("%.2X ", h.block[chkAdr])
				for i := 1; i < 32; i++ {
					if i%8 == 0 {
						shortDesc += " "
					}
					chkAdr++
					if chkAdr >= 0x180000 {
						break
					} else if crawled[chkAdr] != 1 {
						shortDesc += fmt.Sprintf("%.2X ", h.block[chkAdr])
					} else {
						break
					}
				}
				log(address+shortDesc, nil)
			} else {
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
