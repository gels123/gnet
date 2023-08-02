package main

import (
	"fmt"
	"gnet/lib/logzap"
	"gnet/lib/utils"
	"path/filepath"
	"regexp"
	"time"
)

func init() {
}

var patternConversionRegexps = []*regexp.Regexp{
	regexp.MustCompile(`%[%+A-Za-z]`),
	regexp.MustCompile(`\*+`),
}

func main() {

	fmt.Println(utils.GetExeDir())
	fmt.Println(utils.GetCurDir())

	//logOut := filepath.Join(utils.GetCurDir(), conf.LogsConf.FilePath, conf.LogsConf.FileName)
	//fmt.Println("===================df===", logOut)
	//fmt.Println("===================xxx===", conf.LogsConf.FilePath[0] == '.')

	////logzap.Infof("Failed to fetch URL: %s", "xxxx1")
	////logzap.Errorf("Failed to fetch URL: %s", "xxxx2")
	//
	//time.AfterFunc(time.Second*10, func() {
	//	fmt.Println("==============sdfadfadfa===============")
	//})

	fileName := "gamelog" + ".%Y%m%d"
	globPattern := fileName
	for _, re := range patternConversionRegexps {
		globPattern = re.ReplaceAllString(globPattern, "*")
	}

	fullFileName := "D:/work/gnet/log/" + globPattern
	matches, err := filepath.Glob(fullFileName)
	fmt.Println("sdfa=====", matches, err)

	for {
		logzap.Debugw("========sdfadf===", "a=", 100)
		logzap.Debugf("========sdfadf===", "a=", 100)
		time.Sleep(time.Second)
	}
}
