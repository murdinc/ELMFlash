package hexstuff

import (
	"fmt"
	"os"
	"regexp"

	"github.com/murdinc/legacy-cli"
)

// App constants
////////////////..........
const debug = false

type HexStuff struct {
	Calibration string
	block       []byte
}

func New(calName string) (*HexStuff, error) {
	controller := new(HexStuff)

	controller.Calibration = calName

	// Pull in the stuff before the calibration file
	preCalFile := "./calibrations/" + calibrations["pre"]
	log(fmt.Sprintf("TestM1 - Pre-calibration File: %s", preCalFile), nil)

	p, err := os.Open(preCalFile)
	pi, err := p.Stat()
	preFileSize := pi.Size()

	// Pull in the Calibration file
	calFile := "./calibrations/" + calibrations[calName]
	log(fmt.Sprintf("TestM1 - Calibration File: %s", calFile), nil)

	f, err := os.Open(calFile)
	fi, err := f.Stat()
	fileSize := fi.Size()
	if err != nil {
		log("TestM1 - Error opening file", err)
		return nil, err
	}

	log(fmt.Sprintf("TestM1 - [%s] is %d bytes long", calibrations["pre"], preFileSize), nil)
	log(fmt.Sprintf("TestM1 - [%s] is %d bytes long", calibrations[calName], fileSize), nil)

	// Make some buffers
	preBlock := make([]byte, 0x108000)
	calBlock := make([]byte, 0x78000)

	// Read in all the bytes
	n, err := p.Read(preBlock)
	if err != nil {
		log("TestM1 - Error reading calibration", err)
		return &HexStuff{}, err
	}
	log(fmt.Sprintf("TestM1 - reading 0x%X bytes from pre-calibration file.", n), nil)

	n, err = f.Read(calBlock)
	if err != nil {
		log("TestM1 - Error reading calibration", err)
		return &HexStuff{}, err
	}

	log(fmt.Sprintf("TestM1 - reading 0x%X bytes from calibration file.", n), nil)

	block := append(preBlock, calBlock...)

	controller.block = block

	log(fmt.Sprintf("Length: 0x%X", len(block)), nil)

	return controller, nil
}

var calibrations = map[string]string{
	"msp": "MSP.BIN",
	"mp3": "MP3.BIN",
	"p5":  "P5.BIN",
	"pre": "PRE.BIN",
}

func (h *HexStuff) TestM1() ([]int, error) {

	var addresses []int
	count := 0
	sizes := make(map[string]int)
	previous := 0x0000
	index := 0x108000

Loop:
	for {

		matches := FindMatch(h.block[index:])
		base := index

		for _, i := range matches {

			index = base + i[0]

			if index > 0x118000 {
				break Loop
			}

			//fmt.Printf("MATCH: 0x%X\n", index)

			// Round up if odd?
			if index%2 != 0 {
				index++
				//fmt.Printf("Rounding Up: %d > %d\n", index-1, index)
				continue Loop
			}

			height := int(h.block[index+4]) + 1
			width := int(h.block[index+5]) + 1

			// Hope?
			if height <= 1 || height > 32 || width > 50 {
				index += 2
				continue Loop
			}

			h2 := h.block[index+1] + 1
			h4 := h.block[index+3] + 1

			h7 := h.block[index+6] + 1

			h8 := h.block[index+7] + 1

			if index%2 == 0 && index >= previous && index < 0x118000 && height > 1 && h8 <= 11 && h8 > 0 && !(h2 == 2 && h4 == 2) && h7 < 100 {

				count++

				size := width * height
				start := index + 8
				end := start + size

				sixteen := (int(h.block[index+7]) << 8) | int(h.block[index+6])

				fmt.Printf("0x%X\n", previous)
				if previous != index && previous > 0 {
					missing := fmt.Sprintf(" MISSED: %X \n", h.block[previous:index])
					log(missing, nil)
				}

				match := fmt.Sprintf(" MATCH: -1 0x%X 0 0x%X +1 0x%X  ADDRESS: 0x%X  END: 0x%X  SIZE: %d x %d		[%d]	H2: %d      H4: %d      L: 0x%X	L: %d	|	START #%X | ROWS: %d x COLS: %d", h.block[index-8:index], h.block[index:index+8], h.block[end:end+8], index, end, width, height, size, h2, h4, h.block[index+6:index+8], sixteen, start, height, width)
				log(match, nil)

				sizeName := fmt.Sprintf("%d	x	%d", width, height)
				var sizeCount int
				if sizes[sizeName] < 1 {
					sizeCount = 1
				} else {
					sizeCount = sizes[sizeName]
					sizeCount++
				}
				sizes[sizeName] = sizeCount

				addresses = append(addresses, index)

				previous = end

				// Round up if odd?
				if previous%2 != 0 {
					previous++
				}

				index = previous

				continue Loop
			}

			if index%2 == 0 && index >= previous && index < 0x118000 && height > 1 && h8 <= 11 && h8 > 0 && !(h2 == 2 && h4 == 2) && h7 < 100 {

			}
		}
	}

	log(fmt.Sprintf("COUNT: %d ", count), nil)

	fmt.Println("")

	for size, count := range sizes {
		log(fmt.Sprintf("SIZE: %s		COUNT: %d", size, count), nil)
	}

	fmt.Println("")

	for _, address := range addresses {
		fmt.Printf("0x%X, ", address)
	}

	fmt.Println("")
	//fmt.Println(matches)
	fmt.Println("")

	return addresses, nil

}

func (h *HexStuff) TestM2() ([]int, error) {

	var addresses []int
	count := 0
	//sizes := make(map[string]int)

	// startIP := 0x10A500
	startIP := 0x108B00

	previousEnd := startIP

Loop:
	for index := startIP; index < len(h.block); index++ {

		if count >= 111 {
			break Loop
		}

		if index >= 0x110000 {
			break Loop
		}

		dataType := ""
		cols := 0
		rows := 0
		size := 0
		padCount := 0

		end := index

		h1 := h.block[index]
		h2 := h.block[index+1]
		h3 := h.block[index+2]
		h4 := h.block[index+3]
		h5 := h.block[index+4]
		//h6 := h.block[index+5]
		h7 := h.block[index+6]
		//h8 := h.block[index+7]

		// Skip the weird stuff that the code touches.
		if h1 == 0x00 && h2 == 0x02 && h3 == 0x01 && h4 == 0x01 && h5 == 0x00 {
			fmt.Printf("Skipping [0x%X] - [0x%X]...\n", index, index+279)
			index += 278
			previousEnd = index + 1
		}

		if h1 == 0x33 && h2 == 0x4F && h3 == 0x01 {
			fmt.Printf("Skipping [0x%X] - [0x%X]...\n", index, index+3)
			index += 3
			previousEnd = index + 1
		}

		if h1 == 0xA0 && h2 == 0xC0 {
			fmt.Printf("Skipping [0x%X] - [0x%X]...\n", index, index+5)
			index += 5
			previousEnd = index + 1
		}

		if h1 == 0x01 && h2 == 0x10 {
			fmt.Printf("Skipping [0x%X] - [0x%X]...\n", index, index+61)
			index += 61
			previousEnd = index + 1
		}

		if h1 == 0x33 && h2 == 0x03 && h3 == 0x70 {
			fmt.Printf("Skipping [0x%X] - [0x%X]...\n", index, index+216)
			index += 216
			previousEnd = index + 1
		}

		// skip 0xFF's
		if h1 == 0xFF && h2 == 0xFF {
			index++
			previousEnd = index + 1
			continue Loop
		}

		//fmt.Printf("Index: 0x%X	-	TESTING: H1: 0x%X	H2: 0x%X	H3: 0x%X	H4: 0x%X	H5: 0x%X	H7: 0x%X	H8: 0x%X\n\n", index, h1, h2, h3, h4, h5, h7, h8)

		// 4 Byte Match
		if h1 == 0x02 {
			//fmt.Printf("If #3 Index: 0x%X...\n", index)

			size = 4
			end = index + size + 1

			count++

			if previousEnd > 1 && previousEnd < index { // byte match
				fmt.Printf("###############################################################\n PREVIOUS END: 0x%X [%d bytes]\n MISSED: %X \n###############################################################\n\n\n\n", previousEnd, len(h.block[previousEnd:index]), h.block[previousEnd:index])
			}

			match := fmt.Sprintf(" 4 BYTES MATCH # %d	-	ADDRESS: 0x%X	END: 0x%X	[%d bytes]\n 0x %X\n", count, index, end, size, h.block[index:end])
			log(match, nil)

			addresses = append(addresses, index)

			previousEnd = end + padCount

			index = previousEnd - 1

			continue Loop
		}

		//fmt.Printf("HERE Index: 0x%X...\n", index)

		// Word Aligned Stuff
		if index%2 == 0 && (h1 == 0x00 || h1 == 0x40 || h1 == 0x60 || h1 == 0x80) && h2&0xF8 == 0x00 {
			//fmt.Printf("If #1 Index: 0x%X...\n", index)

			// 40 00 60 00 0C
			// A0 00 00 08 06 09 04 09 00
			// 00 02 00 02 16 0D 08 09 FE FE
			if ((h3 == 0x00 && h4&0xF8 == 0x00) || (h3&0x1F == 0x00 && h4 == 0x00)) && h5 > 0x00 { // 3D Maps!
				//fmt.Printf("If #2 Index: 0x%X...\n", index)
				dataType = "3D Map"

				cols = int(h.block[index+4]) + 1
				rows = int(h.block[index+5]) + 1

				size = cols * rows

				// two byte mode?
				if h2 == 0x01 {
					size *= 2
				}

				end = index + size + 8

				if (h2 == 0x02 || h2 == 0x04) && h7&0x01 == 0x01 {
					// Account for escaped 0x00
					for di := index + 8; di < end+padCount; di += 2 {
						if h.block[di] == 0x00 {
							padCount += 1
						}
					}
				}

				count++

				if previousEnd > 1 && previousEnd < index {
					fmt.Printf("###############################################################\n PREVIOUS END: 0x%X [%d bytes]\n MISSED: %X \n###############################################################\n\n\n\n", previousEnd, len(h.block[previousEnd:index]), h.block[previousEnd:index])
				}

				match := fmt.Sprintf(" 3D MATCH # %d	-	ADDRESS: 0x%X	END: 0x%X	%d x %d	[%d bytes]	[padded: %d]	Type: %s\n 0x %X \n", count, index, end, cols, rows, size, padCount, dataType, h.block[index:end+padCount])
				log(match, nil)

				addresses = append(addresses, index)

				previousEnd = end + padCount

				// Round up if not word aligned
				if previousEnd%2 != 0 {
					previousEnd++
				}

				index = previousEnd - 1

				continue Loop

				// 00 02 0D 09 60 50
			} else if h3&0xF0 == 0x00 && h3 > 0x00 { // Arrays!
				//00 00 00 02 00
				//fmt.Printf("If #3 Index: 0x%X...\n", index)
				dataType = "Array"

				size = int(h3) + 1

				// two byte mode?
				if h1 == 0x40 && h4&0x04 == 0x04 {
					fmt.Printf("If #4 Index: 0x%X...\n", index)
					size *= 2
				}

				end = index + size + 4

				// Account for escaped 0x0*
				if h2 == 0x04 && h4 == 0x05 {
					for di := index + 4; di < end+padCount-1; di++ {
						if di%2 == 0x00 && h.block[di]&0x0F == 0x00 {
							padCount += 1
						}
					}
				}

				// alternate pad 1
				if h1 == 0x00 && (h2 == 0x00 && h3 == 0x07 && h4 == 0x05) {
					for di := index + 4; di < end+padCount; di += 1 {
						if h.block[di] < 0x05 && padCount < size {
							padCount += 1
						}
					}
				}

				// alternate pad 2
				// 00 02 0D 09
				if h1 == 0x00 && (h2 == 0x02 && h3 == 0x0D && h4 == 0x09) {
					for di := index + 4; di < end+padCount; di++ {
						if h.block[di] < h2 && h.block[di] > 0x00 {
							padCount += 1
							di++
						}
					}
				}

				count++

				if previousEnd < index {
					fmt.Printf("###############################################################\n PREVIOUS END: 0x%X [%d bytes]\n MISSED: %X \n###############################################################\n\n\n\n", previousEnd, len(h.block[previousEnd:index]), h.block[previousEnd:index])
				}

				match := fmt.Sprintf(" ARRAY MATCH # %d	-	ADDRESS: 0x%X	END: 0x%X	[%d bytes]	[padded: %d]	Type: %s\n 0x %X\n", count, index, end+padCount, size, padCount, dataType, h.block[index:end+padCount])
				log(match, nil)

				addresses = append(addresses, index)

				previousEnd = end + padCount

				// Round up if not word aligned
				if previousEnd%2 != 0 {
					previousEnd++
				}

				index = previousEnd - 1

				continue Loop

			} else {
				index++

				continue Loop
			}

		}

	}

	for _, address := range addresses {
		fmt.Printf("0x%X, ", address)
	}

	fmt.Println("")

	return addresses, nil

}

func FindMatch(block []byte) [][]int {
	regex := string(0x00) + "[" + string(0x00) + "-" + string(0x05) + "]" + string(0x00) + "[" + string(0x00) + "-" + string(0x0F) + "]"

	re := regexp.MustCompile(regex)
	matches := re.FindAllIndex(block, -1)

	return matches
}

func printTable(width, height int, table []byte) {
	rows := [][]string{}

	log("TABLE 8", nil)

	for h := 0; h < height; h++ {
		row := []string{}
		//rowHex := []byte{}
		for w := 0; w < width; w++ {

			offset := h * width

			row = append(row, fmt.Sprintf("%d", table[w+offset]))
			//	rowHex = append(rowHex, table[w+offset])
		}
		//log(fmt.Sprintf("%X", rowHex), nil)
		rows = append(rows, row)
	}

	t := cli.NewTable(rows, &cli.TableOptions{
		Padding:      1,
		UseSeparator: true,
	})
	//t.SetHeader(collumns)
	fmt.Print("\n")

	fmt.Println(t.Render())

	fmt.Println("================================================================================================================")

	fmt.Print("\n\n\n")

}

func printTable16(width, height int, table []byte) {

	if width%2 == 0 {

		rows := [][]string{}

		log("TABLE 16", nil)

		for h := 0; h < height; h++ {
			row := []string{}
			//rowHex := []uint16{}

			for w := 0; w < width; w = w + 2 {

				offset := h * width

				//sixteen := (uint16(table[w+offset]) << 8) | uint16(table[w+offset+1])
				sixteen := (uint16(table[w+offset+1]) << 8) | uint16(table[w+offset])

				row = append(row, fmt.Sprintf("%d", sixteen))
				//rowHex = append(rowHex, sixteen)
			}
			//log(fmt.Sprintf("%X", rowHex), nil)
			rows = append(rows, row)
		}

		t := cli.NewTable(rows, &cli.TableOptions{
			Padding:      1,
			UseSeparator: true,
		})
		//t.SetHeader(collumns)
		fmt.Print("\n")

		fmt.Println(t.Render())

	}

	fmt.Print("\n\n\n")

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
