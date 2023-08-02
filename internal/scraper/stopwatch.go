package scraper

import "time"

type Stopwatch struct {
	startTime time.Time
}

func NewStopwatch() *Stopwatch {
	return &Stopwatch{time.Date(2000, 1, 1, 1, 1, 1, 651387237, time.UTC)}
}

func (stopwatch *Stopwatch) Reset() {
	stopwatch.startTime = time.Now()
}

func (stopwatch *Stopwatch) Elapsed() time.Duration {
	return time.Since(stopwatch.startTime)
}

func (stopwatch *Stopwatch) Finish(time time.Duration) bool {
	return stopwatch.Elapsed() >= time
}
