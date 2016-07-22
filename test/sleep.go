// run

package main

import (
	"time"
)

func main() {
	const delay time.Duration = 123456789
	go func() {
		time.Sleep(delay / 2)
	}()
	start := time.Now()
	time.Sleep(delay)
	delayadj := delay
	duration := time.Now().Sub(start)
	if duration < delayadj {
		print("time.Sleep(", delay, ") slept for only ", duration, "ns\n")
		panic("FAIL")
	}
}
