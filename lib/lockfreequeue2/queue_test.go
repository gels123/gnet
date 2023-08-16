package lockfreequeue2

import (
	"fmt"
	"testing"
	"testing/quick"
)

func ExampleQueue() {
	q := NewQueue()

	q.Put("1st item")
	q.Put("2nd item")
	q.Put("3rd item")

	fmt.Println(q.Get())
	fmt.Println(q.Get())
	fmt.Println(q.Get())
	// Output:
	// 1st item
	// 2nd item
	// 3rd item
}

func runQueueInterface(inputs []int, q queueInterface) (outputs []interface{}) {
	for _, v := range inputs {
		if v >= 0 {
			q.Put(v)
		} else {
			outputs = append(outputs, q.Get())
		}
	}
	return outputs
}

func runQueue(inputs []int) (outputs []interface{}) {
	return runQueueInterface(inputs, NewQueue())
}

func runSliceQueue(inputs []int) (outputs []interface{}) {
	return runQueueInterface(inputs, NewSliceQueue())
}

func TestMatchWithSliceQueue(t *testing.T) {
	if err := quick.CheckEqual(runQueue, runSliceQueue, nil); err != nil {
		t.Error(err)
	}
}
