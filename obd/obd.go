package obd

import (
    "errors"
    "io"
    "strconv"
    "strings"
//    "fmt"
    "bufio"

    serial "github.com/huin/goserial"
)

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

    // Attempt to open serial connection
    conn, err := serial.OpenPort(config)
    if err != nil {
        return Connection{conn}, err
    }

    // Create OBD-II connection
    obd := Connection{conn}

    // AT E0 - Disable device echo
    // AT L0 - Disable line feed
    // AT S0 - Disable spaces
    // AT AT2 - Enable faster responses
    // AT SP 00 - Automatically select protocol
    commands := []string{"AT E0", "AT L0", "AT S0", "AT AT2", "AT SP 00"}
    for _, c := range commands {
        // Send command, verify command received
        buf, err := obd.command(c)
        if err != nil {
            return obd, ErrSetup
        }

        // Check for OK
        if !strings.Contains(string(buf), "OK") {
            return obd, errors.New("obd: received bad response: \"" + string(buf) + "\"")
        }
    }

    // Return new OBD-II connection
    return obd, nil
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

// Identify returns the identity of the current OBD-II device
func (c Connection) Identify() (string, error) {
    // AT I - Identify device
    buf, err := c.command("AT I")
    return string(buf), err
}

// Reset resets the state of the OBD-II device
func (c Connection) Reset() error {
    // AT Z - Reset device to defaults
    _, err := c.command("AT Z")
    return err
}

// Voltage returns the current battery voltage as reported by OBD-II device
func (c Connection) Voltage() (string, error) {
    // AT RV - Return battery voltage
    volt, err := c.command("AT RV")
    if err != nil {
        return "0.0V", err
    }

    // Return a float64
    return string(volt), nil
}

// Speed returns the current vehicle speed as reported by OBD-II device
func (c Connection) Speed() (int64, error) {
    // 01 0D - Return current vehicle speed
    speed, err := c.command("01 0D")
    if err != nil {
        return 0.0, err
    }

    // Convert from hex to decimal
    return strconv.ParseInt("0x"+string(speed[4:6]), 0, 32)
}

// command issues a command and retrieves a response from an OBD-II device
func (c Connection) command(cmd string) ([]byte, error) {
    // Check for open connection
    if c.serial == nil {
        return nil, ErrConnClosed
    }

    // Issue command to device
    //fmt.Printf("Sending command: %v\n", cmd)
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

        out := []byte(strings.Trim(string(reply), "\r\n>"))

        //fmt.Println(string(out))

    // Return trimmed response buffer
    return out, nil
}
