package main

import (
	"fmt"
	"os"

	"github.com/murdinc/ELMFlash/calibrate"
	"github.com/murdinc/ELMFlash/compare"
	"github.com/murdinc/ELMFlash/disasm"
	"github.com/murdinc/ELMFlash/hexstuff"
	"github.com/murdinc/ELMFlash/iso9141"
	"github.com/murdinc/ELMFlash/j3"
	"github.com/murdinc/legacy-cli"
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
			Name:        "j3",
			ShortName:   "j3",
			Example:     "j3",
			Description: "J3 Firmata Test",
			Action: func(c *cli.Context) {
				//connection := j3.New()
				//connection.Run()
				j3.Test()

			},
		},
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
			Description: "Map Test 2",
			Arguments: []cli.Argument{
				cli.Argument{Name: "calibration", Usage: "maptest2 msp", Description: "The name of the calibration to run map tests on", Optional: false},
			},

			Action: func(c *cli.Context) {
				hs := hexstuff.New()
				hs.TestM2(c.NamedArg("calibration"))
			},
		},
		{
			Name:        "maptest3",
			ShortName:   "m3",
			Example:     "maptest2",
			Description: "Map Test 3",
			Arguments: []cli.Argument{
				cli.Argument{Name: "calibration", Usage: "maptest2 msp", Description: "The name of the calibration to run map tests on", Optional: false},
			},

			Action: func(c *cli.Context) {
				hs := hexstuff.New()
				hs.TestM3(c.NamedArg("calibration"))
			},
		},
		{
			Name:        "compare",
			ShortName:   "cmp",
			Example:     "compare mp3 mp3-2",
			Description: "Compare",
			Arguments: []cli.Argument{
				cli.Argument{Name: "pre1", Usage: "compare pre mp3 pre2 mp3x2", Description: "The name of the first pre calibrations to compare", Optional: false},
				cli.Argument{Name: "calibration1", Usage: "compare pre mp3 pre2 mp3x2", Description: "The name of the first calibrations to compare", Optional: false},
				cli.Argument{Name: "pre2", Usage: "compare pre mp3 pre2 mp3x2", Description: "The name of the second pre calibrations to compare", Optional: false},
				cli.Argument{Name: "calibration2", Usage: "compare pre mp3 pre2 mp3x2", Description: "The name of the second calibrations to compare", Optional: false},
			},

			Action: func(c *cli.Context) {
				cmp := compare.New(c.NamedArg("pre1"), c.NamedArg("calibration1"), c.NamedArg("pre2"), c.NamedArg("calibration2"))
				cmp.Compare()
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
				d := disasm.New(c.NamedArg("calibration"))
				d.DisAsm()
			},
		},
		{
			Name:        "interrupt",
			ShortName:   "int",
			Example:     "interrupt",
			Description: "List Interrupt Vector Addresses in Calibration File",
			Arguments: []cli.Argument{
				cli.Argument{Name: "calibration", Usage: "interrupt msp", Description: "The name of the calibration to scan", Optional: false},
			},
			Action: func(c *cli.Context) {
				d := disasm.New(c.NamedArg("calibration"))
				d.GetInterrupts()
			},
		},
		{
			Name:        "calibrate",
			ShortName:   "cal",
			Example:     "interrupt",
			Description: "Calibrate",
			Arguments: []cli.Argument{
				cli.Argument{Name: "calibration", Usage: "calibrate msp", Description: "The name of the calibration to edit", Optional: false},
			},
			Action: func(c *cli.Context) {
				calibrate.Calibrate(true)
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
