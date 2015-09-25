# ELMFlash
ELMFlash is a CLI tool for reverse engineering the ECU on my car. It is written in Go, and is not yet complete. 

This tool utilizes an ELM 327 chip for communicating to the ECU through ISO9141 over the OBD-II port. These are readily available, but I would suggest trying to stay away from the ones with bootleg chips - as they have questionable circuit designs. The MCU on this car is a Intel 196EA variant (manufacturer proprietary) and memory is a vanilla Intel 28F400 Flash chip. 

Much of this was built by sniffing the packets being sent by the OEM supplied ECU flashing tool, and comparing that to datasheets for the specific OBD protocol. The disassembly and pseudo-code output was built by referencing the datasheets for the 196 and making (hopefully) informed assumptions. 

Currently I am trying to make sense of the disassembly and make that output more verbose. I am using a desk rig for testing that includes an electronic engine simulator (JimStim), and a modified ECU with cold-swappable Flash chips. 

Current capabilities: 
* Enters Security Mode
* Download the entire memory address block
* Upload a new tune (currently broken, but probably just an issue with the adjusted offset)
* Scan all Common ID's and Local ID's 
* Disassemble BIN tunes
* Tries to identify patterns of hex that seem to represent Map/Table data. 
