package io

type IO interface {
	readBytes(file string, offset, count int) []byte
	writeBytes(file string, offset int, bytes []byte)
	createFile(file string)
	fileExists(file string) bool
}

// there should be 3 types:
// LocalIO, DistributedIO, and Cache

type LRUCache struct {
	baseIO IO
}

func (c *LRUCache) readBytes(file string, offset, count int) []byte {
	panic("not implemented exception")
}

func (c *LRUCache) writeBytes(file string, offset int, bytes []byte) {
	panic("not implemented exception")
}

func (c *LRUCache) createFile(file string) {
	c.baseIO.createFile(file)
}

func (c *LRUCache) fileExists(file string) bool {
	return c.baseIO.fileExists(file)
}

type LocalIO struct {
}

func (c LocalIO) readBytes(file string, offset, count int) []byte {
	panic("not implemented exception")
}

func (c LocalIO) writeBytes(file string, offset int, bytes []byte) {
	panic("not implemented exception")
}

func (c LocalIO) createFile(file string) {
	panic("not implemented exception")
}

func (c LocalIO) fileExists(file string) bool {
	panic("not implemented exception")
}
