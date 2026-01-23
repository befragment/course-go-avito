package fulljitter_test

import (
	"courier-service/pkg/delay/fulljitter"
	"math"
	"math/rand"
	"testing"
	"time"
)

func TestFullJitter_Table(t *testing.T) {
	seed := func(v int64) *int64 { return &v }

	type tc struct {
		name       string
		base       time.Duration
		max        time.Duration
		multiplier float64
		attempt    int
		wantUpper  time.Duration
		seed       *int64 // nil => не фиксируем rand
	}

	tests := []tc{
		{
			name:       "attempt_1_upper_is_base_seeded",
			base:       100 * time.Millisecond,
			max:        10 * time.Second,
			multiplier: 2,
			attempt:    1,
			wantUpper:  100 * time.Millisecond,
			seed:       seed(1),
		},
		{
			name:       "attempt_3_upper_is_base_times_4_seeded",
			base:       100 * time.Millisecond,
			max:        10 * time.Second,
			multiplier: 2,
			attempt:    3,
			wantUpper:  400 * time.Millisecond,
			seed:       seed(2),
		},
		{
			name:       "upper_capped_by_max_seeded",
			base:       500 * time.Millisecond,
			max:        2 * time.Second,
			multiplier: 3,
			attempt:    10,
			wantUpper:  2 * time.Second,
			seed:       seed(3),
		},
		{
			name:       "unseeded_still_within_range",
			base:       200 * time.Millisecond,
			max:        5 * time.Second,
			multiplier: 2,
			attempt:    2,
			wantUpper:  400 * time.Millisecond,
			seed:       nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.seed != nil {
				rand.Seed(*tt.seed)
			}

			f := fulljitter.NewFullJitter(tt.base, tt.max, tt.multiplier)

			// 1) проверяем формулу upper bound
			upper := float64(f.BaseDelay) * math.Pow(f.Multiplier, float64(tt.attempt-1))
			if upper > float64(f.MaxDelay) {
				upper = float64(f.MaxDelay)
			}
			gotUpper := time.Duration(upper)
			if gotUpper != tt.wantUpper {
				t.Fatalf("upper bound mismatch: want %v, got %v", tt.wantUpper, gotUpper)
			}

			// 2) проверяем диапазон NextDelay
			got := f.NextDelay(tt.attempt)
			if got < 0 {
				t.Fatalf("delay must be >= 0, got %v", got)
			}
			if got >= tt.wantUpper {
				t.Fatalf("delay must be < %v, got %v", tt.wantUpper, got)
			}
		})
	}
}