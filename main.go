package main

import (
	"fmt"
	"os"

	"./disasm"
	"./hexstuff"
	"./iso9141"
	"github.com/murdinc/cli"
)

// Main Function
////////////////..........
func main() {

	app := cli.NewApp()
	app.Name = "ELMFlash"
	app.Usage = "Command Line Interface for programming the 3rd Generation Mazda Protege"
	app.Version = "1.0"
	app.Commands = []cli.Command{
		{
			Name:        "download",
			ShortName:   "d",
			Example:     "download",
			Description: "Download the calibration from the ECU",
			Action: func(c *cli.Context) {
				obd := iso9141.New(false)
				obd.DownloadBIN("DOWNLOAD")
			},
		},
		{
			Name:        "dump",
			ShortName:   "du",
			Example:     "dump",
			Description: "Dump the calibration from the ECU without security mode (slow)",
			Action: func(c *cli.Context) {
				obd := iso9141.New(false)
				obd.DumpBIN("DUMP")
			},
		},
		{
			Name:        "upload",
			ShortName:   "u",
			Example:     "upload msp",
			Description: "Upload a calibration to the ECU",
			Arguments: []cli.Argument{
				cli.Argument{Name: "calibration", Usage: "upload msp", Description: "The name of the calibration to upload", Optional: false},
			},
			Flags: []cli.Flag{
				cli.BoolFlag{Name: "test", Usage: "Test upload"},
			},
			Action: func(c *cli.Context) {

				obd := iso9141.New(c.Bool("test"))
				obd.UploadBIN(c.NamedArg("calibration"))
			},
		},
		{
			Name:        "common",
			ShortName:   "c",
			Example:     "common",
			Description: "Crawls all Common ID's",
			Action: func(c *cli.Context) {
				obd := iso9141.New(false)
				obd.CommonIdDump("COMMON_ID")
			},
		},
		{
			Name:        "local",
			ShortName:   "l",
			Example:     "local",
			Description: "Crawls all Local ID's",
			Action: func(c *cli.Context) {
				obd := iso9141.New(false)
				obd.LocalIdDump("LOCAL_ID")
			},
		},
		{
			Name:        "ecuId",
			ShortName:   "i",
			Example:     "ecuId",
			Description: "Retrieve the ECU ID",
			Action: func(c *cli.Context) {
				obd := iso9141.New(false)
				obd.EcuId()
			},
		},
		{
			Name:        "maptest1",
			ShortName:   "m1",
			Example:     "maptest1",
			Description: "Map Test 1",
			Arguments: []cli.Argument{
				cli.Argument{Name: "calibration", Usage: "maptest1 msp", Description: "The name of the calibration to run map tests on", Optional: false},
			},

			Action: func(c *cli.Context) {
				hs := hexstuff.New()
				hs.TestM1(c.NamedArg("calibration"))
			},
		},
		{
			Name:        "maptest2",
			ShortName:   "m2",
			Example:     "maptest2",
			Description: "Map Test 1",
			Arguments: []cli.Argument{
				cli.Argument{Name: "calibration", Usage: "maptest2 msp", Description: "The name of the calibration to run map tests on", Optional: false},
			},

			Action: func(c *cli.Context) {
				hs := hexstuff.New()
				hs.TestM2(c.NamedArg("calibration"))
			},
		},
		{
			Name:        "disasm",
			ShortName:   "x",
			Example:     "disasm",
			Description: "Disassemble Calibration File",
			Arguments: []cli.Argument{
				cli.Argument{Name: "calibration", Usage: "disasm msp", Description: "The name of the calibration to disassemble", Optional: false},
			},
			Action: func(c *cli.Context) {
				d := disasm.New()
				d.DisAsm(c.NamedArg("calibration"))
			},
		},
		{
			Name:        "vDisasm",
			ShortName:   "v",
			Example:     "vDisasm",
			Description: "Disassemble Calibration File - Verbose",
			Arguments: []cli.Argument{
				cli.Argument{Name: "calibration", Usage: "vDisasm msp", Description: "The name of the calibration to disassemble", Optional: false},
			},
			Action: func(c *cli.Context) {
				d := disasm.New()
				d.VDisAsm(c.NamedArg("calibration"))
			},
		},
	}

	log("ELMFlash - v1.0", nil)
	log("Created by Ahmad A.", nil)
	log("© MVRD INDUSTRIES 2015", nil)
	log("Not for commercial use", nil)
	app.Run(os.Args)
}

// Log Function
////////////////..........
func log(kind string, err error) {
	if err == nil {
		fmt.Printf("====> %s\n", kind)
	} else {
		fmt.Printf("[ERROR - %s]: %s\n", kind, err)
	}
}
