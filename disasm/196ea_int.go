package disasm

/*
	This microcontroller’s flexible interrupt-handling system has three main components:
	- The programmable interrupt controller
	- The peripheral transaction server (PTS)
	- The peripheral interrupt handlers (PIHs).

	The interrupt control circuitry within a microcontroller permits real-time events to control pro- gram flow.
	When an event generates an interrupt, the microcontroller suspends the execution of current instructions while it performs some service in response to the interrupt.
	When the inter- rupt is serviced, program execution resumes at the point where the interrupt occurred.
	An internal peripheral, an external signal, or an instruction can generate an interrupt request.
	In the simplest case, the microcontroller receives the request, performs the service, and returns to the task that was interrupted.

	The interrupt sources fall into two categories:
	- The unimplemented opcode, software trap, and NMI interrupt sources are always enabled.
	- All other sources can be individually enabled.

	Interrupts that go through the interrupt controller are serviced by interrupt service routines that you provide.
	The lower 16 bits of the addresses of these interrupt service routines are stored in the upper and lower interrupt vectors in special-purpose memory.
	The CPU automatically adds FF0000H to the 16-bit vector in special-purpose memory to calculate the address of the interrupt service routine, and then executes the routine.

	The peripheral transaction server (PTS), a microcoded hardware interrupt processor, provides high-speed, low-overhead interrupt handling;
	it does not modify the stack or the PSW.
	You can configure most interrupts (except NMI, software trap (TRAP), unimplemented opcode, stack overflow, PIH0_INT, and PIH1_INT) to be serviced by the PTS instead of the interrupt controller.

	The PTS provides four special microcoded routines that enable it to complete specific tasks faster than an equivalent interrupt service routine.
	It canransfer bytes or words, either individually or in blocks, between any memory locations in page 00H;
	abort PTS service if a dummy PTS request occurs; and test for a missing event in a series of regular events.
	PTS interrupts have a higher pri- ority than standard interrupts and may temporarily suspend interrupt service routines.

	A block of data called the PTS control block (PTSCB) contains the specific details for each PTS routine (see “Initializing the PTS Control Blocks” on page 6-30).
	When a PTS interrupt occurs, the priority resolver selects the appropriate vector from special-purpose memory and fetches the PTS control block (PTSCB).

	To provide support for the large number of event processor array (EPA) channels, the 8XC196EA inrporates two peripheral interrupt handlers (PIHs).
	Each PIH services 16 different interrupt sources. You can select either interrupt controller or PTS service for each H interrupt source.
	When a PIH receives an interrupt request from an enabled source, it generates either a standard interrupt request or PTS service request to the CPU.
	Although the PIH interrupt vectors are stored in special-purpose memory, the PIH must supply the interrupt vector address to the CPU.



*/

type InterruptVector struct {
	InterruptSource string
	Mnemonic        string
	Type            string
}

var interruptVectors = map[int]InterruptVector{
	0x17203E: {
		InterruptSource: "Nonmaskable Interrupt",
		Mnemonic:        "NMI",
		Type:            "Interrupt Controller Service",
	},
	0x17203C: {
		InterruptSource: "Stack Overflow Error",
		Mnemonic:        "Stack",
		Type:            "Interrupt Controller Service",
	},

	// PIH0 PTS Interrupt and PIH0 Standard Interrupt
	// ========================================================================

	0x1720FC: {
		InterruptSource: "(PIH0) EPA Capture/Compare 15",
		Mnemonic:        "(PIH0) EPA15",
		Type:            "Interrupt Controller Service",
	},
	0x1720FE: {
		InterruptSource: "(PIH0) EPA Capture/Compare 15",
		Mnemonic:        "(PIH0) EPA15",
		Type:            "PTS Service",
	},

	0x1720F8: {
		InterruptSource: "(PIH0) EPA Capture/Compare 14",
		Mnemonic:        "(PIH0) EPA14",
		Type:            "Interrupt Controller Service",
	},
	0x1720FA: {
		InterruptSource: "(PIH0) EPA Capture/Compare 14",
		Mnemonic:        "(PIH0) EPA14",
		Type:            "PTS Service",
	},

	0x1720F4: {
		InterruptSource: "(PIH0) EPA Capture/Compare 13",
		Mnemonic:        "(PIH0) EPA13",
		Type:            "Interrupt Controller Service",
	},
	0x1720F6: {
		InterruptSource: "(PIH0) EPA Capture/Compare 13",
		Mnemonic:        "(PIH0) EPA13",
		Type:            "PTS Service",
	},

	0x1720F0: {
		InterruptSource: "(PIH0) EPA Capture/Compare 12",
		Mnemonic:        "(PIH0) EPA12",
		Type:            "Interrupt Controller Service",
	},
	0x1720F2: {
		InterruptSource: "(PIH0) EPA Capture/Compare 12",
		Mnemonic:        "(PIH0) EPA12",
		Type:            "PTS Service",
	},

	0x1720EC: {
		InterruptSource: "(PIH0) EPA Capture/Compare 11",
		Mnemonic:        "(PIH0) EPA11",
		Type:            "Interrupt Controller Service",
	},
	0x1720EE: {
		InterruptSource: "(PIH0) EPA Capture/Compare 11",
		Mnemonic:        "(PIH0) EPA11",
		Type:            "PTS Service",
	},

	0x1720E8: {
		InterruptSource: "(PIH0) EPA Capture/Compare 10",
		Mnemonic:        "(PIH0) EPA10",
		Type:            "Interrupt Controller Service",
	},
	0x1720EA: {
		InterruptSource: "(PIH0) EPA Capture/Compare 10",
		Mnemonic:        "(PIH0) EPA10",
		Type:            "PTS Service",
	},

	0x1720E4: {
		InterruptSource: "(PIH0) EPA Capture/Compare 9",
		Mnemonic:        "(PIH0) EPA9",
		Type:            "Interrupt Controller Service",
	},
	0x1720E6: {
		InterruptSource: "(PIH0) EPA Capture/Compare 9",
		Mnemonic:        "(PIH0) EPA9",
		Type:            "PTS Service",
	},

	0x1720E0: {
		InterruptSource: "(PIH0) EPA Capture/Compare 8",
		Mnemonic:        "(PIH0) EPA8",
		Type:            "Interrupt Controller Service",
	},
	0x1720E2: {
		InterruptSource: "(PIH0) EPA Capture/Compare 8",
		Mnemonic:        "(PIH0) EPA8",
		Type:            "PTS Service",
	},

	0x1720DC: {
		InterruptSource: "(PIH0) EPA Capture/Compare 7",
		Mnemonic:        "(PIH0) EPA7",
		Type:            "Interrupt Controller Service",
	},
	0x1720DE: {
		InterruptSource: "(PIH0) EPA Capture/Compare 7",
		Mnemonic:        "(PIH0) EPA7",
		Type:            "PTS Service",
	},

	0x1720D8: {
		InterruptSource: "(PIH0) EPA Capture/Compare 6",
		Mnemonic:        "(PIH0) EPA6",
		Type:            "Interrupt Controller Service",
	},
	0x1720DA: {
		InterruptSource: "(PIH0) EPA Capture/Compare 6",
		Mnemonic:        "(PIH0) EPA6",
		Type:            "PTS Service",
	},

	0x1720D4: {
		InterruptSource: "(PIH0) EPA Capture/Compare 5",
		Mnemonic:        "(PIH0) EPA5",
		Type:            "Interrupt Controller Service",
	},
	0x1720D6: {
		InterruptSource: "(PIH0) EPA Capture/Compare 5",
		Mnemonic:        "(PIH0) EPA5",
		Type:            "PTS Service",
	},

	0x1720D0: {
		InterruptSource: "(PIH0) EPA Capture/Compare 4",
		Mnemonic:        "(PIH0) EPA4",
		Type:            "Interrupt Controller Service",
	},
	0x1720D2: {
		InterruptSource: "(PIH0) EPA Capture/Compare 4",
		Mnemonic:        "(PIH0) EPA4",
		Type:            "PTS Service",
	},

	0x1720CC: {
		InterruptSource: "(PIH0) EPA Capture/Compare 3",
		Mnemonic:        "(PIH0) EPA3",
		Type:            "Interrupt Controller Service",
	},
	0x1720CE: {
		InterruptSource: "(PIH0) EPA Capture/Compare 3",
		Mnemonic:        "(PIH0) EPA3",
		Type:            "PTS Service",
	},

	0x1720C8: {
		InterruptSource: "(PIH0) EPA Capture/Compare 2",
		Mnemonic:        "(PIH0) EPA2",
		Type:            "Interrupt Controller Service",
	},
	0x1720CA: {
		InterruptSource: "(PIH0) EPA Capture/Compare 2",
		Mnemonic:        "(PIH0) EPA2",
		Type:            "PTS Service",
	},

	0x1720C4: {
		InterruptSource: "(PIH0) EPA Capture/Compare 1",
		Mnemonic:        "(PIH0) EPA1",
		Type:            "Interrupt Controller Service",
	},
	0x1720C6: {
		InterruptSource: "(PIH0) EPA Capture/Compare 1",
		Mnemonic:        "(PIH0) EPA1",
		Type:            "PTS Service",
	},

	0x1720C0: {
		InterruptSource: "(PIH0) EPA Capture/Compare 0",
		Mnemonic:        "(PIH0) EPA0",
		Type:            "Interrupt Controller Service",
	},
	0x1720C2: {
		InterruptSource: "(PIH0) EPA Capture/Compare 0",
		Mnemonic:        "(PIH0) EPA0",
		Type:            "PTS Service",
	},

	// PIH0 PTS Interrupt and PIH0 Standard Interrupt
	// ========================================================================

	0x17213C: {
		InterruptSource: "(PIH1) EPA Capture/Compare 16",
		Mnemonic:        "(PIH1) EPA16",
		Type:            "Interrupt Controller Service",
	},
	0x17213E: {
		InterruptSource: "(PIH1) EPA Capture/Compare 16",
		Mnemonic:        "(PIH1) EPA16",
		Type:            "PTS Service",
	},

	0x172138: {
		InterruptSource: "(PIH1) Output Simulcapture 7",
		Mnemonic:        "(PIH1) OS7",
		Type:            "Interrupt Controller Service",
	},
	0x17213A: {
		InterruptSource: "(PIH1) Output Simulcapture 7",
		Mnemonic:        "(PIH1) OS7",
		Type:            "PTS Service",
	},

	0x172134: {
		InterruptSource: "(PIH1) Output Simulcapture 6",
		Mnemonic:        "(PIH1) OS6",
		Type:            "Interrupt Controller Service",
	},
	0x172136: {
		InterruptSource: "(PIH1) Output Simulcapture 6",
		Mnemonic:        "(PIH1) OS6",
		Type:            "PTS Service",
	},

	0x172130: {
		InterruptSource: "(PIH1) Output Simulcapture 5",
		Mnemonic:        "(PIH1) OS5",
		Type:            "Interrupt Controller Service",
	},
	0x172132: {
		InterruptSource: "(PIH1) Output Simulcapture 5",
		Mnemonic:        "(PIH1) OS5",
		Type:            "PTS Service",
	},

	0x17212C: {
		InterruptSource: "(PIH1) Output Simulcapture 4",
		Mnemonic:        "(PIH1) OS4",
		Type:            "Interrupt Controller Service",
	},
	0x17212E: {
		InterruptSource: "(PIH1) Output Simulcapture 4",
		Mnemonic:        "(PIH1) OS4",
		Type:            "PTS Service",
	},

	0x172128: {
		InterruptSource: "(PIH1) Output Simulcapture 3",
		Mnemonic:        "(PIH1) OS3",
		Type:            "Interrupt Controller Service",
	},
	0x17212A: {
		InterruptSource: "(PIH1) Output Simulcapture 3",
		Mnemonic:        "(PIH1) OS3",
		Type:            "PTS Service",
	},

	0x172124: {
		InterruptSource: "(PIH1) Output Simulcapture 2",
		Mnemonic:        "(PIH1) OS2",
		Type:            "Interrupt Controller Service",
	},
	0x172126: {
		InterruptSource: "(PIH1) Output Simulcapture 2",
		Mnemonic:        "(PIH1) OS2",
		Type:            "PTS Service",
	},

	0x172120: {
		InterruptSource: "(PIH1) Output Simulcapture 1",
		Mnemonic:        "(PIH1) OS1",
		Type:            "Interrupt Controller Service",
	},
	0x172122: {
		InterruptSource: "(PIH1) Output Simulcapture 1",
		Mnemonic:        "(PIH1) OS1",
		Type:            "PTS Service",
	},

	0x17211C: {
		InterruptSource: "(PIH1) Output Simulcapture 0",
		Mnemonic:        "(PIH1) OS0",
		Type:            "Interrupt Controller Service",
	},
	0x17211E: {
		InterruptSource: "(PIH1) Output Simulcapture 0",
		Mnemonic:        "(PIH1) OS0",
		Type:            "PTS Service",
	},

	0x172118: {
		InterruptSource: "(PIH1) Timer 1 Overflow/Underflow",
		Mnemonic:        "(PIH1) OVRTM1",
		Type:            "Interrupt Controller Service",
	},
	0x17211A: {
		InterruptSource: "(PIH1) Timer 1 Overflow/Underflow",
		Mnemonic:        "(PIH1) OVRTM1",
		Type:            "PTS Service",
	},

	0x172114: {
		InterruptSource: "(PIH1) Timer 2 Overflow/Underflow",
		Mnemonic:        "(PIH1) OVRTM2",
		Type:            "Interrupt Controller Service",
	},
	0x172116: {
		InterruptSource: "(PIH1) Timer 2 Overflow/Underflow",
		Mnemonic:        "(PIH1) OVRTM2",
		Type:            "PTS Service",
	},

	0x172110: {
		InterruptSource: "(PIH1) Timer 3 Overflow/Underflow",
		Mnemonic:        "(PIH1) OVRTM3",
		Type:            "Interrupt Controller Service",
	},
	0x172112: {
		InterruptSource: "(PIH1) Timer 3 Overflow/Underflow",
		Mnemonic:        "(PIH1) OVRTM3",
		Type:            "PTS Service",
	},

	0x17210C: {
		InterruptSource: "(PIH1) Timer 4 Overflow/Underflow",
		Mnemonic:        "(PIH1) OVRTM4",
		Type:            "Interrupt Controller Service",
	},
	0x17210E: {
		InterruptSource: "(PIH1) Timer 4 Overflow/Underflow",
		Mnemonic:        "(PIH1) OVRTM4",
		Type:            "PTS Service",
	},

	0x172108: {
		InterruptSource: "(PIH1) EPA0 Capture Overrun",
		Mnemonic:        "(PIH1) OVR0",
		Type:            "Interrupt Controller Service",
	},
	0x17210A: {
		InterruptSource: "(PIH1) EPA0 Capture Overrun",
		Mnemonic:        "(PIH1) OVR0",
		Type:            "PTS Service",
	},

	0x172104: {
		InterruptSource: "(PIH1) EPA1 Capture Overrun",
		Mnemonic:        "(PIH1) OVR1",
		Type:            "Interrupt Controller Service",
	},
	0x172106: {
		InterruptSource: "(PIH1) EPA1 Capture Overrun",
		Mnemonic:        "(PIH1) OVR1",
		Type:            "PTS Service",
	},

	0x172100: {
		InterruptSource: "(PIH1) EPA2 Capture Overrun",
		Mnemonic:        "(PIH1) OVR2",
		Type:            "Interrupt Controller Service",
	},
	0x172102: {
		InterruptSource: "(PIH1) EPA2 Capture Overrun",
		Mnemonic:        "(PIH1) OVR2",
		Type:            "PTS Service",
	},

	0x172032: {
		InterruptSource: "SSIO Channel 1 Transfer",
		Mnemonic:        "SSIO1",
		Type:            "Interrupt Controller Service",
	},
	0x172052: {
		InterruptSource: "SSIO Channel 1 Transfer",
		Mnemonic:        "SSIO1",
		Type:            "PTS Service",
	},

	0x172030: {
		InterruptSource: "SSIO Channel 0 Transfer",
		Mnemonic:        "SSIO0",
		Type:            "Interrupt Controller Service",
	},
	0x172050: {
		InterruptSource: "SSIO Channel 0 Transfer",
		Mnemonic:        "SSIO0",
		Type:            "PTS Service",
	},

	0x172016: {
		InterruptSource: "Dummy PTS Cycle",
		Mnemonic:        "--",
		Type:            "PTS Service",
	},

	0x172014: {
		InterruptSource: "Dummy Standard Interrupt",
		Mnemonic:        "--",
		Type:            "Interrupt Controller Service",
	},

	0x172012: {
		InterruptSource: "Unimplemented Opcode",
		Mnemonic:        "--",
		Type:            "Interrupt Controller Service",
	},

	0x172010: {
		InterruptSource: "Software TRAP Instruction",
		Mnemonic:        "--",
		Type:            "Interrupt Controller Service",
	},

	0x17200E: {
		InterruptSource: "Serial Debug Unit Interrupt",
		Mnemonic:        "SDU",
		Type:            "Interrupt Controller Service",
	},
	0x17204E: {
		InterruptSource: "Serial Debug Unit Interrupt",
		Mnemonic:        "SDU",
		Type:            "PTS Service",
	},

	0x17200C: {
		InterruptSource: "EXTINT Pin",
		Mnemonic:        "EXTINT",
		Type:            "Interrupt Controller Service",
	},
	0x17204C: {
		InterruptSource: "EXTINT Pin",
		Mnemonic:        "EXTINT",
		Type:            "PTS Service",
	},

	0x17200A: {
		InterruptSource: "SIO1 Receive",
		Mnemonic:        "RI1",
		Type:            "Interrupt Controller Service",
	},
	0x17204A: {
		InterruptSource: "SIO1 Receive",
		Mnemonic:        "RI1",
		Type:            "PTS Service",
	},

	0x172008: {
		InterruptSource: "SIO1 Transmit",
		Mnemonic:        "TI1",
		Type:            "Interrupt Controller Service",
	},
	0x172048: {
		InterruptSource: "SIO1 Transmit",
		Mnemonic:        "TI1",
		Type:            "PTS Service",
	},

	0x172006: {
		InterruptSource: "A/D Conversion Complete",
		Mnemonic:        "AD_DONE",
		Type:            "Interrupt Controller Service",
	},
	0x172046: {
		InterruptSource: "A/D Conversion Complete",
		Mnemonic:        "AD_DONE",
		Type:            "PTS Service",
	},

	0x172004: {
		InterruptSource: "EPA Channel 3–16 Overrun",
		Mnemonic:        "EPAx_OVR",
		Type:            "Interrupt Controller Service",
	},
	0x172044: {
		InterruptSource: "EPA Channel 3–16 Overrun",
		Mnemonic:        "EPAx_OVR",
		Type:            "PTS Service",
	},

	0x172002: {
		InterruptSource: "SIO0 Receive",
		Mnemonic:        "RI0",
		Type:            "Interrupt Controller Service",
	},
	0x172042: {
		InterruptSource: "SIO0 Receive",
		Mnemonic:        "RI0",
		Type:            "PTS Service",
	},

	0x172000: {
		InterruptSource: "SIO0 Transmit",
		Mnemonic:        "TI0",
		Type:            "Interrupt Controller Service",
	},
	0x172040: {
		InterruptSource: "SIO0 Transmit",
		Mnemonic:        "TI0",
		Type:            "PTS Service",
	},
}

/*
var interrupts = Interrupts{
	Interrupt{
		Name:  "Lower interrupt vectors",
		Start: 0x172000,
		Stop:  0x17200F,
	},
	Interrupt{
		Name:  "Software TRAP Instruction",
		Start: 0x172010,
		Stop:  0x172012,
	},
	Interrupt{
		Name:  "Unimplemented Opcode",
		Start: 0x172012,
		Stop:  0x172014,
	},
	Interrupt{
		Name:  "PIH dummy interrupt vector (point to a RET instruction)",
		Start: 0x172014,
		Stop:  0x172015,
	},
	Interrupt{
		Name:  "PIH dummy PTS vector (point to ZERO_REG; 0000H)",
		Start: 0x172016,
		Stop:  0x172017,
	},
	Interrupt{
		Name:  "Upper interrupt vectors",
		Start: 0x172030,
		Stop:  0x17203F,
	},
	Interrupt{
		Name:  "PTS vectors",
		Start: 0x172040,
		Stop:  0x17205D,
	},
	Interrupt{
		Name:  "PIH 0 vectors",
		Start: 0x1720C0,
		Stop:  0x1720FF,
	},
	Interrupt{
		Name:  "PIH 1 vectors",
		Start: 0x172100,
		Stop:  0x17213F,
	},
}
*/

func (h *DisAsm) GetInterrupts() error {

	h.vectorAdr = make(map[int]string)       // address of interrupt vector locations and name
	h.intRoutineNames = make(map[int]string) // address of interrupt routine locations and name

	for vec, intr := range interruptVectors {

		rAdr := (int(h.block[vec+1])<<8 | int(h.block[vec]) + 0x170000)
		h.intRoutineAdrs = append(h.intRoutineAdrs, rAdr) // slice of interrupt routine addresses for start locations
		h.vectorAdr[vec] = intr.InterruptSource
		h.intRoutineNames[rAdr] = intr.InterruptSource

		/*
			address := addSpaces(fmt.Sprintf("[0x%X]	Interrupt [%s]", rAdr, intr.InterruptSource), 80)
			shortDesc := fmt.Sprintf("Value: 0x%X ", rAdr)

			log(address+shortDesc, nil)
		*/
	}
	return nil
}
