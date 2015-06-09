package hexstuff

import (
	//"bufio"
	//"encoding/hex"
	//"encoding/binary"
	//"errors"
	"fmt"
	//"github.com/cheggaaa/pb"
	//serial "github.com/huin/goserial"
	//"io"
	//	"io/ioutil"
	//"math/rand"
	"os"
	//"strings"
	//"time"
	"github.com/murdinc/cli"
	"regexp"
)

// App constants
////////////////..........
const debug = false

func New() *HexStuff {
	controller := new(HexStuff)
	controller.location = "MSP.BIN"
	return controller
}

// Connection represents an OBD-II serial connection
type HexStuff struct {
	location string
}

func (h *HexStuff) TestS() error {

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

	regex := string(0x00) + string(0x02) //"[" + string(0x00) + "-" + string(0x01) + "]" + "[" + string(0x01) + "-" + string(0x10) + "]"
	re := regexp.MustCompile(regex)
	matches := re.FindAllIndex(block, -1)

	listCount := 0
	tableCount := 0
	missedCount := 0

	previous := 0x0000

Loop:
	for _, i := range matches {

		// Break if we are over the 0x10000 limit
		if i[0] > 0x10000 {
			break
		}

		// Set the index of the block file
		index := i[0]

		// Skip if withing our previous block of data
		if index < previous {
			//log(fmt.Sprintf("Index [%X] is within previous item ending at [%X], skipping..\n", index, previous), nil)
			continue
		}

		// Print whatever we missed since the last item
		if previous != 0x0 && index != previous {
			log(fmt.Sprintf("MISSED!! PREVIOUS END: [%X] NEXT INDEX: [%X] SIZE: %d", previous, index, len(block[previous:index])), nil)
			log(fmt.Sprintf("********************************** \n %v \n %X \n**********************************\n", block[previous:index], block[previous:index]), nil)
			missedCount++
		}

		// Parts we care about right now
		thirdByte := block[index+2]
		fourthByte := block[index+3]

		//log(fmt.Sprintf("THIRD: %X FOURTH %X ", thirdByte, fourthByte), nil)

		// Third Byte is not 0x00, so must be a list?
		if !(thirdByte == 0x00 && fourthByte == 0x2) &&
			!(thirdByte == 0x00 && fourthByte == 0x5) {

			for b := index + 2; b < len(block); b++ {
				if (block[b] == 0x00 && block[b+1] == 0x02) ||
					//(block[b] == 0x00 && block[b+1] == 0x00) ||
					//(block[b] == 0x01 && block[b+1] == 0x02) ||
					//(block[b] == 0x40 && block[b+1] == 0x00) ||
					b >= 0x8000 { // && b%2 == 0 {

					listCount++
					size := len(block[index:b])
					match := fmt.Sprintf("LIST MATCH: %v   ADDRESS: 0x%X  END: 0x%X SIZE: [%d]", block[index:b], index, b, size)
					log(match, nil)
					match = fmt.Sprintf("LIST MATCH: %X   ADDRESS: 0x%X  END: 0x%X SIZE: [%d]", block[index:b], index, b, size)
					log(match, nil)

					if len(block[index:b]) > 6 && len(block[index+4:b]) == int(block[index+2]+1) {
						log(fmt.Sprintf("LIST MATCH! [%d]", len(block[index+4:b])), nil)

					}

					if size < 11 {
						previous = b
						fmt.Println("================================================================================================================")
						fmt.Println("================================================================================================================")
						fmt.Println("\n")
						continue Loop
					}

					log(fmt.Sprintf("SIZECHECK SIZE: [%d] LENGTH: [%d]", block[index+2]+1, len(block[index+4:b])), nil)

					length := int(block[index+2] + 1)
					if block[index+3] == 0x09 || block[index+3] == 0x08 || block[index+3] == 0x01 {
						singleStart := index + 4
						fmt.Println("\n")
						log(fmt.Sprintf("NOT TABLE: [%v]", block[singleStart:singleStart+length]), nil)
						if length+4 < size {
							log(fmt.Sprintf("EXTRA: [%v]", block[singleStart+length:b]), nil)
						}
					} else {
						height := block[index+3] + 1
						heightCount := 0
						for p := 6; p < size; p = p + length {
							heightCount++
							end := index + p + length
							if end > b {
								end = b
							}
							log(fmt.Sprintf("TABLE: [%v]		WIDTH: [%d]		HEIGHT: [%d]		COUNT: [%d]", block[index+p:end], length, height, heightCount), nil)
						}

					}

					previous = b
					fmt.Println("================================================================================================================")
					fmt.Println("================================================================================================================")
					fmt.Println("\n")
					continue Loop
				}

			}

		}

		// Third byte is 00 so must be a table?
		if (thirdByte == 0x00 && fourthByte == 0x02) ||
			(thirdByte == 0x00 && fourthByte == 0x05) {

			m1 := block[index+1]
			m2 := block[index+3]

			height := int(block[index+4]) + 1
			width := int(block[index+5]) + 1

			tableCount++

			size := width * height
			start := index + 8
			end := start + size
			previous = end

			match := fmt.Sprintf("TABLE MATCH: %X   ADDRESS: 0x%X  END: 0x%X  SIZE: %d x %d     [%d]    M1: %d      M2: %d      L1: 0x%X	L2: 0x%X", block[index:index+8], index, end, width, height, size, m1, m2, block[index+6], block[index+7])
			log(match, nil)

			//log(fmt.Sprintf("%X", block[index:end]), nil)

			fmt.Println("\n")

			//printTable16(width, height, block[start:end])
			printTable(width, height, block[start:end])

		}

	}

	log(fmt.Sprintf("LIST COUNT: %d TABLE COUNT: %d MISSED COUNT: %d", listCount, tableCount, missedCount), nil)

	return nil

}

// TEST
func (h *HexStuff) Test() error {

	// Pull in the Calibration file
	f, err := os.Open("calibrations/MSP.BIN")
	if err != nil {
		log("BIN - Error opening file", err)
		return err
	}

	// Make a 1024 byte buffer
	//block := make([]byte, 1024)
	block := make([]byte, 524288)

	// Read 1024 bytes
	n, err := f.Read(block)
	if err != nil {
		log("UploadBIN - Error reading calibration", err)
		return err
	}
	dbg(fmt.Sprintf("BIN - reading %d bytes.", n), nil)

	//fmt.Print(block)

	//regex := string([]byte{0x00, 0x04, 0x00, 0x02})
	regex := string(0x00) + "[" + string(0x01) + "-" + string(0x04) + "]" + string(0x00) + "[" + string(0x02) + "-" + string(0x09) + "]"
	re := regexp.MustCompile(regex)
	matches := re.FindAllIndex(block, -1)

	count := 0

	sizes := make(map[string]int)

	previous := 0x0000

	for _, i := range matches {
		if i[0]%2 == 0 && i[0] >= previous && i[0] < 0x10000 {
			index := i[0]

			if block[index+6] == 0x08 || block[index+6] == 0x09 || block[index+6] == 0x07 || block[index+6] == 0x06 || block[index+6] == 0x10 {

				height := int(block[index+4]) + 1
				width := int(block[index+5]) + 1

				m1 := block[index+1]
				m2 := block[index+3]

				count++

				size := width * height
				start := index + 8
				end := start + size
				previous = end

				sixteen := (int16(block[index+7]) << 8) | int16(block[index+6])

				match := fmt.Sprintf("MATCH: 0x%X   ADDRESS: 0x%X  END: 0x%X  SIZE: %d x %d		[%d]	M1: %d      M2: %d      L: 0x%X	L: %d", block[index:index+8], index, end, width, height, size, m1, m2, block[index+6:index+8], sixteen)
				log(match, nil)

				//log(fmt.Sprintf("%X", block[index:end]), nil)

				//fmt.Println("\n")

				printTable16(width, height, block[start:end])
				printTable(width, height, block[start:end])

				sizeName := fmt.Sprintf("%d	x	%d", width, height)
				var sizeCount int
				if sizes[sizeName] < 1 {
					sizeCount = 1
				} else {
					sizeCount = sizes[sizeName]
					sizeCount++
				}
				sizes[sizeName] = sizeCount

			}

		}
	}

	log(fmt.Sprintf("COUNT: %d ", count), nil)

	for size, count := range sizes {
		log(fmt.Sprintf("SIZE: %s		COUNT: %d", size, count), nil)
	}

	return nil

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
