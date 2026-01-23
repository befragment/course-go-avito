package retry

import (
	"context"
	"errors"
	"fmt"
	"time"

)

var ErrMaxAttemptsExceeded = errors.New("max retry attempts exceeded")

type Strategy interface {
	NextDelay(attempt int) time.Duration
}

type RetryConfig struct {
	MaxAttempts int
	Strategy    Strategy
	ShouldRetry func(error) bool
}

type RetryExecutor struct {
	config RetryConfig
	logger logger
}

func NewRetryExecutor(config RetryConfig, logger logger) *RetryExecutor {
	if config.MaxAttempts == 0 {
		config.MaxAttempts = 3
	}
	if config.ShouldRetry == nil {
		config.ShouldRetry = func(err error) bool { return err != nil }
	}
	return &RetryExecutor{config: config, logger: logger}
}

func (r *RetryExecutor) Execute(fn func() error) error {
	var lastErr error

	for attempt := 1; attempt <= r.config.MaxAttempts; attempt++ {
		err := fn()
		if err == nil {
			return nil
		}

		lastErr = err

		if !r.config.ShouldRetry(err) {
			return err
		}

		if attempt == r.config.MaxAttempts {
			break
		}

		delay := r.config.Strategy.NextDelay(attempt)
		time.Sleep(delay)
	}

	return fmt.Errorf("%w: %v", ErrMaxAttemptsExceeded, lastErr)
}

func (r *RetryExecutor) ExecuteWithContext(ctx context.Context, fn func(context.Context) error) error {
	var lastErr error

	for attempt := 1; attempt <= r.config.MaxAttempts; attempt++ {

		if err := ctx.Err(); err != nil {
			return err
		}

		err := fn(ctx)
		if err == nil {
			return nil
		}

		lastErr = err

		if !r.config.ShouldRetry(err) {
			return err
		}

		if attempt == r.config.MaxAttempts {
			r.logger.Warnf("Attempt %d failed (last), retrying is stopped", attempt)
			break
		} else {
			r.logger.Warnf("Attempt %d failed, retrying...", attempt)
		}

		delay := r.config.Strategy.NextDelay(attempt)

		select {
		case <-time.After(delay):
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return fmt.Errorf("%w: %v", ErrMaxAttemptsExceeded, lastErr)
}

func (r *RetryExecutor) ExecuteWithCallback(
	fn func() error,
	onRetry func(attempt int, err error, delay time.Duration),
) error {
	var lastErr error

	for attempt := 1; attempt <= r.config.MaxAttempts; attempt++ {
		err := fn()
		if err == nil {
			return nil
		}

		lastErr = err

		if !r.config.ShouldRetry(err) {
			return err
		}

		if attempt == r.config.MaxAttempts {
			break
		}

		delay := r.config.Strategy.NextDelay(attempt)

		if onRetry != nil {
			onRetry(attempt, err, delay)
		}

		time.Sleep(delay)
	}

	return fmt.Errorf("%w: %v", ErrMaxAttemptsExceeded, lastErr)
}
