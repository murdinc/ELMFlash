package main

import (
	"./iso9141"
	"fmt"
	"github.com/murdinc/cli"
	"os"
)

// Main Function
////////////////..........
func main() {

	app := cli.NewApp()
	app.Name = "protege-hack"
	app.Usage = "Command Line Interface for programming the 3rd Generation Mazda Protege"
	app.Version = "1.0"
	app.Commands = []cli.Command{
		{
			Name:      "download",
			ShortName: "d",
			//Usage:     "Download BIN from ECU to computer",
			Action: func(c *cli.Context) {
				obd := iso9141.New()
				obd.DownloadBIN("DUMP")
			},
		},
		{
			Name:      "upload",
			ShortName: "u",
			//Usage:     "Upload BIN from computer to ECU",
			Action: func(c *cli.Context) {
				obd := iso9141.New()
				obd.UploadBIN()
			},
		},
		{
			Name:      "test",
			ShortName: "t",
			Action: func(c *cli.Context) {
				obd := iso9141.New()
				obd.Test()
			},
		},
		{
			Name:      "common",
			ShortName: "c",
			Action: func(c *cli.Context) {
				obd := iso9141.New()
				obd.CommonIdDump("COMMON_ID")
			},
		},
		{
			Name:      "ecuId",
			ShortName: "i",
			Action: func(c *cli.Context) {
				obd := iso9141.New()
				obd.EcuId()
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
