package io

import (
	"os"
)

type IO interface {
	ReadBytes(file string, offset, count int) (data []byte, ok bool)
	WriteBytes(file string, offset int, bytes *[]byte) (ok bool)
	CreateFile(file string) (ok bool)
	FileExists(file string) bool
	DeleteFile(file string) (ok bool)
}

// there should be 3 types:
// LocalIO, DistributedIO, and Cache

type LRUCache struct {
	baseIO IO
}

func (c LRUCache) ReadBytes(file string, offset, count int) ([]byte, bool) {
	return nil, false
}

func (c LRUCache) WriteBytes(file string, offset int, bytes *[]byte) bool {
	return false
}

func (c LRUCache) CreateFile(file string) bool {
	c.baseIO.CreateFile(file)
	return true
}

func (c LRUCache) FileExists(file string) bool {
	return c.baseIO.FileExists(file)
}

func (c LRUCache) DeleteFile(file string) bool {
	return c.baseIO.DeleteFile(file)
}

type LocalIO struct {
}

func (io LocalIO) ReadBytes(filename string, offset, count int) (data []byte, ok bool) {
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
		return nil, false
	}
	return data, true

}

func (io LocalIO) WriteBytes(filename string, offset int, bytes *[]byte) bool {
	var file, err = os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
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
