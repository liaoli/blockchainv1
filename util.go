package main

import (
	"bytes"
	"encoding/binary"
	"log"
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
