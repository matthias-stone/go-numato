package numato

import (
	"flag"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var serialPath = flag.String("serial", "", "Path to serial port")

func init() {
	flag.Parse()
}

func TestHardware(t *testing.T) {
	if *serialPath == "" {
		t.Skip("Skipping hardware test. Run with -serial=serialpath to enable")
	}

	n, err := Open(*serialPath)
	if !assert.NoError(t, err) {
		t.Fatal()
	}

	p := []Port{
		{Relay, 0},
		{Relay, 1},
		{Relay, 2},
		{Relay, 3},
	}

	for _, port := range p {
		assert.NoError(t, n.On(port), "turning port on", port)
		time.Sleep(time.Millisecond * 50)
		isOn, err := n.IsOn(port)
		assert.NoError(t, err, "reading port status", port)
		assert.True(t, isOn, "port not on", port)
	}

	for _, port := range p {
		assert.NoError(t, n.Off(port), "turning port off", port)
		time.Sleep(time.Millisecond * 50)
		isOn, err := n.IsOn(port)
		assert.NoError(t, err, "reading port status", port)
		assert.False(t, isOn, "port not off", port)
	}

	assert.NoError(t, n.Close())
}

func TestGPIOBlink(t *testing.T) {
	if *serialPath == "" {
		t.Skip("Skipping hardware test. Run with -serial=serialpath to enable")
	}

	n, err := Open(*serialPath)
	if !assert.NoError(t, err) {
		t.Fatal()
	}

	p := []Port{
		{GPIO, 0},
		{GPIO, 1},
		{GPIO, 2},
		{GPIO, 3},
		{GPIO, 4},
		{GPIO, 5},
	}

	for i := 0; i < 10; i++ {
		for _, port := range p {
			assert.NoError(t, n.On(port), "turning port on", port)
		}
		for _, port := range p {
			assert.NoError(t, n.Off(port), "turning port off", port)
		}
	}
}

func TestHammer(t *testing.T) {
	if *serialPath == "" {
		t.Skip("Skipping hardware test. Run with -serial=serialpath to enable")
	}

	n, err := Open(*serialPath)
	if !assert.NoError(t, err) {
		t.Fatal()
	}

	port := Port{Relay, 3}

	for i := 0; i < 50; i++ {
		assert.NoError(t, n.On(port), "turning port on", port)
		assert.NoError(t, n.Off(port), "turning port off", port)
	}
}
