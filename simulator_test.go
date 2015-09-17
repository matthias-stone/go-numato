package numato

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type OnChecker interface {
	IsOn(Port) (bool, error)
}

func assertIsOn(t *testing.T, checker OnChecker, p Port) {
	on, err := checker.IsOn(p)
	assert.NoError(t, err)
	assert.True(t, on)
}

func assertIsOff(t *testing.T, checker OnChecker, p Port) {
	on, err := checker.IsOn(p)
	assert.NoError(t, err)
	assert.False(t, on)
}

func TestSimulator(t *testing.T) {
	sim, dummy := OpenSimulator(1, 1, 0)
	p := Port{Relay, 0}

	sim.Set(p, On)
	assertIsOn(t, dummy, p)
	assertIsOn(t, dummy, p)
	assertIsOn(t, sim, p)

	assert.NoError(t, dummy.Off(p))
	assertIsOff(t, sim, p)
	assertIsOff(t, dummy, p)
	assertIsOff(t, dummy, p)
}
