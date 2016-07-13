// Package numato provides a simple interface for controlling a Numato USB digital input-output device.
package numato

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/tarm/serial"
)

// Numato is the core type.
type Numato struct {
	port io.ReadWriteCloser
}

type state string

// State constants represent the status of a DIO port. Only On and Off are exported because API users should not care about internal state.
const (
	unknownState = state("")
	On           = state("on")
	Off          = state("off")
	read         = state("read")
)

// PortType is either a relay, input, or output.
type portType string

// PortTypes describes the types of ports.
const (
	unknownType = portType("")
	Relay       = portType("relay")
	GPIO        = portType("gpio")
	ADC         = portType("adc")
)

var actions = map[portType]map[state]string{
	Relay: {On: "on", Off: "off", read: "read"},
	GPIO:  {On: "set", Off: "clear", read: "read"},
}

// Port is either a relay or GPIO.
type Port struct {
	Class  portType
	Number int
}

// Open attempts to detect and run a numato device at the provided location.
func Open(serialName string) (*Numato, error) {
	c := &serial.Config{Name: serialName, Baud: 9600, ReadTimeout: time.Millisecond * 10}
	s, err := serial.OpenPort(c)
	if err != nil {
		return nil, err
	}
	return &Numato{s}, nil
}

// On activates the provided port.
func (n *Numato) On(p Port) error {
	return n.action(p, On)
}

// Off turns off the provided port.
func (n *Numato) Off(p Port) error {
	return n.action(p, Off)
}

// Set forces a port to a given state.
func (n *Numato) Set(p Port, s state) error {
	return n.action(p, s)
}

// IsOn returns the current state of a relay or GPIO pin.
func (n *Numato) IsOn(p Port) (bool, error) {
	buf := make([]byte, 64)
	_, err := n.port.Read(buf) // Clear data left in buffer
	n.action(p, read)

	c, err := n.port.Read(buf)
	if err != nil {
		return false, err
	}
	// This double read fixes a propogation delay in the numato response
	if _, err = n.port.Read(buf[c:]); err != nil && err != io.EOF {
		return false, err
	}

	return bytes.Contains(buf, []byte("\n\r"+On+"\n\r")), nil
}

// Close releases the serial port.
func (n *Numato) Close() error {
	if err := n.port.Close(); err != nil {
		return err
	}
	return nil
}

func (n *Numato) action(p Port, s state) error {
	_, err := n.port.Write([]byte(fmt.Sprintf("%s %s %d\r", p.Class, actions[p.Class][s], p.Number)))
	if err != nil {
		return err
	}

	return nil
}
