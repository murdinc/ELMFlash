package compare

import (
	"fmt"
	"os"
	"strings"
)

// App constants
////////////////..........
const debug = false

type Compare struct {
	block1 []byte
	block2 []byte
}

var calibrations = map[string]string{
	"msp":   "MSP.BIN",
	"mp3":   "MP3.BIN",
	"mp3x2": "MP3x2.BIN",
	"p5":    "P5.BIN",
	"pre":   "PRE.BIN",
	"pre2":  "PRE2.BIN",
}

func New(preName1 string, calName1 string, preName2 string, calName2 string) *Compare {
	cmp := new(Compare)

	cmp.block1 = buildCal(preName1, calName1)
	cmp.block2 = buildCal(preName2, calName2)

	return cmp
}

func buildCal(preName string, calName string) []byte {
	// Pull in the stuff before the calibration file
	preCalFile := "./calibrations/" + calibrations[preName]
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

	return block
}

func (c *Compare) Compare() error {

	width := 32

	for i := 0; (i < len(c.block1)) && i < len(c.block2); i++ {

		address := addSpaces(fmt.Sprintf("[0x%X] ", i), 20)
		shortDesc := fmt.Sprintf("%.2X ", c.block1[i])

		for j := 1; j < width; j++ {
			if j%8 == 0 {
				shortDesc += " "
			}

			if i >= len(c.block1) {
				break
			} else {
				shortDesc += fmt.Sprintf("%.2X ", c.block1[i])
				i++
			}
		}

		shortDesc += "		"

		i -= (width - 1)

		for j := 1; j < width; j++ {
			if j%8 == 0 {
				shortDesc += " "
			}

			if i >= len(c.block2) {
				break
			} else {
				if c.block1[i] == c.block2[i] {
					shortDesc += "** "
					i++
				} else {
					shortDesc += fmt.Sprintf("%.2X ", c.block2[i])
					i++
				}
			}
		}

		log(address+shortDesc, nil)

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
