package utils

import (
	"io/ioutil"
	"os"
)

// 读取指定文件
func ReadFile(addr string) ([]byte, error) {
	file, err := os.Open(addr)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return ioutil.ReadAll(file)
}

// 创建指定目录, 只在不存在时创建
func Mkdir(addr string) {
	s, err := os.Stat(addr)
	if err != nil || !s.IsDir() {
		os.Mkdir(addr, os.ModePerm)
	}
}
