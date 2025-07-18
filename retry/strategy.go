package retry

import (
	"math"
	"math/rand"
	"time"
)

// FixedDelayPolicy implements a fixed delay between retries
type FixedDelayPolicy struct {
	Delay    time.Duration
	Attempts int
}

func (p *FixedDelayPolicy) NextDelay(_ int) time.Duration {
	return p.Delay
}

func (p *FixedDelayPolicy) MaxAttempts() int {
	return p.Attempts
}

// ExponentialBackoffPolicy implements exponential backoff
type ExponentialBackoffPolicy struct {
	InitialDelay time.Duration
	MaxDelay     time.Duration
	Multiplier   float64
	Attempts     int
}

func (p *ExponentialBackoffPolicy) NextDelay(attempt int) time.Duration {
	if p.Multiplier <= 0 {
		p.Multiplier = 2
	}

	delay := float64(p.InitialDelay) * math.Pow(p.Multiplier, float64(attempt-1))

	if p.MaxDelay > 0 && time.Duration(delay) > p.MaxDelay {
		return p.MaxDelay
	}

	return time.Duration(delay)
}

func (p *ExponentialBackoffPolicy) MaxAttempts() int {
	return p.Attempts
}

// ExponentialBackoffWithJitterPolicy implements exponential backoff with jitter
type ExponentialBackoffWithJitterPolicy struct {
	InitialDelay time.Duration
	MaxDelay     time.Duration
	Multiplier   float64
	Attempts     int
	JitterFactor float64 // 0 to 1, where 0 is no jitter and 1 is full jitter
}

func (p *ExponentialBackoffWithJitterPolicy) NextDelay(attempt int) time.Duration {
	if p.Multiplier <= 0 {
		p.Multiplier = 2
	}

	if p.JitterFactor < 0 {
		p.JitterFactor = 0
	} else if p.JitterFactor > 1 {
		p.JitterFactor = 1
	}

	baseDelay := float64(p.InitialDelay) * math.Pow(p.Multiplier, float64(attempt-1))

	if p.MaxDelay > 0 && time.Duration(baseDelay) > p.MaxDelay {
		baseDelay = float64(p.MaxDelay)
	}

	// Apply jitter
	if p.JitterFactor > 0 {
		jitter := baseDelay * p.JitterFactor * rand.Float64()
		baseDelay = baseDelay - (baseDelay * p.JitterFactor / 2) + jitter
	}

	return time.Duration(baseDelay)
}

func (p *ExponentialBackoffWithJitterPolicy) MaxAttempts() int {
	return p.Attempts
}

// Common policy constructors

// NewFixedDelay creates a fixed delay policy
func NewFixedDelay(delay time.Duration, maxAttempts int) Policy {
	return &FixedDelayPolicy{
		Delay:    delay,
		Attempts: maxAttempts,
	}
}

// NewExponentialBackoff creates an exponential backoff policy
func NewExponentialBackoff(initialDelay, maxDelay time.Duration, maxAttempts int) Policy {
	return &ExponentialBackoffPolicy{
		InitialDelay: initialDelay,
		MaxDelay:     maxDelay,
		Multiplier:   2,
		Attempts:     maxAttempts,
	}
}

// NewExponentialBackoffWithJitter creates an exponential backoff policy with jitter
func NewExponentialBackoffWithJitter(initialDelay, maxDelay time.Duration, maxAttempts int, jitterFactor float64) Policy {
	return &ExponentialBackoffWithJitterPolicy{
		InitialDelay: initialDelay,
		MaxDelay:     maxDelay,
		Multiplier:   2,
		Attempts:     maxAttempts,
		JitterFactor: jitterFactor,
	}
}
