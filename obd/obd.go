package obd

import (
    "errors"
    "io"
    "strings"
    "fmt"
    "bufio"
    "encoding/hex"
    serial "github.com/huin/goserial"
)

const testerAddr = 0xF5
const ecuAddr = 0x10

// DEBUG MODE
const debug = false

// BaudDefault is the default baud rate to connect to a ELM327 OBD-II device
const BaudDefault = 9600

// BaudFast is a faster baud rate available for some ELM327 OBD-II devices
const BaudFast = 38400

// EOL
const EOL = 0x3E

// ErrConnClosed is returned when a command is attempted on a closed OBD-II connection
var ErrConnClosed = errors.New("obd: connection to device is closed")

// ErrIdentify is returned when a device cannot be identified as a OBD-II device
var ErrIdentify = errors.New("obd: device is not a valid OBD-II device")

// ErrSetup is returned when a device cannot be configured to work with gobdii
var ErrSetup = errors.New("obd: cannot configure OBD-II device for use")

// Connection represents a ELM327 OBD-II serial connection
type Connection struct {
    serial io.ReadWriteCloser
}

// New creates a new ELM327 OBD-II serial connection
func New(device string, baud int) (Connection, error) {
    // Create serial connection configuration
    fmt.Printf("Creating new device... [%v %v]\n", device, baud)
    config := &serial.Config{
        Name: device,
        Baud: baud,
    }

    var conn io.ReadWriteCloser

    if debug == true {
        fmt.Print("Debug is on!\n")
        return Connection{nil}, nil
    }
    // Attempt to open serial connection
    fmt.Print("Opening connection...\n")
    conn, err := serial.OpenPort(config)
    if err != nil {
        fmt.Printf("Connection error!: %v\n", err)
        return Connection{conn}, err
    }

    // Create OBD-II connection
    obd := Connection{conn}

    // AT D - Sets All Defaults
    // AT E0 - Disable device echo
    // AT L0 - Disable line feed
    // AT S0 - Disable spaces
    // AT AT2 - Enable faster responses
    // AT SP 00 - Automatically select protocol
    // AT H1 - Turns on headers
    // AT L1 - Enables line feeds
    // AT CA F1 - CAN Automatic Formatting on 
    // AT AL - Allow Long Messages
    commands := []string{"AT D", "AT E0", "AT S0", "AT SP 3", "AT H1", "AT L0", "AT AL"}
    for _, c := range commands {
        // Send command, verify command received
        buf, err := obd.command(c)
        if err != nil {
            fmt.Print(err)
            return obd, ErrSetup
        }

        // Check for OK
        if !strings.Contains(string(buf), "OK") {
            return obd, errors.New("Error: \"" + string(buf) + "\"")
        }
    }

    if debug {
        return obd, nil
    } else {
        // Return new OBD-II connection
        return obd, nil
    }
}

// Close destroys connection to a ELM327 OBD-II device
func (c Connection) Close() error {
    // Reset device
    if err := c.Reset(); err != nil {
        return err
    }

    // Close connection
    return c.serial.Close()
}

// Reset closes the connection
func (c Connection) Reset() error {
    // AT Z - resets the device
    _, err := c.command("AT Z")
    return err
}

// Identify returns the identity of the current OBD-II device
func (c Connection) Identify() (string, error) {
    // AT I - Identify device
    buf, err := c.command("AT I")
    return string(buf), err
}

// Voltage returns the current battery voltage as reported by OBD-II device
func (c Connection) Voltage() (string, error) {
    // AT RV - Return battery voltage
    volt, err := c.command("AT RV")
    if err != nil {
        return "0.0V", err
    }

    return string(volt), nil
}

// Protocol
func (c Connection) Protocol() (string, error) {
    //
    proto, err := c.command("AT DP")
    if err != nil {
        return "0.0V", err
    }

    return string(proto), nil
}

// Speed returns the current vehicle speed as reported by OBD-II device
func (c Connection) Speed() (string, error) {
    // 01 0D - Return current vehicle speed
    speed, err := c.command("01 0D")
    if err != nil {
        return "0.0", err
    }

    // Convert from hex to decimal
    //return strconv.ParseInt("0x"+string(speed[4:6]), 0, 32)
    return string(speed), nil
}

func (c Connection) DumpBIN() error {

    // Set our header once for now
    resp, err := c.command("AT SH 74 10 F1")
    if err != nil {
        return err
    } else {
        fmt.Printf("SH command response: %v\n", resp)
    }

    for i := 0x80; i < 0xFF; i = i + 0x04  {
        dumpCommand := []byte{0x23, 0x10, 0x80, byte(i)}
        resp, err := c.command(toString(dumpCommand))
        if err == nil {
            fmt.Printf("Response: %v\n\n", string(resp))
        } else {
            return err
        }
    }

    return nil

}

func (c Connection) Write(cmd []byte) ([]byte, error) {

    // Prepend our message with the proper headers for hex commands
    h1 := byte((len(cmd) + 3) <<4 ) + 0x04  // length +1 for the checksum
    fullCmd := append([]byte{h1, ecuAddr, testerAddr} , cmd...)
    chks := iso_checksum(fullCmd)
    fmt.Printf("HEX: %.2X %.2X %.2X %2.2X %.2X\n",h1, ecuAddr, testerAddr, cmd, chks)

    if debug {
        return []byte{0x01, 0x2F, 0x2F, 0x10}, nil
    }
    // Check for open connection
    if c.serial == nil {
        return nil, ErrConnClosed
    }

    // Set the headers
    shCmd := append([]byte("AT SH"), h1, ecuAddr, testerAddr)
    shCmd = append(shCmd, []byte("\r")...)
    if _, err := c.serial.Write(shCmd); err != nil {
        return nil, err
    }

    // Read OBD-II response, loop until a response is generated
    reader := bufio.NewReader(c.serial)
    reply, err := reader.ReadBytes(EOL)
    if err != nil {
        panic(err)
    }

    // Issue command to device
    cmd = append(cmd, []byte("\r")...)
    if _, err := c.serial.Write(cmd); err != nil {
        return nil, err
    }

    // Read OBD-II response, loop until a response is generated
    reader = bufio.NewReader(c.serial)
    reply, err = reader.ReadBytes(EOL)
    if err != nil {
        panic(err)
    }

    fmt.Printf("RAW: %X \n\n", reply)
    out := strings.Replace(string(reply[:len(reply)-3]), "\x0d", "\n", -1)
    fmtOut := []byte(strings.Trim(out, ">"))

    return fmtOut, nil

}

func iso_checksum(data []byte) byte {
    crc := byte(0x00);
    for i := 0; i < len(data); i++ {
        crc = crc + data[i]
    }
    return crc;
}

// Encode our hex to a string so we can output it to the ELM 327
func toString(in []byte) string {
    return hex.EncodeToString(in)
}

// Decode our string from the ELM 327 into a byte array we can work with
func toByte(in string) ([]byte) {
    byt, _ := hex.DecodeString(in)
    return byt
}

func (c Connection) setHeader(cmd []byte) (error) {
    // Header automation
    h1 := byte((len(cmd) + 3) <<4 ) + 0x04  // length +1 for the checksum
    fullCmd := append([]byte{h1, ecuAddr, testerAddr} , cmd...)
    chks := iso_checksum(fullCmd) // just for the pretty factor
    fmt.Printf("Sending ECU: %.2X %.2X %.2X %2.2X %.2X\n",h1, ecuAddr, testerAddr, cmd, chks)

    _, err := c.command("AT SH 74 10 F1")
    if err != nil {
        return err
    }
    return nil
}

// command issues a command and retrieves a response from an OBD-II device
func (c Connection) command(cmd string) ([]byte, error) {

    // See if this is for the ELM or ECU
    if !strings.Contains(cmd, "AT") {
        err := c.setHeader(toByte(cmd))
        if err != nil {
            return []byte{}, err
        }
    }

    if debug {
        return []byte("DEBUG\n"), nil
    } else {
        // Check for open connection
        if c.serial == nil {
            fmt.Print("No open connections!")
            return nil, ErrConnClosed
        }

        // Issue command to device
        if _, err := c.serial.Write([]byte(cmd + "\r")); err != nil {
            return nil, err
        }

        // Read OBD-II response, loop until a response is generated
        reader := bufio.NewReader(c.serial)
        reply, err := reader.ReadBytes(EOL)
        reply = []byte(strings.Trim(string(reply[:]), "\r\n>"))
        if err != nil {
            return []byte{}, err
        }

        return reply, nil
    }
}
