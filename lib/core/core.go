package core

import (
	"fmt"
	"gnet/lib/logzap"
	"reflect"
)

// HelperFunctionToUseReflectCall helps to convert realparam like([]interface{}) to reflect.Call's param([]reflect.Value
// and if param is nil, then use reflect.New to create an empty to avoid crash when reflect.Call invokes.
// and genrates more readable error messages if param is not ok.
func HelperFunctionToUseReflectCall(f reflect.Value, callParam []reflect.Value, startNum int, realParam []interface{}) {
	n := len(realParam)
	lastIdx := f.Type().NumIn() - 1
	isVariadic := f.Type().IsVariadic()
	for i := 0; i < n; i++ {
		idx := i + startNum
		if !isVariadic && idx >= f.Type().NumIn() {
			errStr := fmt.Sprintf("InvocationCausedPanic(%v): called param count(%v) is len than reciver function's parma count(%v)", f.Type().String(), len(callParam), f.Type().NumIn())
			logzap.Panicw(errStr)
		}
		var expectedType reflect.Type
		if isVariadic && idx >= lastIdx { //variadic function's last param is []T
			expectedType = f.Type().In(lastIdx)
			expectedType = expectedType.Elem()
		} else {
			expectedType = f.Type().In(idx)
		}
		//if param is nil, create a empty reflect.Value
		if realParam[i] == nil {
			callParam[idx] = reflect.New(expectedType).Elem()
		} else {
			callParam[idx] = reflect.ValueOf(realParam[i])
		}
		actualType := callParam[idx].Type()
		if !actualType.AssignableTo(expectedType) {
			//panic if param is not assignable to Call
			errStr := fmt.Sprintf("InvocationCausedPanic(%v): called with a mismatched parameter type [parameter #%v: expected %v; got %v].", f.Type().String(), idx, expectedType, actualType)
			logzap.Panicw(errStr)
		}
	}
}

func PrintArgListForFunc(f reflect.Value) {
	t := f.Type()
	if t.Kind() != reflect.Func {
		fmt.Println("Not a func")
		return
	}
	inCount := t.NumIn()
	var str string
	for i := 0; i < inCount; i++ {
		et := t.In(i)
		str = str + ":" + et.Name()
	}
	fmt.Println(str)
}

// Parse Node Id parse node sid from service sid
func ParseNodeId(id SID) uint64 {
	return id.NodeId()
}

// Send send a message to dst service no src service.
func Send(dst SID, msgType MsgType, encType EncType, cmd CmdType, data ...interface{}) error {
	return _send(INVALID_SRC_ID, dst, msgType, encType, 0, cmd, data...)
}

// SendCloseToAll simple send a close msg to all service
func SendCloseToAll() {
	mgr.dicMutex.Lock()
	defer mgr.dicMutex.Unlock()
	for _, ser := range mgr.idDict {
		sendLocal(INVALID_SRC_ID, ser, MSG_TYPE_CLOSE, MSG_ENC_TYPE_NO, 0, Cmd_None, false)
	}
}

func Exit() {
	if isStandalone {
		SendCloseToAll()
	} else {
		route(Cmd_Exit)
	}
}

func ExitNodeByName(nodeName string) {
	if isStandalone {
		SendCloseToAll()
	} else {
		route(Cmd_Exit_Node, nodeName)
	}
}

func RefreshSlaveWhiteIPList(ips []string) {
	if isStandalone {
	} else {
		route(Cmd_RefreshSlaveWhiteIPList, ips)
	}
}

// Wait wait on a sync.WaitGroup, until all service is closed.
func Wait() {
	mgr.group.Wait()
}

// CheckIsLocalServiceId heck a given service sid is a local service
func CheckIsLocalServiceId(id SID) bool {
	return isLocalSid(id)
}
