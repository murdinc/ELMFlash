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
	}

	app.Run(os.Args)
}

/*
	obd := iso9141.New()

	cmdResp, err := obd.Cmd("ATDP")
	if err != nil {
		log("ATDP", err)
	} else {
		log("Protocol - ["+cmdResp+"]", nil)
	}

	msg := []byte{0x22, 0x02, 0x00}
	msgResp, err := obd.Msg(msg)
	if err != nil {
		log("CMD 22 17", err)
	} else {
		fmt.Printf("Test Message response: %X\n", msgResp.Message)
	}

	msg = []byte{0x22, 0x17, 0x00}
	msgResp, err = obd.Msg(msg)
	if err != nil {
		log("CMD 22 17", err)
	} else {
		fmt.Printf("Test Message response: %X\n", msgResp.Message)
	}

	obd.EnableSecurity()

	obd.DownloadBlock()

*/

// Log Function
////////////////..........
func log(kind string, err error) {
	if err == nil {
		fmt.Printf("====> %s\n", kind)
	} else {
		fmt.Printf("[ERROR - %s]: %s\n", kind, err)
	}
}
