/*
 * 与文件相关的通用函数
 */
package logzap

import (
	"os"

	"github.com/pkg/errors"
)

// 获取当前工作目录路径
func getCurDir() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return dir
}

// 判断所给路径是否为文件夹
func isDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// 判断所给路径是否为文件
func isFile(path string) bool {
	return !isDir(path)
}

// 创建文件夹
func createDir(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(path, os.ModePerm)
			if err != nil {
				return errors.Wrap(err, "createDir err")
			}
		} else {
			return errors.Wrap(err, "createDir err")
		}
	} else {
		if !s.IsDir() {
			return errors.New("createDir err")
		}
	}
	return nil
}
