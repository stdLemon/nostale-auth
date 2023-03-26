package identitymgr

import (
	"math/rand"

	"github.com/stdLemon/nostale-auth/pkg/blackbox"
)

type TimingRange struct {
	Min int
	Max int
}

func (t *TimingRange) Random() float64 {
	return float64(rand.Intn(t.Max-t.Min) + t.Min)
}

type Identity struct {
	Timing         TimingRange
	Fingerprint    blackbox.Fingerprint
	InstallationId string
}
