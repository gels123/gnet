package core

import (
	"errors"
	"gnet/game/conf"
	"gnet/lib/logzap"
	"go.uber.org/zap"
	"sync"
	"sync/atomic"
)

// 服务管理
type serviceMgr struct {
	dictId   sync.Map // ID-服务字典
	dictName sync.Map // 名称-服务字典
	nodeId   uint64   // 集群节点ID
	curId    uint64   // 当前服务ID
}

var (
	mgr *serviceMgr
	wg  sync.WaitGroup
)

func newServiceMgr() *serviceMgr {
	mgr := &serviceMgr{}
	if conf.NodeID >= NODE_ID_MAX {
		panic("new service mgr error")
	}
	mgr.nodeId = uint64(conf.NodeID)
	mgr.curId = uint64(SERVICE_ID_MIN)
	return mgr
}

// 是否本地服务
func isLocalSid(sid Sid) bool {
	nodeId := sid.NodeId()
	if nodeId == mgr.nodeId {
		return true
	}
	return false
}

// isLocalName checks a given name is a local name.
// a name Start with '.' or empty is a local name. others a all global name
func isLocalName(name string) bool {
	v, ok := mgr.dictName.Load(name)
	if v != nil && ok {
		return true
	}
	return false
}

// 注册服务
func registService(s *ServiceBase) Sid {
	v, ok := mgr.dictName.Load(s.GetName())
	if ok && v != nil {
		panic("regist service error: name exsit")
	}
	id := atomic.AddUint64(&mgr.curId, 1)
	if id >= NODE_ID_MAX {
		panic("regist service error: id invalid")
	}
	id = mgr.nodeId<<NODE_ID_OFF | id
	mgr.dictId.Store(id, s)
	mgr.dictName.Store(s.GetName(), s)
	sid := Sid(id)
	s.setId(sid)
	wg.Add(1)
	return sid
}

// 取消注册服务
func unregistService(s *ServiceBase) {
	sid := s.GetId()
	id := uint64(sid)
	_, ok := mgr.dictId.Load(id)
	if !ok {
		logzap.Warnw("unregist service error", zap.Uint64("sid", id))
		return
	}
	mgr.dictId.Delete(id)
	mgr.dictName.Delete(s.GetName())
	wg.Done()
}

// 根据id查找服务
func findServiceById(sid Sid) (s *ServiceBase, err error) {
	id := uint64(sid)
	v, ok := mgr.dictId.Load(id)
	if !ok {
		err = errors.New("find service by sid failed")
		return nil, err
	}
	s = v.(*ServiceBase)
	return s, err
}

// findServiceByName return a ServiceBase by ServiceBase name, it only return local ServiceBase.
func findServiceByName(name string) (s *ServiceBase, err error) {
	if len(name) == 0 {
		panic("find service by name error: name invalid")
	}
	v, ok := mgr.dictName.Load(name)
	if !ok {
		err = errors.New("find service by sid failed")
		return nil, err
	}
	s = v.(*ServiceBase)
	return s, err
}

// 初始化
func init() {
	mgr = newServiceMgr()
}
