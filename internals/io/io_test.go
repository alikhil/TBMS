package io

import (
	"bytes"
	"testing"
)

func createFile(c IO, t *testing.T) {
	c.CreateFile("nodes.store")
	defer c.DeleteFile("nodes.store")

	if !c.FileExists("nodes.store") {
		t.Errorf("file does not created")
	}
}

func readWriteFile(c IO, t *testing.T) {
	fname := "nodes.store"
	data := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	var offset int32 = 0

	_ = c.CreateFile(fname)
	defer c.DeleteFile(fname)

	c.WriteBytes(fname, offset, &data)

	readData, _ := c.ReadBytes(fname, offset, 13)

	if !bytes.Equal(data, readData) {
		t.Errorf("Expected %v but got %v", data, readData)
	}
}

func TestLocalIOReadWrite(t *testing.T) {
	readWriteFile(LocalIO{}, t)
}

func TestLocalIOCreateFile(t *testing.T) {
	createFile(LocalIO{}, t)
}

func TestCacheReadWrite(t *testing.T) {
	var mapa = map[string]int32{
		"nodes.store":             13,
		"labels.store":            9,
		"labelsStrings.store":     21,
		"relationships.store":     34,
		"properties.store":        14,
		"strings.store":           64,
		"inuse.store":             11,
		"propertykeys.store":      21,
		"relationshiptypes.store": 34,
	}
	cache := LRUCache{}
	cache.Init(LocalIO{}, &mapa, 20)
	readWriteFile(&cache, t)
}

func TestCacheCreateFile(t *testing.T) {
	var mapa = map[string]int32{
		"nodes.store":             13,
		"labels.store":            9,
		"labelsStrings.store":     21,
		"relationships.store":     34,
		"properties.store":        14,
		"strings.store":           64,
		"inuse.store":             11,
		"propertykeys.store":      21,
		"relationshiptypes.store": 34,
	}
	cache := LRUCache{}
	cache.Init(LocalIO{}, &mapa, 2)
	createFile(&cache, t)
}
