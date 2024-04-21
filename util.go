package main

import "time"

func MeasureDuration(fn func()) time.Duration {
	startTime := time.Now()

	fn()

	return time.Since(startTime)
}
