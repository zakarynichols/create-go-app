package clock

import "time"

type Elapse struct {
	Start time.Time
	Since func(start time.Time) time.Duration
}

func Timer() Elapse {
	return Elapse{
		Start: time.Now(),
		Since: func(start time.Time) time.Duration {
			return time.Since(start)
		},
	}
}
