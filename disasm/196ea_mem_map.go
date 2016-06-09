package disasm

import "fmt"

type MemLocations []MemLocation

type MemLocation struct {
	Name        string
	Description string
	Start       int
	Stop        int
	Ignore      bool
}

var memMap = MemLocations{

	MemLocation{
		Name:        "Lower register file",
		Description: "Lower register file (stack pointer and CPU SFRs)",
		Start:       0x000000,
		Stop:        0x000019,
		Ignore:      false,
	},
	MemLocation{
		Name:        "Lower register file",
		Description: "Lower register file (general-purpose register RAM)",
		Start:       0x00001A,
		Stop:        0x0000FF,
		Ignore:      false,
	},
	MemLocation{
		Name:        "Upper register file",
		Description: "Upper register file (general-purpose register RAM)",
		Start:       0x000100,
		Stop:        0x0003FF,
		Ignore:      false,
	},
	MemLocation{
		Name:        "Internal code/data RAM",
		Description: "Internal code/data RAM (identically mapped into page FFH)",
		Start:       0x000400,
		Stop:        0x000FFF,
		Ignore:      false,
	},
	MemLocation{
		Name:        "External device",
		Description: "External device (memory or I/O) connected to address/data bus",
		Start:       0x001000,
		Stop:        0x001BFF,
		Ignore:      false,
	},
	MemLocation{
		Name:        "Peripheral special-function registers (SFRs)",
		Description: "Peripheral special-function registers (SFRs)",
		Start:       0x001C00,
		Stop:        0x001FDF,
		Ignore:      false,
	},
	MemLocation{
		Name:        "Memory-mapped special-function registers (SFRs)",
		Description: "Memory-mapped special-function registers (SFRs)",
		Start:       0x001FE0,
		Stop:        0x001FFB,
		Ignore:      false,
	},
	MemLocation{
		Name:        "Memory-mapped special-function registers (SFRs)",
		Description: "Memory-mapped special-function registers (SFRs); External memory if EA# is low; internal ROM if EA# is high. ",
		Start:       0x001FFC,
		Stop:        0x001FFF,
		Ignore:      false,
	},
	MemLocation{
		Name:        "External device",
		Description: "External device (memory or I/O) connected to address/data bus",
		Start:       0x002000,
		Stop:        0x0023FF,
		Ignore:      false,
	},
	MemLocation{
		Name:        "Internal ROM or External Memory",
		Description: "A copy of internal ROM (FF2400–FF3FFFH) if CCB1.2=0 External memory if CCB1.2=1 ",
		Start:       0x002400,
		Stop:        0x003FFF,
		Ignore:      false,
	},
	MemLocation{
		Name:        "External device",
		Description: "External device (memory or I/O) connected to address/data bus",
		Start:       0x004000,
		Stop:        0xFFFFF,
		Ignore:      false,
	},

	MemLocation{ // NOT SURE ABOUT THIS ONE
		Name:        "Overlaid memory",
		Description: "Overlaid memory (reserved for future microcontrollers); locations xF0000–xF03FFH are reserved for in-circuit emulators",
		Start:       0x100000,
		Stop:        0x16FFFF, // ????
		Ignore:      false,
	},

	MemLocation{
		Name:        "Reserved",
		Description: "Reserved for in-circuit emulators",
		Start:       0x170000,
		Stop:        0x1703FF,
		Ignore:      false,
	},
	MemLocation{
		Name:        "Internal code/data RAM ",
		Description: "Internal code/data RAM (identically mapped from page 00H)",
		Start:       0x170400,
		Stop:        0x170FFF,
		Ignore:      true,
	},
	MemLocation{
		Name:        "External device",
		Description: "External device (memory or I/O) connected to address/data bus",
		Start:       0x171000,
		Stop:        0x171FFF,
		Ignore:      false,
	},
	MemLocation{
		Name:        "Special-purpose memory",
		Description: "Special-purpose memory (CCBs, interrupt vectors, PTS vectors)",
		Start:       0x172000,
		Stop:        0x17207F,
		Ignore:      false,
	},
	MemLocation{
		Name:        "Program Start",
		Description: "After reset, the first instruction is fetched from 0x172080.",
		Start:       0x172080,
		Stop:        0x1720BF,
		Ignore:      false,
	},
	MemLocation{
		Name:        "Special-purpose memory",
		Description: "Special-purpose memory (PIH vectors)",
		Start:       0x1720C0,
		Stop:        0x17213F,
		Ignore:      false,
	},
	MemLocation{
		Name:        "Program memory",
		Description: "Program memory",
		Start:       0x172140,
		Stop:        0x1723FF,
		Ignore:      false,
	},
	MemLocation{
		Name:        "Program memory",
		Description: "Program memory; can be mapped into page 00H (CCB1.2 = 1)",
		Start:       0x172400,
		Stop:        0x173FFF,
		Ignore:      false,
	},
	MemLocation{
		Name:        "External device",
		Description: "External device (memory or I/O) connected to address/data bus",
		Start:       0x174000,
		Stop:        0x17FFFF,
		Ignore:      false,
	},
}

func (h *DisAsm) GetMemoryMap() error {

	h.memStarts = make(map[int]string) // Starts of memory map Locations
	h.memStops = make(map[int]string)  // Ends of memory map Locations
	h.skip = make(map[int]int)         // Start And Stop locations for places to skip

	for _, memLoc := range memMap {
		h.memStarts[memLoc.Start] = memLoc.Name
		h.memStops[memLoc.Stop] = memLoc.Name

		if memLoc.Ignore {
			h.skip[memLoc.Start] = memLoc.Stop
		}
	}

	return nil
}

func (h *DisAsm) doMemoryMap(adr int) int {
	skip := 0

	//Print the end of a location
	if h.memStops[adr-1] != "" {
		location := addSpaces(fmt.Sprintf("\n\n **\n ** [0x%X] END OF %s \n **\n **\n *********************************************\n\n", adr-1, h.memStops[adr-1]), 80)
		log(location, nil)
	}

	// Print out our Memory Location Start Block
	if h.memStarts[adr] != "" {
		location := addSpaces(fmt.Sprintf("\n\n *********************************************\n **\n ** [0x%X] START OF %s \n **\n **\n\n", adr, h.memStarts[adr]), 80)
		log(location, nil)
		if h.skip[adr] != 0 {
			skip = (h.skip[adr] - adr) + 1
			location := addSpaces(fmt.Sprintf("** SKIPPING %d BYTES \n **\n", skip), 80)
			log(location, nil)
			location = addSpaces(fmt.Sprintf("\n\n **\n ** [0x%X] END OF %s \n **\n **\n *********************************************\n\n", adr+skip-1, h.memStops[adr+skip-1]), 80)
			log(location, nil)
		}

	}

	return skip

}
