package io

type IO interface {
	ReadBytes(file string, offset, count int) []byte
	WriteBytes(file string, offset int, bytes []byte)
	CreateFile(file string)
	FileExists(file string) bool
	DeleteFile(file string)
}

// there should be 3 types:
// LocalIO, DistributedIO, and Cache

type LRUCache struct {
	baseIO IO
}

func (c LRUCache) ReadBytes(file string, offset, count int) []byte {
	panic("not implemented exception")
}

func (c LRUCache) WriteBytes(file string, offset int, bytes []byte) {
	panic("not implemented exception")
}

func (c LRUCache) CreateFile(file string) {
	c.baseIO.CreateFile(file)
}

func (c LRUCache) FileExists(file string) bool {
	return c.baseIO.FileExists(file)
}

func (c LRUCache) DeleteFile(file string) {
	c.baseIO.DeleteFile(file)
}

type LocalIO struct {
}

func (c LocalIO) ReadBytes(file string, offset, count int) []byte {
	panic("not implemented exception")
}

func (c LocalIO) WriteBytes(file string, offset int, bytes []byte) {
	panic("not implemented exception")
}

func (c LocalIO) CreateFile(file string) {
	panic("not implemented exception")
}

func (c LocalIO) FileExists(file string) bool {
	panic("not implemented exception")
}

func (c LocalIO) DeleteFile(file string) {
	panic("not implemented exception")
}
