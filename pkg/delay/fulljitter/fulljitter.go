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
}

func NewFullJitter(base, max time.Duration, multiplier float64) *FullJitter {
	return &FullJitter{
		BaseDelay:  base,
		MaxDelay:   max,
		Multiplier: multiplier,
	}
}

func (f *FullJitter) NextDelay(attempt int) time.Duration {
	maxDelay := float64(f.BaseDelay) * math.Pow(f.Multiplier, float64(attempt-1))
	if maxDelay > float64(f.MaxDelay) {
		maxDelay = float64(f.MaxDelay)
	}

	return time.Duration(rand.Float64() * maxDelay)
}
