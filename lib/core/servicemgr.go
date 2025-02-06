/*
 * 服务管理器
 */
package core

import (
	"errors"
	"gnet/game/conf"
	"gnet/lib/logzap"
	"sync"
)

// 服务管理器
type serviceMgr struct {
	idDict map[uint64]*ServiceBase  	// 服务ID-地址字典
	nameDict map[string]*ServiceBase  	// 服务名称-地址字典
	nodeId   uint64   					// 集群节点ID
	curId    uint64   					// 当前服务ID
	mutex    sync.Mutex  				// 互斥锁
	group  sync.WaitGroup  				// 等待组
}

// 服务管理器实例(全局唯一)
var (
	mgr *serviceMgr
)

func newServiceMgr() *serviceMgr {
	mgr := &serviceMgr{}
	if conf.NodeID >= NODE_ID_MAX {
		logzap.Panic("newServiceMgr error: NodeID >= NODE_ID_MAX")
	}
	mgr.nodeId = uint64(conf.NodeID)
	mgr.curId = uint64(MIN_SRC_ID)
	return mgr
}

// 是否本地服务
func isLocalSid(sid SID) bool {
	nodeId := sid.NodeId()
	if nodeId == mgr.nodeId {
		return true
	}
	return false
}

// 是否本地服务
// isLocalName checks a given name is a local name. A name Start with '.' is a local name. otherwise, is a all global name
func isLocalName(name string) bool {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()
	if len(name) == 0 {
		logzap.Panicw("isLocalName error: name invalid", "name=", name)
	}
	_, ok := mgr.nameDict[name]
	if ok {
		return true
	}
	return false
}

// 注册服务
func registService(s *ServiceBase) SID {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()
	name := s.GetName()
	if len(name) == 0 {
		logzap.Panicw("registService error: name invalid", "name=", name)
	}
	s, ok := mgr.nameDict[s.GetName()]
	if ok {
		logzap.Panicw("registService error: name exist", "name=", name)
	}
	if mgr.curId >= NODE_ID_MAX {
		logzap.Panicw("registService error: curId >= NODE_ID_MAX", "name=", name)
	}
	mgr.curId++
	id := mgr.curId
	id = mgr.nodeId<<NODE_ID_OFF | id
	mgr.idDict[id] = s
	mgr.nameDict[name] = s
	sid := SID(id)
	s.setId(sid)
	mgr.group.Add(1)
	return sid
}

// 取消注册服务
func unregistService(s *ServiceBase) {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()
	sid := s.GetId()
	id := uint64(sid)
	_, ok := mgr.idDict[id]
	if !ok {
		logzap.Warnw("unregistService ignore: id invalid", "id=", id)
		return
	}
	delete(mgr.idDict, id)
	delete(mgr.nameDict, s.GetName())
	mgr.group.Done()
}

// 根据id查找本地服务
func findServiceById(sid SID) (s *ServiceBase, err error) {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()
	id := uint64(sid)
	s, ok := mgr.idDict[id]
	if !ok {
		err = errors.New("findServiceById failed")
		return nil, err
	}
	return s, err
}

// 根据名称查找本地服务
func findServiceByName(name string) (s *ServiceBase, err error) {
	mgr.mutex.Lock()
	defer mgr.mutex.Unlock()
	if len(name) == 0 {
		panic("findServiceByName error: name invalid")
	}
	s, ok := mgr.nameDict[name]
	if !ok {
		err = errors.New("findServiceByName failed")
		return nil, err
	}
	return s, err
}

// 初始化
func init() {
	mgr = newServiceMgr()
}
