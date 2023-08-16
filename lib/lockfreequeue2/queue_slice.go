package lockfreequeue2

import (
	"sync"
)

type SliceQueue struct {
	s  []interface{}
	mu sync.Mutex
}

func NewSliceQueue() (q *SliceQueue) {
	return &SliceQueue{s: make([]interface{}, 0)}
}

func (q *SliceQueue) Put(v interface{}) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.s = append(q.s, v)
}

func (q *SliceQueue) Get() interface{} {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.s) == 0 {
		return nil
	}
	v := q.s[0]
	q.s = q.s[1:]
	return v
}

//func exampleSliceQueue() {
//	q := NewSliceQueue()
//
//	q.Put("1st item")
//	q.Put("2nd item")
//	q.Put("3rd item")
//
//	fmt.Println(q.Get())
//	fmt.Println(q.Get())
//	fmt.Println(q.Get())
//	// Output:
//	// 1st item
//	// 2nd item
//	// 3rd item
//}
