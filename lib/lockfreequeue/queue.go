/*
 * lock free queue(性能佳于chan和lockfreequeue2)
 */
package lockfreequeue

import (
	"fmt"
	"runtime"
	"sync/atomic"
)

type node struct {
	putNo uint32
	getNo uint32
	value interface{}
}

// lock free queue
type Queue struct {
	cap    uint32
	capMod uint32
	putPos uint32
	getPos uint32
	nodes  []node
}

func NewQueue(cap uint32) *Queue {
	q := new(Queue)
	q.cap = minQuantity(cap)
	q.capMod = q.cap - 1
	q.putPos = 0
	q.getPos = 0
	q.nodes = make([]node, q.cap)
	for i := range q.nodes {
		nd := &q.nodes[i]
		nd.getNo = uint32(i)
		nd.putNo = uint32(i)
	}
	nd := &q.nodes[0]
	nd.getNo = q.cap
	nd.putNo = q.cap
	return q
}

func (q *Queue) String() string {
	getPos := atomic.LoadUint32(&q.getPos)
	putPos := atomic.LoadUint32(&q.putPos)
	return fmt.Sprintf("Queue{cap: %v, capMod: %v, putPos: %v, getPos: %v}", q.cap, q.capMod, putPos, getPos)
}

func (q *Queue) Capaciity() uint32 {
	return q.cap
}

func (q *Queue) Quantity() uint32 {
	var putPos, getPos uint32
	var quantity uint32
	getPos = atomic.LoadUint32(&q.getPos)
	putPos = atomic.LoadUint32(&q.putPos)

	if putPos >= getPos {
		quantity = putPos - getPos
	} else {
		quantity = q.capMod + (putPos - getPos)
	}

	return quantity
}

// put queue functions
func (q *Queue) Put(val interface{}) (ok bool, quantity uint32) {
	var putPos, putPosNew, getPos, posCnt uint32
	var nd *node
	capMod := q.capMod

	getPos = atomic.LoadUint32(&q.getPos)
	putPos = atomic.LoadUint32(&q.putPos)

	if putPos >= getPos {
		posCnt = putPos - getPos
	} else {
		posCnt = capMod + (putPos - getPos)
	}

	if posCnt >= capMod-1 {
		runtime.Gosched()
		return false, posCnt
	}

	putPosNew = putPos + 1
	if !atomic.CompareAndSwapUint32(&q.putPos, putPos, putPosNew) {
		runtime.Gosched()
		return false, posCnt
	}

	nd = &q.nodes[putPosNew&capMod]

	for {
		getNo := atomic.LoadUint32(&nd.getNo)
		putNo := atomic.LoadUint32(&nd.putNo)
		if putPosNew == putNo && getNo == putNo {
			nd.value = val
			atomic.AddUint32(&nd.putNo, q.cap)
			return true, posCnt + 1
		} else {
			runtime.Gosched()
		}
	}
}

// puts queue functions
func (q *Queue) Puts(values []interface{}) (puts, quantity uint32) {
	var putPos, putPosNew, getPos, posCnt, putCnt uint32
	capMod := q.capMod

	getPos = atomic.LoadUint32(&q.getPos)
	putPos = atomic.LoadUint32(&q.putPos)

	if putPos >= getPos {
		posCnt = putPos - getPos
	} else {
		posCnt = capMod + (putPos - getPos)
	}

	if posCnt >= capMod-1 {
		runtime.Gosched()
		return 0, posCnt
	}

	if capPuts, size := q.cap-posCnt, uint32(len(values)); capPuts >= size {
		putCnt = size
	} else {
		putCnt = capPuts
	}
	putPosNew = putPos + putCnt

	if !atomic.CompareAndSwapUint32(&q.putPos, putPos, putPosNew) {
		runtime.Gosched()
		return 0, posCnt
	}

	for posNew, v := putPos+1, uint32(0); v < putCnt; posNew, v = posNew+1, v+1 {
		var nd *node = &q.nodes[posNew&capMod]
		for {
			getNo := atomic.LoadUint32(&nd.getNo)
			putNo := atomic.LoadUint32(&nd.putNo)
			if posNew == putNo && getNo == putNo {
				nd.value = values[v]
				atomic.AddUint32(&nd.putNo, q.cap)
				break
			} else {
				runtime.Gosched()
			}
		}
	}
	return putCnt, posCnt + putCnt
}

// get queue functions
func (q *Queue) Get() (val interface{}, ok bool, quantity uint32) {
	var putPos, getPos, getPosNew, posCnt uint32
	var nd *node
	capMod := q.capMod

	putPos = atomic.LoadUint32(&q.putPos)
	getPos = atomic.LoadUint32(&q.getPos)

	if putPos >= getPos {
		posCnt = putPos - getPos
	} else {
		posCnt = capMod + (putPos - getPos)
	}

	if posCnt < 1 {
		runtime.Gosched()
		return nil, false, posCnt
	}

	getPosNew = getPos + 1
	if !atomic.CompareAndSwapUint32(&q.getPos, getPos, getPosNew) {
		runtime.Gosched()
		return nil, false, posCnt
	}

	nd = &q.nodes[getPosNew&capMod]

	for {
		getNo := atomic.LoadUint32(&nd.getNo)
		putNo := atomic.LoadUint32(&nd.putNo)
		if getPosNew == getNo && getNo == putNo-q.cap {
			val = nd.value
			nd.value = nil
			atomic.AddUint32(&nd.getNo, q.cap)
			return val, true, posCnt - 1
		} else {
			runtime.Gosched()
		}
	}
}

// gets queue functions
func (q *Queue) Gets(values []interface{}) (gets, quantity uint32) {
	var putPos, getPos, getPosNew, posCnt, getCnt uint32
	capMod := q.capMod

	putPos = atomic.LoadUint32(&q.putPos)
	getPos = atomic.LoadUint32(&q.getPos)

	if putPos >= getPos {
		posCnt = putPos - getPos
	} else {
		posCnt = capMod + (putPos - getPos)
	}

	if posCnt < 1 {
		runtime.Gosched()
		return 0, posCnt
	}

	if size := uint32(len(values)); posCnt >= size {
		getCnt = size
	} else {
		getCnt = posCnt
	}
	getPosNew = getPos + getCnt

	if !atomic.CompareAndSwapUint32(&q.getPos, getPos, getPosNew) {
		runtime.Gosched()
		return 0, posCnt
	}

	for posNew, v := getPos+1, uint32(0); v < getCnt; posNew, v = posNew+1, v+1 {
		var nd *node = &q.nodes[posNew&capMod]
		for {
			getNo := atomic.LoadUint32(&nd.getNo)
			putNo := atomic.LoadUint32(&nd.putNo)
			if posNew == getNo && getNo == putNo-q.cap {
				values[v] = nd.value
				nd.value = nil
				getNo = atomic.AddUint32(&nd.getNo, q.cap)
				break
			} else {
				runtime.Gosched()
			}
		}
	}

	return getCnt, posCnt - getCnt
}

// round 到最近的2的倍数
func minQuantity(v uint32) uint32 {
	v--
	v |= v >> 1
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16
	v++
	return v
}

func Delay(z int) {
	for x := z; x > 0; x-- {
	}
}
