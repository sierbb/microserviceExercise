package sliding_window

import "sync"

type Bucket struct{
	TotalCount int
	FailedCount int
	sync.RWMutex
}


func (b *Bucket) Incr(status bool) {
	b.Lock()
	defer b.Unlock()
	// record request status
	if status == false{
		b.FailedCount += 1
	}
	b.TotalCount += 1
}

func NewBucket () *Bucket{
	return &Bucket{}
}