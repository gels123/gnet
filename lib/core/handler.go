package core

import (
	"errors"
	"gnet/lib/utils"
	"gnet/lib/vector"
	"sync"
)

type handleDic map[uint64]*BaseService

// a storage that stores all local services
type handleStorage struct {
	dicMutex           sync.Mutex
	dic                handleDic
	nodeId             uint64
	curId              uint64
	baseServiceIdCache *vector.Vector
}

var (
	h                   *handleStorage
	ServiceNotFindError = errors.New("service is not find.")
	exitGroup           sync.WaitGroup
)

func newHandleStorage() *handleStorage {
	h := &handleStorage{}
	h.nodeId = DEFAULT_NODE_ID
	h.dic = make(map[uint64]*BaseService)
	h.curId = uint64(INIT_SERVICE_ID)
	h.baseServiceIdCache = vector.NewCap(1000)
	return h
}

// checkIsLocalId checks a given BaseService id is a local BaseService's id
// a serviceId's node id is equal to DEFAULT_NODE_ID or NodeId is a local BaseService's id
func checkIsLocalId(id sid) bool {
	nodeId := id.NodeId()
	if nodeId == DEFAULT_NODE_ID {
		return true
	}
	if nodeId == h.nodeId {
		return true
	}
	return false
}

// checkIsLocalName checks a given name is a local name.
// a name start with '.' or empty is a local name. others a all global name
func checkIsLocalName(name string) bool {
	if len(name) == 0 {
		return true
	}
	if name[0] == '.' {
		return true
	}
	return false
}

func init() {
	h = newHandleStorage()
}

// registerService register a BaseService and allocate a BaseService id to the given BaseService.
func registerService(s *BaseService) sid {
	h.dicMutex.Lock()
	defer h.dicMutex.Unlock()
	var baseServiceId uint64
	//if h.baseServiceIdCache.Empty() {
	h.curId++
	baseServiceId = h.curId
	//} else {
	//	baseServiceId = h.baseServiceIdCache.Get().(uint64)
	//}
	id := h.nodeId<<NODE_ID_OFF | baseServiceId
	h.dic[id] = s
	sid := sid(id)
	s.setId(sid)
	exitGroup.Add(1)
	return sid(sid)
}

// unregisterService delete a BaseService and put it's to cache which can be resued again when register
func unregisterService(s *BaseService) {
	h.dicMutex.Lock()
	defer h.dicMutex.Unlock()
	id := uint64(s.getId())
	if _, ok := h.dic[id]; !ok {
		return
	}
	delete(h.dic, id)
	h.baseServiceIdCache.Put((sid(id)).BaseId())
	exitGroup.Done()
}

// findServiceById return a BaseService by BaseService id
func findServiceById(id sid) (s *BaseService, err error) {
	h.dicMutex.Lock()
	defer h.dicMutex.Unlock()
	s, ok := h.dic[uint64(id)]
	if !ok {
		err = ServiceNotFindError
	}
	return s, err
}

// findServiceByName return a BaseService by BaseService name, it only return local BaseService.
func findServiceByName(name string) (s *BaseService, err error) {
	utils.PanicWhen(len(name) == 0, "name must not empty.")
	h.dicMutex.Lock()
	defer h.dicMutex.Unlock()
	for _, value := range h.dic {
		if value.getName() == name {
			s = value
			return s, nil
		}
	}
	return nil, ServiceNotFindError
}
