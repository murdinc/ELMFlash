package calibrate

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/goware/cors"
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

	r := chi.NewRouter()

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowCredentials: true,
	})

	r.Use(cors.Handler)
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.Logger)

	// Static Asset Folders
	FileServer(r, "/js", http.Dir(guiLocation))
	FileServer(r, "/css", http.Dir(guiLocation))
	FileServer(r, "/fonts", http.Dir(guiLocation))
	FileServer(r, "/static", http.Dir(guiLocation))

	// Index and Dashboard
	r.Get("/", index)
	r.Get("/table/:type/:address", getTable)

	if !devMode {
		webbrowser.Open("http://localhost:8080/") // TODO race condition?
	}

	http.ListenAndServe(":8080", r)

}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {

	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.FileServer(root)

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
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

func index(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseFiles("calibrate/ui/templates/calibrate.html")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	payload := make(map[string]interface{})

	var tables []Table

	calName := "mp3"
	calibration := New(calName)

	hs, _ := hexstuff.New(calName)

	addresses, _ := hs.TestM1()

	for _, address := range addresses {
		table := calibration.GetTable(address)

		tables = append(tables, *table)

	}

	payload["Tables"] = tables
	payload["Errors"] = nil

	/*ctx.Render("calibrate.html", Content{Title: "Calibrate", Payload: payload, Calibration: "MSP/MP3", RenderLayout: true})*/
	w.Header().Set("Content-Type", "text/html")
	t.Execute(w, payload)
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

func getTable(w http.ResponseWriter, r *http.Request) {
	payload := make(map[string]interface{})

	calibration := New("mp3")
	table := calibration.GetTable(0x10AA60)

	payload["Table"] = table
	payload["Errors"] = nil

	/*ctx.JSONP(iris.StatusOK, "callbackName", Content{Title: "Calibrate", Payload: payload, Calibration: "MSP/MP3", RenderLayout: true})*/

	js, err := json.Marshal(Content{Title: "Calibrate", Payload: payload, Calibration: "MSP/MP3", RenderLayout: true})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

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
