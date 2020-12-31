package main

import (
	"log"
	"sync"
	"time"
)

//1. 参考 Hystrix 实现一个滑动窗口计数器。
type Number struct {
	Buckets map[int64]*numberBucket
	Mutex   *sync.RWMutex
}

type numberBucket struct {
	Value int64
}

func NewNumber() *Number {
	r := &Number{
		Buckets: make(map[int64]*numberBucket),
		Mutex:   &sync.RWMutex{},
	}
	return r
}

func (r *Number) getCurrentBucket() *numberBucket {
	now := time.Now().Unix()
	var bucket *numberBucket
	var ok bool

	if bucket, ok = r.Buckets[now]; !ok {
		bucket = &numberBucket{}
		r.Buckets[now] = bucket
	}

	return bucket
}

func (r *Number) removeOldBuckets() {
	now := time.Now().Unix() - 5

	for timestamp := range r.Buckets {
		if timestamp <= now {
			delete(r.Buckets, timestamp)
		}
	}
}

func (r *Number) Increment(i int64) {
	if i == 0 {
		return
	}

	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	b := r.getCurrentBucket()
	b.Value += i
	r.removeOldBuckets()
}

func (r *Number) Sum(now time.Time) int64 {
	sum := int64(0)

	r.Mutex.RLock()
	defer r.Mutex.RUnlock()

	for timestamp, bucket := range r.Buckets {
		if timestamp >= now.Unix()-5 {
			sum += bucket.Value
		}
	}

	return sum
}

func (r *Number) Avg(now time.Time) int64 {
	return r.Sum(now) / 5
}

func main() {
	n := NewNumber()
	for _, x := range []int64{1, 2, 3, 4, 5, 6} {
		n.Increment(x)
		time.Sleep(1 * time.Second)
	}
	log.Println(n.Avg(time.Now())) //4
}
