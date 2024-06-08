package util

import (
	"os"
)

// CreateFolder 创建文件夹
func CreateFolder(path string) error {
	if !PathOrFileExist(path) {
		return os.MkdirAll(path, os.ModePerm)
	}
	return nil
}

// GetExecutablePath 获取应用主目录
func GetExecutablePath() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	return dir + string(os.PathSeparator)
}

// PathOrFileExist 判断文件目录或文件是否存在
func PathOrFileExist(target string) bool {
	if _, err := os.Stat(target); os.IsNotExist(err) {
		return false
	}
	return true
}
