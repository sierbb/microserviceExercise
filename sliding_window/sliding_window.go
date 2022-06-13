package sliding_window

import (
	"fmt"
	"sync"
	"time"
)

type SlidingWindowCounter struct {
	// bucket is a list of buckets that records the request status within a period of time in the past
	buckets []*Bucket
	fused bool
	// The timestamp of last time there was a fuse
	lastFusedTimestamp time.Time
	// size is number of buckets to be present in the sliding window
	size int
	// totalThreshold is the number of total requests that will trigger a fuse
	totalThreshold int
	// failedThreshold is the number of failed requests that will trigger a fuse
	failedThreshold float64
	// freeze time after each fuse
	freezeTime time.Duration
	sync.RWMutex
}


func NewSlidingWindow(size int, totalThreshold int, failedThreshold float64, freezeTime time.Duration) *SlidingWindowCounter{
	return &SlidingWindowCounter{
		buckets: make([]*Bucket, 0, size),
		size: size,
		totalThreshold: totalThreshold,
		failedThreshold: failedThreshold,
		freezeTime: freezeTime,
	}
}


func (s *SlidingWindowCounter) Start(){
	go func(){
		for {
			s.AddBucket()
			time.Sleep(time.Second)
		}
	}()
}

// AddBucket adds a new bucket to the end of the sliding window,
// and remove the oldest bucket from the head of sliding window if needed to maintain the size
func (s *SlidingWindowCounter) AddBucket (){
	s.Lock()
	defer s.Unlock()
	s.buckets = append(s.buckets, NewBucket())
	if len(s.buckets) > s.size{
		s.buckets = s.buckets[1:]
	}
}

func (s *SlidingWindowCounter) Incr(status bool){
	s.GetLastBucket().Incr(status)
}

// CheckCounter stats an infinite process to check in interval whether counter hits the fuse thresholds
// If we hit a fuse, we will record the fuse time
// If the fuse has lasted for longer than freezeTime, we turn off the fuse
func (s *SlidingWindowCounter) CheckCounter(){
	go func() {
		s.checkFused()
		time.Sleep(time.Millisecond * 100)
	}()
}


func (s *SlidingWindowCounter) checkFused() {
	// Only do read lock since we are not writing
	s.RLock()
	total := 0
	failed := 0
	for _, bucket := range s.buckets{
		total += bucket.TotalCount
		failed += bucket.FailedCount
	}
	s.RUnlock()

	if !s.fused{
		if total >= s.totalThreshold && float64(failed/total) >= s.failedThreshold {
			s.Lock()
			s.fused = true
			s.lastFusedTimestamp = time.Now()
			s.Unlock()
		}
	} else {
		if time.Since(s.lastFusedTimestamp) >= s.freezeTime{
			s.Lock()
			s.fused = false
			s.Unlock()
		}
	}
}

func (s *SlidingWindowCounter) GetLastBucket () *Bucket{
	if len(s.buckets) > 0{
		return s.buckets[len(s.buckets)-1]
	}
	return nil
}

func (s *SlidingWindowCounter) PrintStatus() {
	go func(){
		fmt.Printf("Fused: %v", s.fused)
		time.Sleep(time.Second)
	}()
}
