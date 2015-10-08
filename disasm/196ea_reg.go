package disasm

import (
//"errors"
)

func regName(s string, v int) string {
	if v == 0x00 {
		return s
	}

	if name, okk := RegObjs[v]; okk {
		// Return from the list below
		return s + " ~(" + name.Description + ")"

	} else if v >= 0x400 && v <= 0xFFF {
		// Internal RAM
		return s + " ~( Internal RAM - Code RAM )"

	} else if v >= 0x00 && v <= 0x17 {
		// Special Function Registers
		return s + " ~( SFR - Special Function Registers )"

	} else if v >= 0x18 && v <= 0x19 {
		// Stack Pointer
		return s + " ~( SP - Stack Pointer )"

	} else if v >= 0x1A && v <= 0xFF {
		// General Purpose Register RAM - lower register
		return s + " ~( GP Reg RAM )"

	} else if v >= 0x100 && v <= 0x3FFF {
		// Special Function Registers - upper register
		return s + " ~( GP Reg RAM )"

	}

	return s + " ~"
}

type Register struct {
	Mnemonic        string
	Description     string
	LongDescription string
}

var RegObjs = map[int]Register{
	0x1E72: {
		Mnemonic:    "AD_RESULT",
		Description: "A/D Result",
	},

	0x1E74: {
		Mnemonic:    "AD_COMMAND",
		Description: "A/D Command",
	},

	0x1E50: {
		Mnemonic:    "AD_RESULT0",
		Description: " A/D Result 0",
	},

	0x1E52: {
		Mnemonic:    "AD_RESULT1",
		Description: " A/D Result 1",
	},

	0x1E54: {
		Mnemonic:    "AD_RESULT2 ",
		Description: "A/D Result 2",
	},

	0x1E56: {
		Mnemonic:    "AD_RESULT3",
		Description: "A/D Result 3",
	},

	0x1E58: {
		Mnemonic:    "AD_RESULT4",
		Description: "A/D Result 4",
	},

	0x1E5A: {
		Mnemonic:    "AD_RESULT5",
		Description: "A/D Result 5",
	},

	0x1E5C: {
		Mnemonic:    "AD_RESULT6",
		Description: "A/D Result 6",
	},

	0x1E5E: {
		Mnemonic:    "AD_RESULT7",
		Description: "A/D Result 7",
	},

	0x1E60: {
		Mnemonic:    "AD_RESULT8",
		Description: "A/D Result 8",
	},

	0x1E62: {
		Mnemonic:    "AD_RESULT9",
		Description: "A/D Result 9",
	},

	0x1E64: {
		Mnemonic:    "AD_RESULT10",
		Description: "A/D Result 10",
	},

	0x1E66: {
		Mnemonic:    "AD_RESULT11",
		Description: "A/D Result 11",
	},

	0x1E68: {
		Mnemonic:    "AD_RESULT12",
		Description: "A/D Result 12",
	},

	0x1E6A: {
		Mnemonic:    "AD_RESULT13",
		Description: "A/D Result 13",
	},

	0x1E6C: {
		Mnemonic:    "AD_RESULT14",
		Description: "A/D Result 14",
	},

	0x1E6E: {
		Mnemonic:    "AD_RESULT15",
		Description: "A/D Result 15",
	},

	0x1E70: {
		Mnemonic:    "AD_SCAN",
		Description: "A/D Scan",
	},

	0x1E76: {
		Mnemonic:    "AD_TEST",
		Description: "A/D Test",
	},

	0x1E77: {
		Mnemonic:    "AD_TIME",
		Description: "A/D Time",
	},

	0x1E78: {
		Mnemonic:    "ADDRCOM0",
		Description: "Address Compare 0",
	},

	0x1E80: {
		Mnemonic:    "ADDRCOM1",
		Description: "Address Compare 1",
	},

	0x1E88: {
		Mnemonic:    "ADDRCOM2",
		Description: "Address Compare 2",
	},

	0x1E7A: {
		Mnemonic:    "ADDRMSK0",
		Description: "Address Mask 0",
	},

	0x1E82: {
		Mnemonic:    "ADDRMSK1",
		Description: "Address Mask 1",
	},

	0x1E8A: {
		Mnemonic:    "ADDRMSK2",
		Description: "Address Mask 2",
	},

	0x1E7C: {
		Mnemonic:    "BUSCON0",
		Description: "Bus Control 0",
	},

	0x1E84: {
		Mnemonic:    "BUSCON1",
		Description: "Bus Control 1",
	},

	0x1E8C: {
		Mnemonic:    "BUSCON2",
		Description: "Bus Control 2",
	},

	0x1F80: {
		Mnemonic:    "CLKOUT_CON",
		Description: "Clock Out Control",
	},

	0x1FE3: {
		Mnemonic:    "EP_DIR",
		Description: "Extended Port I/O Direction",
	},

	0x1FE1: {
		Mnemonic:    "EP_MODE",
		Description: "Extended Port Mode",
	},

	0x1FE7: {
		Mnemonic:    "EP_PIN",
		Description: "Extended Port Pin Input",
	},

	0x1FE5: {
		Mnemonic:    "EP_REG",
		Description: "Extended Port Data Output",
	},

	0x1F5C: {
		Mnemonic:    "EPA0_CON",
		Description: "EPA Capture/Compare 0 Control",
	},

	0x1F58: {
		Mnemonic:    "EPA1_CON",
		Description: "EPA Capture/Compare 1 Control",
	},

	0x1F54: {
		Mnemonic:    "EPA2_CON",
		Description: "EPA Capture/Compare 2 Control",
	},

	0x1F50: {
		Mnemonic:    "EPA3_CON",
		Description: "EPA Capture/Compare 3 Control",
	},

	0x1F4C: {
		Mnemonic:    "EPA4_CON",
		Description: "EPA Capture/Compare 4 Control",
	},

	0x1F48: {
		Mnemonic:    "EPA5_CON",
		Description: "EPA Capture/Compare 5 Control",
	},

	0x1F44: {
		Mnemonic:    "EPA6_CON",
		Description: "EPA Capture/Compare 6 Control",
	},

	0x1F40: {
		Mnemonic:    "EPA7_CON",
		Description: "EPA Capture/Compare 7 Control",
	},

	0x1F3C: {
		Mnemonic:    "EPA8_CON",
		Description: "EPA Capture/Compare 8 Control",
	},

	0x1F38: {
		Mnemonic:    "EPA9_CON",
		Description: "EPA Capture/Compare 9 Control",
	},

	0x1F34: {
		Mnemonic:    "EPA10_CON",
		Description: "EPA Capture/Compare 10 Control",
	},

	0x1F30: {
		Mnemonic:    "EPA11_CON",
		Description: "EPA Capture/Compare 11 Control",
	},

	0x1F2C: {
		Mnemonic:    "EPA12_CON",
		Description: "EPA Capture/Compare 12 Control",
	},

	0x1F28: {
		Mnemonic:    "EPA13_CON",
		Description: "EPA Capture/Compare 13 Control",
	},

	0x1F24: {
		Mnemonic:    "EPA14_CON",
		Description: "EPA Capture/Compare 14 Control",
	},

	0x1F20: {
		Mnemonic:    "EPA15_CON",
		Description: "EPA Capture/Compare 15 Control",
	},

	0x1F1C: {
		Mnemonic:    "EPA16_CON",
		Description: "EPA Capture/Compare 16 Control",
	},

	0x1F5E: {
		Mnemonic:    "EPA0_TIME",
		Description: "EPA Capture/Compare 0 Time",
	},

	0x1F5A: {
		Mnemonic:    "EPA1_TIME",
		Description: "EPA Capture/Compare 1 Time",
	},

	0x1F56: {
		Mnemonic:    "EPA2_TIME",
		Description: "EPA Capture/Compare 2 Time",
	},

	0x1F52: {
		Mnemonic:    "EPA3_TIME",
		Description: "EPA Capture/Compare 3 Time",
	},

	0x1F4E: {
		Mnemonic:    "EPA4_TIME",
		Description: "EPA Capture/Compare 4 Time",
	},

	0x1F4A: {
		Mnemonic:    "EPA5_TIME",
		Description: "EPA Capture/Compare 5 Time",
	},

	0x1F46: {
		Mnemonic:    "EPA6_TIME",
		Description: "EPA Capture/Compare 6 Time",
	},

	0x1F42: {
		Mnemonic:    "EPA7_TIME",
		Description: "EPA Capture/Compare 7 Time",
	},

	0x1F3E: {
		Mnemonic:    "EPA8_TIME",
		Description: "EPA Capture/Compare 8 Time",
	},

	0x1F3A: {
		Mnemonic:    "EPA9_TIME",
		Description: "EPA Capture/Compare 9 Time",
	},

	0x1F36: {
		Mnemonic:    "EPA10_TIME",
		Description: "EPA Capture/Compare 10 Time",
	},

	0x1F32: {
		Mnemonic:    "EPA11_TIME",
		Description: "EPA Capture/Compare 11 Time",
	},

	0x1F2E: {
		Mnemonic:    "EPA12_TIME",
		Description: "EPA Capture/Compare 12 Time",
	},

	0x1F2A: {
		Mnemonic:    "EPA13_TIME",
		Description: "EPA Capture/Compare 13 Time",
	},

	0x1F26: {
		Mnemonic:    "EPA14_TIME",
		Description: "EPA Capture/Compare 14 Time",
	},

	0x1F22: {
		Mnemonic:    "EPA15_TIME",
		Description: "EPA Capture/Compare 15 Time",
	},

	0x1F1E: {
		Mnemonic:    "EPA16_TIME",
		Description: "EPA Capture/Compare 16 Time",
	},

	0x0008: {
		Mnemonic:    "INT_MASK",
		Description: "Interrupt Mask",
	},

	0x0013: {
		Mnemonic:    "INT_MASK1",
		Description: "Interrupt Mask 1",
	},

	0x0009: {
		Mnemonic:    "INT_PEND",
		Description: "Interrupt Pending",
	},

	0x0012: {
		Mnemonic:    "INT_PEND1",
		Description: "Interrupt Pending 1",
	},

	0x1FE0: {
		Mnemonic:    "IRAM_CON",
		Description: "Internal RAM Control",
	},

	0x0002: {
		Mnemonic:    "ONES_REG",
		Description: "Ones Register",
	},

	0x1EFC: {
		Mnemonic:    "OS0_CON",
		Description: "Output Simulcapture 0 Control",
	},

	0x1EF8: {
		Mnemonic:    "OS1_CON",
		Description: "Output Simulcapture 1 Control",
	},

	0x1EF4: {
		Mnemonic:    "OS2_CON",
		Description: "Output Simulcapture 2 Control",
	},

	0x1EF0: {
		Mnemonic:    "OS3_CON",
		Description: "Output Simulcapture 3 Control",
	},

	0x1EEC: {
		Mnemonic:    "OS4_CON",
		Description: "Output Simulcapture 4 Control",
	},

	0x1EE8: {
		Mnemonic:    "OS5_CON",
		Description: "Output Simulcapture 5 Control",
	},

	0x1EE4: {
		Mnemonic:    "OS6_CON",
		Description: "Output Simulcapture 6 Control",
	},

	0x1EE0: {
		Mnemonic:    "OS7_CON",
		Description: "Output Simulcapture 7 Control",
	},

	0x1EFE: {
		Mnemonic:    "OS0_TIME",
		Description: "Output Simulcapture 0 Time",
	},

	0x1EFA: {
		Mnemonic:    "OS1_TIME",
		Description: "Output Simulcapture 1 Time",
	},

	0x1EF6: {
		Mnemonic:    "OS2_TIME",
		Description: "Output Simulcapture 2 Time",
	},

	0x1EF2: {
		Mnemonic:    "OS3_TIME",
		Description: "Output Simulcapture 3 Time",
	},

	0x1EEE: {
		Mnemonic:    "OS4_TIME",
		Description: "Output Simulcapture 4 Time",
	},

	0x1EEA: {
		Mnemonic:    "OS5_TIME",
		Description: "Output Simulcapture 5 Time",
	},

	0x1EE6: {
		Mnemonic:    "OS6_TIME",
		Description: "Output Simulcapture 6 Time",
	},

	0x1EE2: {
		Mnemonic:    "OS7_TIME",
		Description: "Output Simulcapture 7 Time",
	},

	0x1FD2: {
		Mnemonic:    "P2_DIR",
		Description: "Port 2 I/O Direction",
	},

	0x1FF3: {
		Mnemonic:    "P5_DIR",
		Description: "Port 5 I/O Direction",
	},

	0x1FCA: {
		Mnemonic:    "P7_DIR",
		Description: "Port 7 I/O Direction",
	},

	0x1FCB: {
		Mnemonic:    "P8_DIR",
		Description: "Port 8 I/O Direction",
	},

	0x1FC2: {
		Mnemonic:    "P9_DIR",
		Description: "Port 9 I/O Direction",
	},

	0x1FC3: {
		Mnemonic:    "P10_DIR",
		Description: "Port 10 I/O Direction",
	},

	0x1FBA: {
		Mnemonic:    "P11_DIR",
		Description: "Port 11 I/O Direction",
	},

	0x1FEA: {
		Mnemonic:    "P12_DIR",
		Description: "Port 12 I/O Direction",
	},

	0x1FD0: {
		Mnemonic:    "P2_MODE",
		Description: "Port 2 Mode",
	},

	0x1FF1: {
		Mnemonic:    "P5_MODE",
		Description: "Port 5 Mode",
	},

	0x1FC8: {
		Mnemonic:    "P7_MODE",
		Description: "Port 7 Mode",
	},

	0x1FC9: {
		Mnemonic:    "P8_MODE",
		Description: "Port 8 Mode",
	},

	0x1FC0: {
		Mnemonic:    "P9_MODE",
		Description: "Port 9 Mode",
	},

	0x1FC1: {
		Mnemonic:    "P10_MODE",
		Description: "Port 10 Mode",
	},

	0x1FB8: {
		Mnemonic:    "P11_MODE",
		Description: "Port 11 Mode",
	},

	0x1FE8: {
		Mnemonic:    "P12_MODE",
		Description: "Port 12 Mode",
	},

	0x1FD6: {
		Mnemonic:    "P2_PIN",
		Description: "Port 2 Pin Input",
	},

	0x1FFE: {
		Mnemonic:    "P3_PIN",
		Description: "Port 3 Pin Input",
	},

	0x1FFF: {
		Mnemonic:    "P4_PIN",
		Description: "Port 4 Pin Input",
	},

	0x1FF7: {
		Mnemonic:    "P5_PIN",
		Description: "Port 5 Pin Input",
	},

	0x1FCE: {
		Mnemonic:    "P7_PIN",
		Description: "Port 7 Pin Input",
	},

	0x1FCF: {
		Mnemonic:    "P8_PIN",
		Description: "Port 8 Pin Input",
	},

	0x1FC6: {
		Mnemonic:    "P9_PIN",
		Description: "Port 9 Pin Input",
	},

	0x1FC7: {
		Mnemonic:    "P10_PIN",
		Description: "Port 10 Pin Input",
	},

	0x1FBE: {
		Mnemonic:    "P11_PIN",
		Description: "Port 11 Pin Input",
	},

	0x1FEE: {
		Mnemonic:    "P12_PIN",
		Description: "Port 12 Pin Input",
	},

	0x1FD4: {
		Mnemonic:    "P2_REG",
		Description: "Port 2 Data Output",
	},

	0x1FFC: {
		Mnemonic:    "P3_REG",
		Description: "Port 3 Data Output",
	},

	0x1FFD: {
		Mnemonic:    "P4_REG",
		Description: "Port 4 Data Output",
	},

	0x1FF5: {
		Mnemonic:    "P5_REG",
		Description: "Port 5 Data Output",
	},

	0x1FCC: {
		Mnemonic:    "P7_REG",
		Description: "Port 7 Data Output",
	},

	0x1FCD: {
		Mnemonic:    "P8_REG",
		Description: "Port 8 Data Output",
	},

	0x1FC4: {
		Mnemonic:    "P9_REG",
		Description: "Port 9 Data Output",
	},

	0x1FC5: {
		Mnemonic:    "P10_REG",
		Description: "Port 10 Data Output",
	},

	0x1FBC: {
		Mnemonic:    "P11_REG",
		Description: "Port 11 Data Output",
	},

	0x1FEC: {
		Mnemonic:    "P12_REG",
		Description: "Port 12 Data Output",
	},

	0x1FF4: {
		Mnemonic:    "P34_DRV",
		Description: "Port 3/4 Push-pull Enable",
	},

	0x1E98: {
		Mnemonic:    "PIH0_INT_MASK",
		Description: "Peripheral Int Handler 0 Int Mask",
	},

	0x1EA8: {
		Mnemonic:    "PIH1_INT_MASK",
		Description: "Peripheral Int Handler 1 Int Mask",
	},

	0x1E9A: {
		Mnemonic:    "PIH0_INT_PEND",
		Description: "Peripheral Int Handler 0 Int Pending",
	},

	0x1EAA: {
		Mnemonic:    "PIH1_INT_PEND",
		Description: "Peripheral Int Handler 1 Int Pending",
	},

	0x1E96: {
		Mnemonic:    "PIH0_PTSSEL",
		Description: "Peripheral Int Handler 0 PTS Select",
	},

	0x1EA6: {
		Mnemonic:    "PIH1_PTSSEL",
		Description: "Peripheral Int Handler 1 PTS Select",
	},

	0x1E94: {
		Mnemonic:    "PIH0_PTSSRV",
		Description: "Peripheral Int Handler 0 PTS Service",
	},

	0x1EA4: {
		Mnemonic:    "PIH1_PTSSRV",
		Description: "Peripheral Int Handler 1 PTS Service",
	},

	0x1E92: {
		Mnemonic:    "PIH0_VEC_BASE",
		Description: "Peripheral Int Handler 0 Vector Base",
	},

	0x1EA2: {
		Mnemonic:    "PIH1_VEC_BASE",
		Description: "Peripheral Int Handler 1 Vector Base",
	},

	0x1E90: {
		Mnemonic:    "PIH0_VEC_IDX",
		Description: "Peripheral Int Handler 0 Vector Index",
	},

	0x1EA0: {
		Mnemonic:    "PIH1_VEC_IDX",
		Description: "Peripheral Int Handler 1 Index",
	},

	0x0004: {
		Mnemonic:    "PTSSEL",
		Description: "PTS Select",
	},

	0x0006: {
		Mnemonic:    "PTSSRV",
		Description: "PTS Service",
	},

	0x1EDD: {
		Mnemonic:    "PWM0_1_COUNT",
		Description: "PWM 0 and 1 Count",
	},

	0x1ED9: {
		Mnemonic:    "PWM2_3_COUNT",
		Description: "PWM 2 and 3 Count",
	},

	0x1ED5: {
		Mnemonic:    "PWM4_5_COUNT",
		Description: "PWM 4 and 5 Count",
	},

	0x1ED1: {
		Mnemonic:    "PWM6_7_COUNT",
		Description: "PWM 6 and 7 Count",
	},

	0x1EDF: {
		Mnemonic:    "PWM0_1_PERIOD",
		Description: "PWM 0 and 1 Period",
	},

	0x1EDB: {
		Mnemonic:    "PWM2_3_PERIOD",
		Description: "PWM 2 and 3 Period",
	},

	0x1ED7: {
		Mnemonic:    "PWM4_5_PERIOD",
		Description: "PWM 4 and 5 Period",
	},

	0x1ED3: {
		Mnemonic:    "PWM6_7_PERIOD",
		Description: "PWM 6 and 7 Period",
	},

	0x1EDE: {
		Mnemonic:    "PWM0_CONTROL",
		Description: "PWM 0 Control",
	},

	0x1EDC: {
		Mnemonic:    "PWM1_CONTROL",
		Description: "PWM 1 Control",
	},

	0x1EDA: {
		Mnemonic:    "PWM2_CONTROL",
		Description: "PWM 2 Control",
	},

	0x1ED8: {
		Mnemonic:    "PWM3_CONTROL",
		Description: "PWM 3 Control",
	},

	0x1ED6: {
		Mnemonic:    "PWM4_CONTROL",
		Description: "PWM 4 Control",
	},

	0x1ED4: {
		Mnemonic:    "PWM5_CONTROL",
		Description: "PWM 5 Control",
	},

	0x1ED2: {
		Mnemonic:    "PWM6_CONTROL",
		Description: "PWM 6 Control",
	},

	0x1ED0: {
		Mnemonic:    "PWM7_CONTROL",
		Description: "PWM 7 Control",
	},

	0x1FA4: {
		Mnemonic:    "RSTSRC",
		Description: "Reset Source Indicator",
	},

	0x1F88: {
		Mnemonic:    "SBUF0_RX",
		Description: "Serial Port Receive Buffer 0",
	},

	0x1F98: {
		Mnemonic:    "SBUF1_RX",
		Description: "Serial Port Receive Buffer 1",
	},

	0x1F8A: {
		Mnemonic:    "SBUF0_TX",
		Description: "Serial Port Transmit Buffer 0",
	},

	0x1F9A: {
		Mnemonic:    "SBUF1_TX",
		Description: "Serial Port Transmit Buffer 1",
	},

	0x0018: {
		Mnemonic:    "SP",
		Description: "Stack Pointer",
	},

	0x0019: {
		Mnemonic:    "SP",
		Description: "Stack Pointer",
	},

	0x1F8C: {
		Mnemonic:    "SP0_BAUD",
		Description: "Serial Port 0 Baud Rate",
	},

	0x1F9C: {
		Mnemonic:    "SP1_BAUD",
		Description: "Serial Port 1 Baud Rate",
	},

	0x1F8B: {
		Mnemonic:    "SP0_CON",
		Description: "Serial Port 0 Control",
	},

	0x1F9B: {
		Mnemonic:    "SP1_CON",
		Description: "Serial Port 1 Control",
	},

	0x1F89: {
		Mnemonic:    "SP0_STATUS",
		Description: "Serial Port 0 Status",
	},

	0x1F99: {
		Mnemonic:    "SP1_STATUS",
		Description: "Serial Port 1 Status",
	},

	0x1F94: {
		Mnemonic:    "SSIO_BAUD ",
		Description: "Synchronous Serial Port Baud Rate",
	},

	0x1F90: {
		Mnemonic:    "SSIO0_BUF",
		Description: "Synchronous Serial Port 0 Buffer",
	},

	0x1F92: {
		Mnemonic:    "SSIO1_BUF",
		Description: "Synchronous Serial Port 1 Buffer",
	},

	0x1F95: {
		Mnemonic:    "SSIO0_CLK",
		Description: "Synchronous Serial Port 0 Clock",
	},

	0x1F97: {
		Mnemonic:    "SSIO1_CLK",
		Description: "Synchronous Serial Port 1 Clock",
	},

	0x1F91: {
		Mnemonic:    "SSIO0_CON",
		Description: "Synchronous Serial Port 0 Control",
	},

	0x1F93: {
		Mnemonic:    "SSIO1_CON",
		Description: "Synchronous Serial Port 1 Control",
	},

	0x1FA0: {
		Mnemonic:    "STACK_BOTTOM",
		Description: "Stack Bottom",
	},

	0x1FA2: {
		Mnemonic:    "STACK_TOP",
		Description: "Stack Top",
	},

	0x1F7C: {
		Mnemonic:    "T1CONTROL",
		Description: "Timer 1 Control",
	},

	0x1F78: {
		Mnemonic:    "T2CONTROL",
		Description: "Timer 2 Control",
	},

	0x1F74: {
		Mnemonic:    "T3CONTROL",
		Description: "Timer 3 Control",
	},

	0x1F70: {
		Mnemonic:    "T4CONTROL",
		Description: "Timer 4 Control",
	},

	0x1F7E: {
		Mnemonic:    "TIMER1",
		Description: "Timer 1 Value",
	},

	0x1F7A: {
		Mnemonic:    "TIMER2",
		Description: "Timer 2 Value",
	},

	0x1F76: {
		Mnemonic:    "TIMER3",
		Description: "Timer 3 Value",
	},

	0x1F72: {
		Mnemonic:    "TIMER4",
		Description: "Timer 4 Value",
	},

	0x1F6E: {
		Mnemonic:    "TIMER_MUX",
		Description: "Timer Multiplexer",
	},

	0x000A: {
		Mnemonic:    "WATCHDOG",
		Description: "Watchdog Timer",
	},

	0x0014: {
		Mnemonic:    "WSR",
		Description: "Window Selection",
	},

	0x0015: {
		Mnemonic:    "WSR1",
		Description: "Window Selection 1",
	},

	0x0000: {
		Mnemonic:    "ZERO_REG",
		Description: "Zero Register",
	},

	/*
		CCR0
		Chip Configuration 0
	*/

	/*
		CCR1
		Chip Configuration 1
	*/

	/*
		PSW
		Processor Status Word
		no direct access
	*/

}
