package logzap

import (
	"fmt"
	"gnet/lib/utils"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/lestrrat-go/strftime"
	"github.com/pkg/errors"
)

// divWriter represents a log file that gets automatically rotated as you write to it.
type divWriter struct {
	fileDir      string             // 文件路径
	fileName     string             // 文件名称
	curFileName  string             // 当前文件名称
	curFileIdx   int                // 当前文件索引
	pattern      *strftime.Strftime // 文件名称正则fileName
	patternClean string             // 文件清理正则
	maxSize      int64              // 最大文件大小
	curSize      int64              // 当前文件大小
	maxAge       time.Duration      // 文件过期时间
	rotateTime   time.Duration      // 当前文件滚动时长
	bNewFile     bool               // 文件是否滚动
	out          *os.File           // 输出文件对象
	clock        Clock              // 时钟
	mutex        sync.RWMutex       // 读写锁
}

// New creates a new divWriter object.
func newDivWriter(filePath, fileName string, maxSize int64) *divWriter {
	globPattern := fileName
	for _, re := range patternConversionRegexps {
		globPattern = re.ReplaceAllString(globPattern, "*")
	}

	pattern, err := strftime.New(fileName)
	if err != nil {
		panic(errors.Wrap(err, `newDivWriter err: invalid strftime pattern`))
	}

	return &divWriter{
		fileDir:     filePath,
		fileName:    fileName,
		curFileName: "",
		curFileIdx:  0,
		pattern:     pattern,
		bNewFile:    false,
		maxSize:     maxSize,
		curSize:     0,
		out:         nil,

		clock:        clockFunc(time.Now),
		patternClean: globPattern,

		maxAge:     time.Hour * 24 * 7,
		rotateTime: time.Hour * 24,
	}
}

// Write satisfies the io.Writer interface.
func (dw *divWriter) Write(p []byte) (n int, err error) {
	// Guard against concurrent writes
	dw.mutex.Lock()
	defer dw.mutex.Unlock()

	if dw.out == nil || dw.bNewFile {
		out, err := dw.genFile(false, false)
		if err != nil {
			return 0, errors.Wrap(err, `genFile err`)
		}
		dw.out = out
	}
	n, err = dw.out.Write(p)
	dw.curSize += int64(n)
	if dw.curSize >= dw.maxSize {

	}
	return n, err
}

// must be locked during this operation
func (dw *divWriter) genFile(bailOnRotateFail, useGenerationalNames bool) (io.Writer, error) {
	filename := dw.genFileName()
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, errors.Wrap(err, `os.OpenFile err`)
	}

	if err := dw.cleanFile(filename); err != nil {
		err = errors.Wrap(err, "failed to rotate")
		if bailOnRotateFail {
			// Failure to rotate is a problem, but it's really not a great
			// idea to stop your application just because you couldn't rename
			// your log.
			//
			// We only return this error when explicitly needed (as specified by bailOnRotateFail)
			//
			// However, we *NEED* to close `file` here
			if file != nil { // probably can't happen, but being paranoid
				file.Close()
			}
			return nil, err
		}
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
	}
	if dw.out != nil {
		dw.out.Close()
	}
	dw.out = file

	return file, nil
}

func (dw *divWriter) genFileName() string {
	now := dw.clock.Now()
	now = now.Truncate(time.Duration(dw.rotateTime))
	newFileName := dw.pattern.FormatString(now)
	if dw.curFileName != newFileName {
		dw.curFileName = newFileName
		dw.curFileIdx = 0
	}
	if dw.curFileIdx <= 0 {
		tmp := filepath.Join(filePath, newFileName+".*")
		matches, err := filepath.Glob(tmp)
		if err != nil || len(matches) <= 0 {
			dw.curFileIdx = 1
		} else {
			for _, fname := range matches {
				if !utils.IsDir(fname) {
					ok := strings.HasPrefix(fname, newFileName) // 过滤指定格式
					if ok {
						n := strings.LastIndex(fname, ".")
						fname = fname[n+1:]
						idx, err := strconv.Atoi(fname)
						if err == nil && (idx+1) > dw.curFileIdx {
							dw.curFileIdx = idx + 1
						}
					}
				}
			}
			if dw.curFileIdx <= 0 {
				dw.curFileIdx = 1
			}
		}
	} else {
		dw.curFileIdx++
	}
	newFileName = fmt.Sprintf("%s.%d", newFileName, dw.curFileIdx)
	return filepath.Join(filePath, newFileName)
}

// return the current file name
func (dw *divWriter) CurrentFileName() string {
	dw.mutex.RLock()
	defer dw.mutex.RUnlock()
	return dw.curFileName
}

var patternConversionRegexps = []*regexp.Regexp{
	regexp.MustCompile(`%[%+A-Za-z]`),
	regexp.MustCompile(`\*+`),
}

// call for a force rotation
func (dw *divWriter) Rotate() error {
	dw.mutex.Lock()
	defer dw.mutex.Unlock()

	if _, err := dw.genFile(true, true); err != nil {
		return err
	}
	return nil
}

// clean log files with modify time before maxAge
func (dw *divWriter) cleanFile(filename string) error {
	if dw.maxAge > 0 {
		cutoff := dw.clock.Now().Add(-dw.maxAge)
		matches, err := filepath.Glob(dw.patternClean)
		if err != nil {
			return err
		}
		var files []string
		for _, fname := range matches {
			fi, err := os.Stat(fname)
			if err != nil {
				continue
			}
			if fi.ModTime().After(cutoff) {
				continue
			}
			files = append(files, fname)
		}
		if len(files) > 0 {
			// remove files on a separate goroutine
			go func() {
				for _, fname := range files {
					os.Remove(fname)
				}
			}()
		}
	}

	return nil
}

// Close divWriter
func (dw *divWriter) Close() error {
	dw.mutex.Lock()
	defer dw.mutex.Unlock()

	if dw.out != nil {
		dw.out.Close()
		dw.out = nil
	}
	return nil
}
