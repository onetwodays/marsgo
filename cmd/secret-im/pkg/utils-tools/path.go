package utils

import (
	"os"
	"path/filepath"
)

// 当前路径
func CurrentPath() string {
	return filepath.Dir(os.Args[0])
}

// 路径是否存在
func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
