package utils

import (
	"github.com/pkg/errors"
	"os"
	"path/filepath"
)

// 获取可执行文件路径
func GetExeDir() string {
	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	return path
}

// 获取当前工作目录路径
func GetCurDir() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return dir
}

// 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// 判断所给路径是否为文件
func IsFile(path string) bool {
	return !IsDir(path)
}

// 创建文件夹
func CreateDir(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(path, os.ModePerm)
			if err != nil {
				return errors.Wrap(err, "CreateDir err")
			}
		} else {
			return errors.Wrap(err, "CreateDir err")
		}
	} else {
		if !s.IsDir() {
			return errors.New("CreateDir err")
		}
	}
	return nil
}
