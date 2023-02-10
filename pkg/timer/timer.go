package timer

import "time"

type Timer struct {
	Start   time.Time
	Elapsed func() time.Duration
}

func Start() *Timer {
	t := &Timer{Start: time.Now()}
	t.Elapsed = func() time.Duration {
		return time.Since(t.Start)
	}
	return t
}
