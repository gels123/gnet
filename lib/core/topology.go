package core

import (
	"gnet/lib/logsimple"
	"gnet/lib/vector"
	"sync"
)

type nameRet struct {
	id   Sid
	ok   bool
	name string
}

const (
	StartingNodeID = 10
)

var (
	once                   sync.Once
	isStandalone, isMaster bool
	registerNodeChan       chan uint64
	nameChanMap            map[uint]chan *nameRet
	nameMapMutex           sync.Mutex
	nameRequestId          uint
	beginNodeId            uint64
	validNodeIdVec         *vector.Vector
)

func init() {
	registerNodeChan = make(chan uint64, 1)

	nameChanMap = make(map[uint]chan *nameRet)
	nameRequestId = 0

	isStandalone = true

	beginNodeId = StartingNodeID
	validNodeIdVec = vector.New()
}

// InitNode set node's information.
func InitNode(_isStandalone, _isMaster bool) {
	isStandalone = _isStandalone
	isMaster = _isMaster
	if !isStandalone && isMaster {
		mgr.nodeId = MASTER_NODE_ID
	}
}

// RegisterNode : register slave node to master, and get a node sid
// block until register success
func RegisterNode(nodeName string) {
	once.Do(func() {
		if !isStandalone && !isMaster {
			route(Cmd_RegisterNode, nodeName)
			mgr.nodeId = <-registerNodeChan
			worker, _ = NewIdWorker(int64(mgr.nodeId))
			logsimple.Info("SlaveNode register ndoe success: NodeId: %v, nodeName: {%v}", mgr.nodeId, nodeName)
		}
	})
}

// DispatchRegisterNodeRet send RegisterNode's return to channel which RegisterNode is wait for
func DispatchRegisterNodeRet(id uint64) {
	registerNodeChan <- id
}

// globalName regist name to master
// it will notify all exist ServiceBase through distribute msg.
func globalName(id Sid, name string) {
	route(Cmd_RegisterName, uint64(id), name)
}

// route send msg to master
// if node is not a master node, it send to .slave node first, .slave will forward msg to master.
func route(cmd CmdType, data ...interface{}) bool {
	router, err := findServiceByName(".router")
	if err != nil {
		return false
	}
	localSendWithoutMutex(INVALID_SERVICE_ID, router, MSG_TYPE_NORMAL, MSG_ENC_TYPE_NO, 0, cmd, data...)
	return true
}

// NameToId couldn't guarantee get the correct sid for name.
// it will return err if the named server is until now.
func NameToId(name string) (Sid, error) {
	ser, err := findServiceByName(name)
	if err == nil {
		return ser.getId(), nil
	}
	if !isLocalName(name) {
		nameMapMutex.Lock()
		nameRequestId++
		tmp := nameRequestId
		nameMapMutex.Unlock()

		ch := make(chan *nameRet, 1)
		nameMapMutex.Lock()
		nameChanMap[tmp] = ch
		nameMapMutex.Unlock()

		route(Cmd_GetIdByName, name, tmp)
		ret := <-ch
		close(ch)
		if !ret.ok {
			return INVALID_SERVICE_ID, ServiceNotFindError
		}
		return ret.id, nil
	}
	return INVALID_SERVICE_ID, ServiceNotFindError
}

func DispatchGetIdByNameRet(id Sid, ok bool, name string, rid uint) {
	nameMapMutex.Lock()
	ch := nameChanMap[rid]
	delete(nameChanMap, rid)
	nameMapMutex.Unlock()
	ch <- &nameRet{id, ok, name}
}

func GenerateNodeId() uint64 {
	var ret uint64
	if validNodeIdVec.Empty() {
		ret = beginNodeId
		beginNodeId++
	} else {
		ret = validNodeIdVec.Get().(uint64)
	}
	return ret
}

func CollectNodeId(recycledNodeID uint64) {
	if recycledNodeID >= StartingNodeID {
		logsimple.Info("Recycled of NodeID<%v>", recycledNodeID)
		validNodeIdVec.Put(recycledNodeID)
	}
}
