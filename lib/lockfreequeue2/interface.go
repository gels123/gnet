package lockfreequeue2

type queueInterface interface {
	Put(interface{})
	Get() interface{}
}
