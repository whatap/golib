package queue

import (
	"sync"
	"time"

	"github.com/whatap/golib/util/dateutil"
	"github.com/whatap/golib/util/list"
)

type RequestQueue struct {
	queue      list.LinkedList
	capacity   int
	lock       *sync.Cond
	Failed     func(interface{})
	Overflowed func(interface{})
}

func NewRequestQueue(size int) *RequestQueue {
	q := new(RequestQueue)
	q.lock = sync.NewCond(new(sync.Mutex))
	q.capacity = size
	return q
}
func (this *RequestQueue) Get() interface{} {
	this.lock.L.Lock()
	defer this.lock.L.Unlock()

	for this.queue.Size() <= 0 {
		this.lock.Wait()
	}
	x := this.queue.RemoveFirst()
	return x
}
func (this *RequestQueue) GetTimeout(timeout int) interface{} {
	t := timeout
	timeto := dateutil.SystemNow() + int64(timeout)

	// 최대 3~ 4 회 루프
	var v interface{} = nil
	for v = this.GetNoWait(); v == nil; v = this.GetNoWait() {
		time.Sleep(time.Duration(t/3) * time.Millisecond)

		t = int(timeto - dateutil.SystemNow())
		if t <= 0 {
			break
		}
	}

	return v
}

func (this *RequestQueue) GetNoWait() interface{} {
	this.lock.L.Lock()
	defer this.lock.L.Unlock()

	if this.queue.Size() > 0 {
		return this.queue.RemoveFirst()
	} else {
		return nil
	}
}

func (this *RequestQueue) Put(v interface{}) bool {
	this.lock.L.Lock()
	defer this.lock.L.Unlock()
	if this.capacity <= 0 || this.queue.Size() < this.capacity {
		this.queue.Add(v)
		this.lock.Broadcast()
		return true
	} else {
		if this.Failed != nil {
			this.Failed(v)
		}
		//this.lock.Signal()
		this.lock.Broadcast()
		return false
	}
}
func (this *RequestQueue) PutForce(v interface{}) bool {
	this.lock.L.Lock()
	defer this.lock.L.Unlock()

	if this.capacity <= 0 || this.queue.Size() < this.capacity {
		this.queue.Add(v)
		this.lock.Broadcast()
		return true
	} else {
		for this.queue.Size() >= this.capacity {
			o := this.queue.RemoveFirst()
			if this.Overflowed != nil {
				this.Overflowed(o)
			}
		}
		this.queue.Add(v)
		this.lock.Broadcast()
		return false
	}
}

func (this *RequestQueue) Clear() {
	this.lock.L.Lock()
	defer this.lock.L.Unlock()
	this.queue.Clear()
}
func (this *RequestQueue) Size() int {
	return this.queue.Size()
}

func (this *RequestQueue) GetCapacity() int {
	return this.capacity
}

func (this *RequestQueue) SetCapacity(size int) {
	this.capacity = size
}
