package io

import (
	"os"
)

type IO interface {
	ReadBytes(file string, offset, count int32) (data []byte, ok bool)
	WriteBytes(file string, offset int32, bytes *[]byte) (ok bool)
	CreateFile(file string) (ok bool)
	FileExists(file string) bool
	DeleteFile(file string) (ok bool)
}

// there should be 3 types:
// LocalIO, DistributedIO, and Cache

type LocalIO struct {
}

func (io LocalIO) ReadBytes(filename string, offset, count int32) (data []byte, ok bool) {
	if !io.FileExists(filename) {
		return nil, false
	}
	var file, err = os.Open(filename)
	defer file.Close()
	if err != nil {
		return nil, false
	}

	data = make([]byte, count, count)
	var _, er = file.ReadAt(data, int64(offset))
	if er != nil {
		if er.Error() == "EOF" {
			resultData := make([]byte, 0, count)
			for i := 0; int32(i) <= int32(count)/recordSize; i++ {
				data = make([]byte, recordSize, recordSize)
				var _, error = file.ReadAt(data, int64(offset+recordSize*int32(i)))
				if error != nil {
					data = make([]byte, cap(resultData)-len(resultData), cap(resultData)-len(resultData))
						return append(resultData, (data)...), true
				}
				resultData = append(resultData, (data)...)
			}
		}
		return nil, false
	}
	return data, true

}

func (io LocalIO) WriteBytes(filename string, offset int32, bytes *[]byte) bool {
	var file, err = os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return false
	}
	_, err = file.WriteAt(*bytes, int64(offset))
	return err == nil

}

func (io LocalIO) CreateFile(filename string) bool {
	var _, err = os.Create(filename)
	return err == nil
}

func (LocalIO) FileExists(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}

func (io LocalIO) DeleteFile(filename string) bool {
	if io.FileExists(filename) {
		var err = os.Remove(filename)
		return err == nil
	}
	return false
}
