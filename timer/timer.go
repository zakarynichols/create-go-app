package timer

import "time"

type Timer struct {
	Start time.Time
	Since func(start time.Time) time.Duration
}

func New() Timer {
	return Timer{
		Start: time.Now(),
		Since: func(start time.Time) time.Duration {
			return time.Since(start)
		},
	}
}
