package fulljitter

import (
	"math"
	"math/rand"
	"time"
)

// FullJitter implements exponential backoff with full jitter strategy.
// It calculates delay as: random(0, min(base * multiplier^(attempt-1), max))
type FullJitter struct {
	BaseDelay  time.Duration
	MaxDelay   time.Duration
	Multiplier float64
	rnd        *rand.Rand
}

func NewFullJitter(base, max time.Duration, multiplier float64, initialSeed *int64) *FullJitter {
	var seed int64
	if initialSeed == nil {
		seed = time.Now().UnixNano()
	} else {
		seed = *initialSeed
	}
	return &FullJitter{
		BaseDelay:  base,
		MaxDelay:   max,
		Multiplier: multiplier,
		rnd:        rand.New(rand.NewSource(seed)),
	}
}

func (f *FullJitter) NextDelay(attempt int) time.Duration {
	maxDelay := float64(f.BaseDelay) * math.Pow(f.Multiplier, float64(attempt-1))
	if maxDelay > float64(f.MaxDelay) {
		maxDelay = float64(f.MaxDelay)
	}

	return time.Duration(f.rnd.Float64() * maxDelay)
}
