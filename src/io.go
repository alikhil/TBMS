package engine

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

func (c *LRUCache) readBytes() []byte {
	panic("not implemented exception")
}

func (c *LRUCache) writeBytes() {
	panic("not implemented exception")
}

func (c *LRUCache) createFile(file string) {
	c.baseIO.createFile(file)
}

func (c *LRUCache) fileExists(file string) bool {
	return c.baseIO.fileExists(file)
}
