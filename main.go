package main

import (
	"serviceAvailability/sliding_window"
	"time"
)

func main(){
	freezeTime := time.Second
	size := 100
	totalThreshold := 1000
	failedThreshold := 0.8

	counter := sliding_window.NewSlidingWindow(size, totalThreshold, failedThreshold, freezeTime)
	counter.Start()
	counter.CheckCounter()
	counter.PrintStatus()
}