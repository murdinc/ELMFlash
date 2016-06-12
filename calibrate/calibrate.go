package calibrate

import (
	"fmt"
	"os"
	"regexp"

	"github.com/kataras/iris"
	"github.com/murdinc/ELMFlash/hexstuff"
	"github.com/toqueteos/webbrowser"
)

// App constants
////////////////..........
const debug = false

type Calibration struct {
	block []byte
}

type Content struct {
	Title        string
	Calibration  string
	Type         string
	Payload      map[string]interface{}
	Configs      interface{}
	AZList       []string
	RenderLayout bool
	Errors       []string
	ClassFormURL string
}

var calibrations = map[string]string{
	"msp":   "MSP.BIN",
	"mp3":   "MP3.BIN",
	"mp3x2": "MP3x2.BIN",
	"p5":    "P5.BIN",
	"pre":   "PRE2.BIN",
}

func New(calName string) *Calibration {
	controller := new(Calibration)

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

func Calibrate(devMode bool) {

	guiLocation := "calibrate/ui/"

	api := iris.New()

	// Template Configuration
	api.Config().Render.Template.Directory = guiLocation
	api.Config().Render.Template.Layout = "templates/layout.html"

	// Static Asset Folders
	api.StaticWeb("/js", guiLocation, 0)
	api.StaticWeb("/css", guiLocation, 0)
	api.StaticWeb("/fonts", guiLocation, 0)
	api.StaticWeb("/static", guiLocation, 0)

	// Index and Dashboard
	api.Get("/", index)
	api.Get("/table/:type/:address", getTable)

	if !devMode {
		webbrowser.Open("http://localhost:8080/") // TODO race condition?
	}

	api.Listen(":8080") // TODO optionally configurable port #

}

func GetMaps(block []byte) []int {
	regex := string(0x00) + "[" + string(0x00) + "-" + string(0x05) + "]" + string(0x00) + "[" + string(0x00) + "-" + string(0x0F) + "]"
	re := regexp.MustCompile(regex)
	matches := re.FindAllIndex(block, -1)

	count := 0

	sizes := make(map[string]int)

	previous := 0x0000

	var addresses []int

	for _, i := range matches {

		index := i[0]

		height := int(block[index+4]) + 1
		width := int(block[index+5]) + 1

		if index%2 == 0 && index >= previous && index > 0x108000 && index < 0x118000 && height > 1 {

			//if block[index+6] == 0x08 || block[index+6] == 0x09 || block[index+6] == 0x07 || block[index+6] == 0x06 || block[index+6] == 0x10 {

			h2 := block[index+1] + 1
			h4 := block[index+3] + 1

			count++

			size := width * height
			start := index + 8
			end := start + size
			previous = end + 1

			sixteen := (int(block[index+7]) << 8) | int(block[index+6])

			match := fmt.Sprintf(" MATCH: -1 0x%X 0 0x%X +1 0x%X  ADDRESS: 0x%X  END: 0x%X  SIZE: %d x %d		[%d]	H2: %d      H4: %d      L: 0x%X	L: %d	|	START #%X | ROWS: %d x COLS: %d", block[index-8:index], block[index:index+8], block[index+8:index+16], index, end, width, height, size, h2, h4, block[index+6:index+8], sixteen, start, height, width)
			log(match, nil)

			//log(fmt.Sprintf("%X", block[index:end]), nil)

			//fmt.Println("\n")

			//printTable16(width, height, block[start:end])
			//printTable(width, height, block[start:end])

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

			//}
		}
	}

	return addresses
}

func index(ctx *iris.Context) {
	payload := make(map[string]interface{})

	var tables []Table

	// msp
	//calibration := New("msp")
	//addresses := []int{0x108F7E, 0x1090C8, 0x109212, 0x1093F4, 0x10953E, 0x109770, 0x109876, 0x109A2E, 0x109A98, 0x109B64, 0x109C88, 0x109CDC, 0x109D1C, 0x109F30, 0x10A2B8, 0x10A440, 0x10AB1E, 0x10AC68, 0x10ADB2, 0x10AEFC, 0x10B046, 0x10B20A, 0x10B218, 0x10B226, 0x10B234, 0x10B242, 0x10B250, 0x10B25E, 0x10B26C, 0x10B27A, 0x10B288, 0x10B296, 0x10B2A4, 0x10B2B2, 0x10B2C0, 0x10B2CE, 0x10B2DC, 0x10B2EA, 0x10B316, 0x10B48E, 0x10B618, 0x10B644, 0x10B680, 0x10B716, 0x10B752, 0x10B75E, 0x10B76A, 0x10B776, 0x10B872, 0x10B944, 0x10BE36, 0x10C14A, 0x10C1EE, 0x10C32E, 0x10C46E, 0x10C5AE, 0x10C6F6, 0x10C836, 0x10CA5A, 0x10CB5A, 0x10CC22, 0x10CCEA, 0x10CDB2, 0x10CE22, 0x10CE92, 0x10CF02, 0x10CF72, 0x10CFE2, 0x10D052, 0x10D0C2, 0x10D396, 0x10D3C0, 0x10D45E, 0x10D52E, 0x10D814, 0x10E204, 0x10E2CE, 0x10E3BA, 0x10E68C, 0x10E6EE, 0x10E766, 0x10E7DE, 0x10EA54, 0x10EDCA, 0x10F0FE, 0x10F158, 0x10F634, 0x10F6FC}

	// mp3
	calibration := New("mp3")
	//addresses := []int{0x10AB3A, 0x10A586, 0x10A660, 0x10A73A, 0x10A814, 0x10A986, 0x10AA60, 0x10ABEA, 0x10AC98, 0x10ADF8, 0x10AE62, 0x10AF2E, 0x10B042, 0x10B0E6, 0x10BA60, 0x10BB1E, 0x10BCDE, 0x10BDB8, 0x10BE92, 0x10BF6C, 0x10C046, 0x10C120, 0x10C1FA, 0x10C2D4, 0x10C3AE, 0x10C60E, 0x10C706, 0x10C7FE, 0x10CF5A, 0x10D3E0, 0x10D450, 0x10D528, 0x10D600, 0x10D6D8, 0x10D7B8, 0x10D890, 0x10DA84, 0x10DB04, 0x10DB84, 0x10DC04, 0x10E11E, 0x10E196, 0x10E20E, 0x10E6F8, 0x10E7C2, 0x10E876, 0x10EB20, 0x10F346, 0x10F55C, 0x10F624}
	//addresses := []int{0x10A444, 0x10A470, 0x10A660, 0x10A73A, 0x10A814, 0x10A986, 0x10AA60, 0x10AB3A, 0x10ABEA, 0x10AC98, 0x10ADF8, 0x10AE62, 0x10AF2E, 0x10B0E6, 0x10B1D4, 0x10B27E, 0x10B306, 0x10B4AA, 0x10E11E, 0x10E196, 0x10E20E, 0x10E6F8, 0x10E7C2, 0x10E876, 0x10E902, 0x10EB20, 0x10ED8A, 0x10EE68, 0x10EF9A, 0x10EFCA, 0x10EFFA, 0x10F346, 0x10F55C, 0x10F624, 0x10F782, 0x10FF92, 0x10FFFC}

	//addresses := GetMaps(calibration.block)

	hs := hexstuff.New()

	addresses, _ := hs.TestM3("mp3")

	for _, address := range addresses {
		table := calibration.GetTable(address)

		tables = append(tables, *table)

	}

	payload["Tables"] = tables
	payload["Errors"] = nil

	ctx.Render("templates/calibrate.html", Content{Title: "Calibrate", Payload: payload, Calibration: "MSP/MP3", RenderLayout: true})
}

type Table struct {
	Address    int
	AddressStr string
	EndStr     string
	Height     int
	Width      int

	H1    int
	H1Str string
	H2    int
	H2Str string
	H3    int
	H3Str string
	H4    int
	H4Str string
	H5    int
	H5Str string
	H6    int
	H6Str string
	H7    int
	H7Str string
	H8    int
	H8Str string

	Size  int
	Start int
	End   int
	Data  []int
}

func getTable(ctx *iris.Context) {
	payload := make(map[string]interface{})

	calibration := New("mp3")
	table := calibration.GetTable(0x10AA60)

	payload["Table"] = table
	payload["Errors"] = nil

	ctx.JSONP(iris.StatusOK, "callbackName", Content{Title: "Calibrate", Payload: payload, Calibration: "MSP/MP3", RenderLayout: true})
}

// 0x10AA60
func (c *Calibration) GetTable(index int) *Table {

	h1 := int(c.block[index])
	h1Str := fmt.Sprintf("0x%.2X", c.block[index])

	h2 := int(c.block[index+1])
	h2Str := fmt.Sprintf("0x%.2X", c.block[index+1])

	h3 := int(c.block[index+2])
	h3Str := fmt.Sprintf("0x%.2X", c.block[index+2])

	h4 := int(c.block[index+3])
	h4Str := fmt.Sprintf("0x%.2X", c.block[index+3])

	h5 := int(c.block[index+4]) // rows
	h5Str := fmt.Sprintf("0x%.2X", c.block[index+4])

	h6 := int(c.block[index+5]) // cols
	h6Str := fmt.Sprintf("0x%.2X", c.block[index+5])

	h7 := int(c.block[index+6])
	h7Str := fmt.Sprintf("0x%.2X", c.block[index+6])

	h8 := int(c.block[index+7])
	h8Str := fmt.Sprintf("0x%.2X", c.block[index+7])

	height := int(c.block[index+4]) + 1
	width := int(c.block[index+5]) + 1

	// axe flip™
	if h5&0x10 == 0x10 {
		width, height = height, width
	}

	size := width * height
	start := index + 8
	end := start + size

	table := new(Table)

	data := make([]int, size)

	for i := range data {

		val := int(c.block[start+i])
		/*
			if h7_1 == 1 && i%2 == 0 {
				val = int(c.block[start+i]) >> 4
				val2 := int(c.block[start+i]) & 0xF
				data[i] = val
				data[i+1] = val2

			} else if h7_1 == 0 {
				data[i] = val
			}
		*/
		data[i] = val
	}

	table.H1 = h1
	table.H1Str = h1Str

	table.H2 = h2
	table.H2Str = h2Str

	table.H3 = h3
	table.H3Str = h3Str

	table.H4 = h4
	table.H4Str = h4Str

	table.H5 = h5
	table.H5Str = h5Str

	table.H6 = h6
	table.H6Str = h6Str

	table.H7 = h7
	table.H7Str = h7Str

	table.H8 = h8
	table.H8Str = h8Str

	table.Size = size
	table.Start = start
	table.End = end

	table.Address = index
	table.AddressStr = fmt.Sprintf("0x%.6X", index)
	table.EndStr = fmt.Sprintf("0x%.6X", end)
	table.Height = height
	table.Width = width
	table.Data = data

	return table

}

func log(kind string, err error) {
	if err == nil {
		fmt.Printf(" %s\n", kind)
	} else {
		fmt.Printf("[ERROR - %s]: %s\n", kind, err)
	}
}
