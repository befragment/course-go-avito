package tokenbucket_test

import (
	"testing"
	"time"

	"courier-service/pkg/ratelimiter/tokenbucket"
)

func TestTokenbucket(t *testing.T) {
	type step struct {
		advance        time.Duration // на сколько "промотать" время
		calls          int           // сколько раз вызвать Allow()
		wantOK         int           // сколько должно быть true среди этих calls
		expectedTokens *int          // nil => don't check
	}

	tokens := func(v int) *int { return &v }

	tests := []struct {
		name       string
		capacity   int
		refillRate int
		steps      []step
	}{
		{
			name:       "simple_case",
			capacity:   2,
			refillRate: 1,
			steps: []step{
				{advance: 0, calls: 5, wantOK: 2, expectedTokens: tokens(0)},
			},
		},
		{
			name:       "long_complex_case_limit_not_exceeded",
			capacity:   200,
			refillRate: 1,
			steps: []step{
				{advance: 0, calls: 100, wantOK: 100, expectedTokens: tokens(100)},
				{advance: 4 * time.Second, calls: 2, wantOK: 2, expectedTokens: tokens(102)},
			},
		},
		{
			name:       "long_complex_case_limit_exceeded",
			capacity:   200,
			refillRate: 1,
			steps: []step{
				{advance: 0, calls: 100, wantOK: 100, expectedTokens: tokens(100)},
				{advance: 20 * time.Second, calls: 10, wantOK: 10, expectedTokens: tokens(110)},
				{advance: 0, calls: 120, wantOK: 110, expectedTokens: tokens(0)},
			},
		},
		{
			name:       "no refill without Allow() call",
			capacity:   10,
			refillRate: 1,
			steps: []step{
				{advance: 0, calls: 15, wantOK: 10, expectedTokens: tokens(0)},
				{advance: 10 * time.Second, calls: 0, wantOK: 0, expectedTokens: tokens(0)},
			},
		},
		{
			name:       "refill after Allow() call",
			capacity:   10,
			refillRate: 1,
			steps: []step{
				{advance: 0, calls: 15, wantOK: 10, expectedTokens: tokens(0)},
				{advance: 10 * time.Second, calls: 1, wantOK: 1, expectedTokens: tokens(9)},
			},
		},
		{
			name:       "it's been a while since last refill",
			capacity:   20,
			refillRate: 5,
			steps: []step{
				{advance: 0, calls: 20, wantOK: 20, expectedTokens: tokens(0)},
				{advance: 20 * time.Second, calls: 1, wantOK: 1},
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			fake := time.Unix(0, 0)
			nowFn := func() time.Time { return fake }

			tb := tokenbucket.NewTokenBucket(tc.capacity, tc.refillRate, nowFn)

			for si, st := range tc.steps {
				fake = fake.Add(st.advance)

				ok := 0
				for i := 0; i < st.calls; i++ {
					if tb.Allow() {
						ok++
					}
				}
				if ok != st.wantOK {
					t.Fatalf("step %d: expected %d allows, got %d", si, st.wantOK, ok)
				}

				if st.expectedTokens != nil && tb.Tokens() != *st.expectedTokens {
					t.Fatalf("step %d: expected tokens=%d, got %d", si, *st.expectedTokens, tb.Tokens())
				}
			}
		})
	}
}
