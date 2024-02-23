package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"os"
)

func uint64ToByte(num uint64) []byte {
	var buffer bytes.Buffer
	//将二进制形式数据存到 buffer中
	err := binary.Write(&buffer, binary.BigEndian, num)

	if err != nil {
		log.Panic(err)
	}

	return buffer.Bytes()
}

//判断文件是否存在
func isFileExist(filename string) bool {
	// func Stat(name string) (FileInfo, error) {
	_, err := os.Stat(filename)

	//os.IsExist不要使用，不可靠
	if os.IsNotExist(err) {
		return false
	}

	return true
}
