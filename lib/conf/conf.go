package conf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"reflect"
)

// 日志配置
var (
	LogFilePath   string = "./log"               //log's file path
	LogFileName   string = "game"                //log's file name
	LogFileLevel  int    = logsimple.LEVEL_MAX   //log's file level
	LogShellLevel int    = logsimple.DEBUG_LEVEL //log's shell level
	LogMaxLine    int    = 10000                 //log's max line per file
	LogBufferSize int    = 2000                  //log's max buffer size
)

// ~
func assignTo(r map[string]interface{}, target interface{}, name string) {
	if t := reflect.TypeOf(target); t.Kind() != reflect.Ptr {
		fmt.Println("Kind is", t.Kind())
		panic("Not a proper type")
	}
	if s, ok := r[name]; ok {
		t0 := reflect.TypeOf(s)
		t1 := reflect.TypeOf(target).Elem()
		// var doOK bool = true
		tSet := reflect.ValueOf(target).Elem()
		if t0.AssignableTo(t1) {
			tSet.Set(reflect.ValueOf(s))
		} else if t1.Kind() == reflect.Int {
			if t0.Kind() == reflect.Float64 {
				f := reflect.ValueOf(s).Float()
				tSet.Set(reflect.ValueOf(int(f)))
			} else {
				// doOK = false
			}
		} else {
			// doOK = false
		}
		// if doOK {
		// 	fmt.Printf("Set of <%v> to %v\n", name, s)
		// } else {
		// 	fmt.Printf("Cannot assign %v to %v\n", t0.Name(), t1.Name())
		// }
	}
}

const (
	privateConfiguraPath = ".private/svrconf.json"
)

// ~ If you wish to alter the configuration filepath
// ~ Overwrite the configures by configuration file.
func init() {
	goPath := os.ExpandEnv("$GOPATH")
	if len(goPath) <= 0 {
		return
	}
	fname := path.Join(goPath, privateConfiguraPath)
	fin, err := os.Open(fname)
	if err != nil {
		//The configure file may not exist
		return
	}
	defer fin.Close()
	chunk, r_err := ioutil.ReadAll(fin)
	if r_err != nil {
		return
	}
	var intfs interface{}
	jErr := json.Unmarshal(chunk, &intfs)
	if jErr != nil {
		return
	}

	if mIntfs, ok := intfs.(map[string]interface{}); ok {
		assignTo(mIntfs, &LogFilePath, `LogFilePath`)
		assignTo(mIntfs, &CoreIsStandalone, `IsStandalone`)
		assignTo(mIntfs, &CoreIsMaster, `IsMaster`)
		assignTo(mIntfs, &MasterListenIp, `MasterIP`)
		assignTo(mIntfs, &SlaveConnectIp, `SlaveIP`)
		assignTo(mIntfs, &MultiNodePort, `MasterPort`)
		assignTo(mIntfs, &CallTimeOut, `CallTimeout`)
	}
}

func SetMasterMode() {
	CoreIsMaster = true
	CoreIsStandalone = false
}

func SetStandaloneMode() {
	CoreIsMaster = false
	CoreIsStandalone = true
}

func SetSlaveMode() {
	CoreIsMaster = false
	CoreIsStandalone = false
}
