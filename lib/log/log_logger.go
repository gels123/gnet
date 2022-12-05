package log

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

type SimpleLogger struct {
	//file log
	fileLevel  int
	fileLogger *log.Logger
	baseFile   *os.File
	outputPath string
	fileName   string
	curName    string
	curIdx     int
	maxLine    int
	logLine    int
	//shell log
	shellLevel int
	isColored  bool
	buffer     chan *Msg
	wg         sync.WaitGroup
}

const (
	COLOR_DEBUG_LEVEL_DESC = "[\x1b[32mdebug\x1b[0m] "
	COLOR_INFO_LEVEL_DESC  = "[\x1b[36minfo\x1b[0m] "
	COLOR_WARN_LEVEL_DESC  = "[\x1b[33mwarn\x1b[0m] "
	COLOR_ERROR_LEVEL_DESC = "[\x1b[31merror\x1b[0m] "
	COLOR_FATAL_LEVEL_DESC = "[\x1b[31mfatal\x1b[0m] "
)

var color_format = []string{
	COLOR_DEBUG_LEVEL_DESC,
	COLOR_INFO_LEVEL_DESC,
	COLOR_WARN_LEVEL_DESC,
	COLOR_ERROR_LEVEL_DESC,
	COLOR_FATAL_LEVEL_DESC,
}

var std_format = []string{
	DEBUG_LEVEL_DESC,
	INFO_LEVEL_DESC,
	WARN_LEVEL_DESC,
	ERROR_LEVEL_DESC,
	FATAL_LEVEL_DESC,
}

func (self *SimpleLogger) SetColored(colored bool) {
	self.isColored = colored
}

func (self *SimpleLogger) doPrintf(level int, levelDesc, msg string) {
	nformat := levelDesc + msg
	if level >= self.fileLevel {
		if self.logLine > self.maxLine {
			if self.baseFile != nil {
				self.baseFile.Sync()
				self.baseFile.Close()
			}
			self.baseFile, _ = self.createLogFile(self.outputPath, self.fileName)
			self.fileLogger = log.New(self.baseFile, "", log.Ldate|log.Lmicroseconds)
			self.logLine = 0
		}
		if self.fileLogger != nil {
			self.fileLogger.Printf(nformat)
			self.logLine++
		}
	}
	if level >= self.shellLevel {
		sel_fmt := color_format
		if !self.isColored {
			sel_fmt = std_format
		}
		nformat := sel_fmt[level] + msg
		log.Printf(nformat)
	}
}

func (self *SimpleLogger) DoPrintf(level int, levelDesc, msg string) {
	if self.buffer == nil {
		self.doPrintf(level, levelDesc, msg)
		return
	}
	self.buffer <- &Msg{level, levelDesc, msg}
}

func (self *SimpleLogger) setFileOutDir(path string) {
	self.outputPath = path
}

func (self *SimpleLogger) setFileName(fileName string) {
	self.fileName = fileName
}

func (self *SimpleLogger) setFileLogLevel(level int) {
	self.fileLevel = level
}
func (self *SimpleLogger) setShellLogLevel(level int) {
	self.shellLevel = level
}

func (self *SimpleLogger) createLogFile(dir string, fileName string) (*os.File, error) {
	now := time.Now()
	newFileName := fmt.Sprintf("%s-%d-%02d-%02d.log",
		fileName,
		now.Year(),
		now.Month(),
		now.Day())
	if dir != "" {
		_, err := os.Stat(dir)
		if os.IsNotExist(err) {
			os.MkdirAll(dir, os.ModePerm)
		}
	}
	if self.curName != newFileName {
		self.curName = newFileName
		self.curIdx = 0
	}
	if self.curIdx <= 0 {
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			self.curIdx = 1
		} else {
			for _, fi := range files {
				if !fi.IsDir() {
					// 过滤指定格式
					fname := fi.Name()
					ok := strings.HasPrefix(fname, newFileName)
					if ok {
						n := strings.LastIndex(fname, ".")
						fname = fname[n+1:]
						idx, err := strconv.Atoi(fname)
						if err == nil && (idx+1) > self.curIdx {
							self.curIdx = idx + 1
						}
					}
				}
			}
			if self.curIdx <= 0 {
				self.curIdx = 1
			}
		}
	} else {
		self.curIdx = self.curIdx + 1
	}
	newFileName = fmt.Sprintf("%s.%d", newFileName, self.curIdx)
	file, err := os.Create(path.Join(dir, newFileName))
	if err != nil {
		log.Printf("create file failed %v", err)
		return nil, err
	}
	return file, nil
}

func (self *SimpleLogger) Close() {
	for len(self.buffer) > 0 {
		runtime.Gosched()
	}
	close(self.buffer)
	self.wg.Wait()
}

func (self *SimpleLogger) run() {
	go func() {
		self.wg.Add(1)
		for {
			m, ok := <-self.buffer
			if ok {
				self.doPrintf(m.level, m.levelDesc, m.msg)
			} else {
				break
			}
		}
		self.wg.Done()
	}()
}

func CreateLogger(path string, fileName string, fileLevel, shellLevel, maxLine, bufSize int) *SimpleLogger {
	logger := &SimpleLogger{}
	log.SetFlags(log.Ldate | log.Lmicroseconds)
	logger.setFileOutDir(path)
	logger.setFileName(fileName)
	logger.setFileLogLevel(fileLevel)
	logger.setShellLogLevel(shellLevel)
	logger.SetColored(true)
	logger.maxLine = maxLine
	logger.logLine = maxLine + 1
	if bufSize > 10 {
		log.Printf("log start with async mode.")
		logger.buffer = make(chan *Msg, bufSize)
		logger.run()
	}
	return logger
}

var DefaultLogger = CreateLogger("log", "game", FATAL_LEVEL, DEBUG_LEVEL, 10000, 1000)
