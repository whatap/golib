package queue

import (
	"context"
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
func (this *RequestQueue) GetTimeout2(timeout int) interface{} {
	this.lock.L.Lock()

	// ✅ 1단계: 먼저 데이터가 있는지 확인
	if this.queue.Size() > 0 {
		x := this.queue.RemoveFirst()
		this.lock.L.Unlock()
		return x
	}

	// ✅ 2단계: 타임아웃 처리
	timeoutCh := time.After(time.Duration(timeout) * time.Millisecond)
	waitDone := make(chan struct{}, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // ✅ 타임아웃 시 고루틴 취소 보장

	// ✅ 3단계: 별도 고루틴에서 Wait() 호출
	go func() {
		// ✅ 락을 다시 획득 (메인 고루틴이 Unlock() 후)
		this.lock.L.Lock()
		defer this.lock.L.Unlock()

		// ✅ 조건이 만족될 때까지 대기 (Spurious Wakeup 방지)
		for this.queue.Size() <= 0 {
			// ✅ context 취소 확인 (고루틴 누수 방지)
			select {
			case <-ctx.Done():
				return // ✅ 타임아웃 시 고루틴 종료
			default:
			}

			this.lock.Wait() // 락 해제하고 대기, Broadcast()로 깨어남

			// ✅ 다시 확인 (Broadcast() 후에도 취소되었을 수 있음)
			select {
			case <-ctx.Done():
				return // ✅ 타임아웃 시 고루틴 종료
			default:
			}
		}

		// ✅ 데이터가 도착함
		select {
		case waitDone <- struct{}{}:
		case <-ctx.Done():
			// ✅ 타임아웃 시 고루틴 종료
		}
	}()

	// ✅ 4단계: 락을 해제하고 대기
	this.lock.L.Unlock()

	select {
	case <-waitDone:
		// ✅ 데이터가 도착함, 락을 다시 획득
		this.lock.L.Lock()
		if this.queue.Size() > 0 {
			x := this.queue.RemoveFirst()
			this.lock.L.Unlock()
			return x
		}
		this.lock.L.Unlock()
		return nil
	case <-timeoutCh:
		// ✅ 타임아웃 발생, 하지만 마지막으로 한 번 더 확인 (race condition 방지)
		this.lock.L.Lock()
		if this.queue.Size() > 0 {
			// ✅ 타임아웃 직전에 데이터가 추가되었을 수 있음
			x := this.queue.RemoveFirst()
			this.lock.L.Unlock()
			return x
		}
		this.lock.Broadcast() // Wait() 중인 고루틴 깨우기
		this.lock.L.Unlock()
		return nil
	}
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
