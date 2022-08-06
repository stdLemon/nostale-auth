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

type Timing struct {
	Dp TimingRange
	Df TimingRange
	Dw TimingRange
	Dc TimingRange
	D  TimingRange
}

type Identity struct {
	Timing         Timing
	Fingerprint    blackbox.Fingerprint
	InstallationId string
}
