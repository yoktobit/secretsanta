package logic

import (
	"math/rand"
	"time"
)

// A Randomizer generated random numbers
type Randomizer interface {
	NextInt(upper int) int
}

type randomizer struct {
	rand *rand.Rand
}

type mockRandomizer struct {
	generatingSelf bool
}

// NewRandomizer creates a real randomizer
func NewRandomizer() Randomizer {
	rndSource := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(rndSource)
	return &randomizer{rand: rnd}
}

// NewMockRandomizer creates a new mocked randomizer
func NewMockRandomizer() Randomizer {

	return &mockRandomizer{generatingSelf: true}
}

func (rnd *randomizer) NextInt(upper int) int {
	return rnd.rand.Intn(upper)
}

func (rnd *mockRandomizer) NextInt(upper int) int {
	if upper == 4 && rnd.generatingSelf {
		rnd.generatingSelf = false
		return 0
	}
	if upper == 0 {
		return 3
	}
	return upper - 1
}
