package retry

import (
	"math/rand"
	"strings"
	"time"
)

func Retry(attempts int, baseDelay time.Duration, maxDelay time.Duration, fn func() error) error {
	var err error
	delay := baseDelay

	for i := 0; i < attempts; i++ {
		if err = fn(); err == nil {
			return nil
		}

		jitter := time.Duration(rand.Int63n(int64(delay / 2)))
		sleep := delay + jitter

		time.Sleep(sleep)

		delay *= 2
		if delay > maxDelay {
			delay = maxDelay
		}
	}
	return err
}

func RetryIf(attempts int, baseDelay, maxDelay time.Duration, shouldRetry func(error) bool, fn func() error) error {
	var err error
	delay := baseDelay

	for range attempts {
		if err = fn(); err == nil {
			return nil
		}

		if !shouldRetry(err) {
			return err
		}

		jitter := time.Duration(rand.Int63n(int64(delay / 2)))
		sleep := delay + jitter

		time.Sleep(sleep)

		delay *= 2
		if delay > maxDelay {
			delay = maxDelay
		}
	}

	return err
}

func IsTransientError(err error) bool {
	if err == nil {
		return false
	}
	errMsg := strings.ToLower(err.Error())

	return strings.Contains(errMsg, "serialization_failure") ||
		strings.Contains(errMsg, "deadlock") ||
		strings.Contains(errMsg, "timeout") ||
		strings.Contains(errMsg, "connection reset") ||
		strings.Contains(errMsg, "temporary")
}
