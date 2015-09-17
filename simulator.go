package numato

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

// Simulator controls a dummy Numato device.
// It only deals with the input/output of the numato.Numato object and does not
// handle all valid inputs to a real Numato.
type Simulator struct {
	relays, GPIOs, ADCs uint8
	state               map[portType][]bool

	buf     bytes.Buffer
	pending []byte
}

// OpenSimulator returns a Simulator and a Numato object under its control.
func OpenSimulator(relays, GPIOs, ADCs uint8) (*Simulator, *Numato) {
	sim := &Simulator{
		relays, GPIOs, ADCs,
		map[portType][]bool{
			Relay: make([]bool, relays),
			GPIO:  make([]bool, GPIOs),
		},
		bytes.Buffer{},
		[]byte{},
	}

	dummy := &Numato{sim}

	return sim, dummy
}

// Read can be used to receive responses from the Simulator.
func (sim *Simulator) Read(b []byte) (int, error) {
	return sim.buf.Read(b)
}

// Write acts as a dummy serial port and processes any completed command.
// Incomplete commands will be buffered and handled once a '\r' is written.
func (sim *Simulator) Write(b []byte) (int, error) {
	commands := bytes.Split(b, []byte("\r"))
	commands[0] = append(sim.pending, commands[0]...)
	for i := 0; i < len(commands)-1; i++ {
		sim.process(commands[i])
	}
	sim.pending = commands[len(commands)-1]
	return sim.buf.Write(b)
}

func (sim *Simulator) process(cmd []byte) {
	// Simulate the echo behaviour
	sim.buf.Write(cmd)
	sim.buf.Write([]byte("\r"))

	components := strings.Split(string(cmd), " ")
	if len(components) != 3 {
		return
	}
	num, err := strconv.Atoi(components[2])
	if err != nil {
		return
	}
	p := Port{
		portType(components[0]),
		num,
	}
	s := state(components[1])
	switch s {
	case On:
		fallthrough
	case Off:
		sim.Set(p, s)
	case read:
		on, err := sim.IsOn(p)
		if err != nil {
			break
		}
		status := "on"
		if !on {
			status = "off"
		}
		sim.buf.Write([]byte(fmt.Sprintf("\n\r%s\n\r", status)))
	default:
		// an error happened
	}
	sim.buf.Write([]byte("\n\r> "))
}

// Close is a noop.
func (sim *Simulator) Close() error {
	return nil
}

// Off turns the simulated port off.
func (sim *Simulator) Off(p Port) {
	sim.Set(p, Off)
}

// On turns the simulated port on.
func (sim *Simulator) On(p Port) {
	sim.Set(p, On)
}

// Set sets a port to the profided state.
func (sim *Simulator) Set(p Port, s state) error {
	set, ok := sim.state[p.Class]
	if !ok {
		panic("invalid type")
	}
	if p.Number >= len(set) {
		panic("port out of range")
	}

	set[p.Number] = s == On
	return nil
}

// IsOn reads the status of the port as seen by the simulator.
func (sim *Simulator) IsOn(p Port) (bool, error) {
	set, ok := sim.state[p.Class]
	if !ok {
		panic("invalid type")
	}
	if p.Number >= len(set) {
		panic("port out of range")
	}

	return set[p.Number], nil
}
