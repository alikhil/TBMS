package io

import (
	"os"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func ReadFromFile(file *os.File, offset, count int) (data []byte, ok bool) {
	data = make([]byte, count, count)
	curCnt := 0
	for curCnt < count {
		var lastRead []byte
		var readCnt, err = file.ReadAt(lastRead, int64(offset))
		if err == nil {
			return nil, false
		}
		n := min(readCnt, count-curCnt)
		data = append(data, lastRead[:n]...)
		curCnt += n
	}

	return data, len(data) == count
}
