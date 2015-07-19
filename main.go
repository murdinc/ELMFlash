package main

import (
	"./disasm"
	"./hexstuff"
	"./iso9141"
	"fmt"
	"github.com/murdinc/cli"
	"os"
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
				obd := iso9141.New()
				obd.DownloadBIN("DOWNLOAD")
			},
		},
		{
			Name:        "dump",
			ShortName:   "du",
			Example:     "dump",
			Description: "Dump the calibration from the ECU without security mode (slow)",
			Action: func(c *cli.Context) {
				obd := iso9141.New()
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
			Action: func(c *cli.Context) {
				obd := iso9141.New()
				obd.UploadBIN(c.NamedArg("calibration"))
			},
		},
		{
			Name:      "tests",
			ShortName: "s",
			Action: func(c *cli.Context) {
				hs := hexstuff.New()
				hs.TestS()
			},
		},
		{
			Name:        "common",
			ShortName:   "c",
			Example:     "common",
			Description: "Crawls all Common ID's",
			Action: func(c *cli.Context) {
				obd := iso9141.New()
				obd.CommonIdDump("COMMON_ID")
			},
		},
		{
			Name:        "local",
			ShortName:   "l",
			Example:     "local",
			Description: "Crawls all Local ID's",
			Action: func(c *cli.Context) {
				obd := iso9141.New()
				obd.LocalIdDump("LOCAL_ID")
			},
		},
		{
			Name:        "ecuId",
			ShortName:   "i",
			Example:     "ecuId",
			Description: "Retrieve the ECU ID",
			Action: func(c *cli.Context) {
				obd := iso9141.New()
				obd.EcuId()
			},
		},
		{
			Name:      "test",
			ShortName: "t",
			Action: func(c *cli.Context) {
				hs := hexstuff.New()
				hs.Test()
			},
		},
		{
			Name:      "tests",
			ShortName: "s",
			Action: func(c *cli.Context) {
				hs := hexstuff.New()
				hs.TestS()
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
	}

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
