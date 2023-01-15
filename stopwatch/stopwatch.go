package stopwatch

import "time"

type Stopwatch struct {
	Start   time.Time
	Elapsed func() time.Duration
}

func Start() *Stopwatch {
	sw := &Stopwatch{Start: time.Now()}
	sw.Elapsed = func() time.Duration {
		return time.Since(sw.Start)
	}
	return sw
}
