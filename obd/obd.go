package obd

import (
    "errors"
    "io"
    "strings"
    "fmt"
    "bufio"
    "time"

    serial "github.com/huin/goserial"
)

const testerAddr = 0xF5
const ecuAddr = 0x10

// DEBUG MODE
const debug = true

// BaudDefault is the default baud rate to connect to a ELM327 OBD-II device
const BaudDefault = 9600

// BaudFast is a faster baud rate available for some ELM327 OBD-II devices
const BaudFast = 38400

// EOL
const EOL = '\x3e'

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
    config := &serial.Config{
        Name: device,
        Baud: baud,
    }

    var conn io.ReadWriteCloser
    if debug {
        return Connection{nil}, nil
    } else {
        // Attempt to open serial connection
        conn, err := serial.OpenPort(config)
        if err != nil {
            return Connection{conn}, err
        }
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
    // AT SI - Slow Init
    commands := []string{"AT D", "AT E1", "AT S1", "AT SP 00", "AT H1", "AT L1", "AT AL", "AT SI"}
    for _, c := range commands {
        // Send command, verify command received
        buf, err := obd.command(c)
        if err != nil {
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

func (c Connection) Write(cmd []byte) ([]byte, error) {
    start := time.Now()

    // Prepend our message with the proper headers for hex commands
    h1 := byte((len(cmd) + 3) <<4 ) + 0x04  // length +1 for the checksum
    fullCmd := append([]byte{h1, ecuAddr, testerAddr} , cmd...)
    chks := iso_checksum(fullCmd)
    fmt.Printf("HEX: %X %X %X %X %X\n",h1, ecuAddr, testerAddr, cmd, chks)

    if debug {
        return []byte{0x01, 0x2F, 0x2F, 0x10}, nil
    } else {
        // Check for open connection
        if c.serial == nil {
            return nil, ErrConnClosed
        }

        // Issue command to device
        if _, err := c.serial.Write(cmd); err != nil {
            //fmt.Printf("Error sending command: %v\n", err)
            return nil, err
        }

        // Read OBD-II response, loop until a response is generated
        reader := bufio.NewReader(c.serial)
        reply, err := reader.ReadBytes(EOL)
        if err != nil {
            panic(err)
        }

        out := strings.Replace(string(reply[:len(reply)-3]), "\x0d", "\n", -1)
        fmtOut := []byte(strings.Trim(out, ">"))

        // Return trimmed response buffer
        elapsed := time.Since(start)
        fmt.Printf("Command Time %s\n", elapsed)
        return fmtOut, nil



    }
}

func iso_checksum(data []byte) byte {
    crc := byte(0x00);
    for i := 0; i < len(data); i++ {
        crc = crc + data[i]
    }
    return crc;
}

// command issues a command and retrieves a response from an OBD-II device
func (c Connection) command(cmd string) ([]byte, error) {
    start := time.Now()

    fmt.Printf("Sending command: %v %v\n", cmd)

    if debug {
        return []byte("DEBUG\n"), nil
    } else {
        // Check for open connection
        if c.serial == nil {
            return nil, ErrConnClosed
        }

        // Issue command to device
        if _, err := c.serial.Write([]byte(cmd + "\r")); err != nil {
            //fmt.Printf("Error sending command: %v\n", err)
            return nil, err
        }

        // Read OBD-II response, loop until a response is generated
        reader := bufio.NewReader(c.serial)
        reply, err := reader.ReadBytes(EOL)
        if err != nil {
            panic(err)
        }

        out := strings.Replace(string(reply[:len(reply)-3]), "\x0d", "\n", -1)
        fmtOut := []byte(strings.Trim(out, ">"))

        // Return trimmed response buffer
        elapsed := time.Since(start)
        fmt.Printf("Command Time %s\n", elapsed)
        return fmtOut, nil
    }
}
