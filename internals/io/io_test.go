package io

import (
	"bytes"
	"testing"
)

func createFile(c IO, t *testing.T) {
	c.CreateFile("test")
	defer c.DeleteFile("test")

	if !c.FileExists("test") {
		t.Errorf("file does not created")
	}
}

func readWriteFile(c IO, t *testing.T) {
	fname := "testfile"
	data := []byte{0, 1, 2, 3, 4}
	offset := 42

	_ = c.CreateFile(fname)
	defer c.DeleteFile(fname)

	c.WriteBytes(fname, offset, data)

	readData, _ := c.ReadBytes(fname, offset, 5)

	if !bytes.Equal(data, readData) {
		t.Errorf("Read data is not equal to that which was written")
	}
}

func TestLocalIOReadWrite(t *testing.T) {
	readWriteFile(LocalIO{}, t)
}

func TestLocalIOCreateFile(t *testing.T) {
	createFile(LocalIO{}, t)
}

func TestCacheReadWrite(t *testing.T) {
	readWriteFile(LRUCache{LocalIO{}}, t)
}

func TestCacheCreateFile(t *testing.T) {
	createFile(LRUCache{LocalIO{}}, t)
}
