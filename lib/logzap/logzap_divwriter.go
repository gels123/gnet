package logzap

import (
	"fmt"
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
	maxSize      int64              // 最大文件大小(B)
	curSize      int64              // 当前文件大小
	maxAge       time.Duration      // 文件过期时长
	rotateTime   time.Duration      // 文件滚动时长
	bNewFile     bool               // 文件是否滚动
	out          *os.File           // 输出文件对象
	timer        *time.Timer        // 文件滚动计时器
	mutex        sync.RWMutex       // 读写锁
}

var patternConversionRegexps = []*regexp.Regexp{
	regexp.MustCompile(`%[%+A-Za-z]`),
	regexp.MustCompile(`\*+`),
}

// New creates a new divWriter object.
func newDivWriter(fileDir, fileName string, maxSize int64, maxAge, rotateTime time.Duration) *divWriter {
	pattern, err := strftime.New(fileName)
	if err != nil {
		panic(errors.Wrap(err, `newDivWriter err: invalid strftime pattern`))
	}
	patternClean := fileName
	for _, re := range patternConversionRegexps {
		patternClean = re.ReplaceAllString(patternClean, "*")
	}
	dw := &divWriter{
		fileDir:      fileDir,
		fileName:     fileName,
		curFileName:  "",
		curFileIdx:   0,
		pattern:      pattern,
		patternClean: patternClean,
		maxSize:      maxSize,
		curSize:      0,
		maxAge:       maxAge,
		rotateTime:   rotateTime,
		bNewFile:     false,
		out:          nil,
		timer:        nil,
	}
	dw.rotateByTime()
	return dw
}

// Write satisfies the io.Writer interface.
func (dw *divWriter) Write(p []byte) (n int, err error) {
	// Guard against concurrent writes
	dw.mutex.Lock()
	defer dw.mutex.Unlock()

	if dw.out == nil || dw.bNewFile {
		out, err := dw.genFile()
		if err != nil {
			return 0, errors.Wrap(err, `Write err`)
		}
		if dw.out != nil {
			dw.out.Close()
		}
		dw.out = out
		dw.curSize = 0
		dw.bNewFile = false
	}
	n, err = dw.out.Write(p)
	dw.curSize += int64(n)
	if dw.curSize >= dw.maxSize {
		dw.bNewFile = true
	}
	return n, err
}

// must be locked during this operation
func (dw *divWriter) genFile() (*os.File, error) {
	filename := dw.genFileName()
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, errors.Wrap(err, "genFile err")
	}
	if err := dw.cleanFile(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
	}
	return file, nil
}

func (dw *divWriter) genFileName() string {
	now := time.Now().UTC()
	now = now.Truncate(time.Duration(dw.rotateTime))
	newFileName := dw.pattern.FormatString(now)
	if dw.curFileName != newFileName {
		dw.curFileName = newFileName
		dw.curFileIdx = 0
	}
	if dw.curFileIdx <= 0 {
		tmp := filepath.Join(dw.fileDir, newFileName+".*")
		matches, err := filepath.Glob(tmp)
		if err != nil || len(matches) <= 0 {
			dw.curFileIdx = 1
		} else {
			tmp = filepath.Join(dw.fileDir, newFileName)
			for _, fname := range matches {
				if isFile(fname) {
					ok := strings.HasPrefix(fname, tmp) // fname begin with tmp
					if ok {
						n := strings.LastIndex(fname, ".")
						if n >= 0 {
							fname = fname[n+1:]
							idx, err := strconv.Atoi(fname)
							if err == nil && (idx+1) > dw.curFileIdx {
								dw.curFileIdx = idx + 1
							}
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
	return filepath.Join(dw.fileDir, newFileName)
}

// clean log files with modify time before maxAge
func (dw *divWriter) cleanFile() error {
	if dw.maxAge > 0 {
		cutoff := time.Now().UTC().Add(-dw.maxAge)
		tmp := filepath.Join(dw.fileDir, dw.patternClean)
		matches, err := filepath.Glob(tmp)
		if err != nil {
			return errors.Wrap(err, "cleanFile err")
		}
		var files []string
		for _, fname := range matches {
			fi, err := os.Stat(fname)
			if err != nil {
				continue
			}
			if fi.ModTime().UTC().After(cutoff) {
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

// Rotate by time
func (dw *divWriter) rotateByTime() {
	t1 := time.Now().UTC().Truncate(time.Second)
	t2 := t1.Truncate(time.Duration(dw.rotateTime))
	t3 := t2.Add(dw.rotateTime).Sub(t1) // t2 + dw.rotateTime - t1
	dw.timer = time.AfterFunc(t3, func() {
		dw.Rotate()
		dw.rotateByTime()
	})
}

// Call for a rotation
func (dw *divWriter) Rotate() {
	dw.mutex.Lock()
	defer dw.mutex.Unlock()
	dw.bNewFile = true
}

// Close divWriter
func (dw *divWriter) Close() error {
	dw.mutex.Lock()
	defer dw.mutex.Unlock()

	if dw.out != nil {
		dw.out.Close()
		dw.out = nil
	}
	if dw.timer != nil {
		dw.timer.Stop()
		dw.timer = nil
	}
	return nil
}
