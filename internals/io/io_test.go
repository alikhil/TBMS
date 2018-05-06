package io

import (
	"bytes"
	"fmt"
	"github.com/alikhil/distributed-fs/utils"
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

func completeTest(c IO, t *testing.T) {
	fname := "nodes.store"
	var recordSize byte = 13
	const n byte = 19
	var m = n * recordSize

	data := make([]byte, m)
	var offset int32
	for i := byte(0); i < m; i++ {
		data[i] = i
	}

	_ = c.CreateFile(fname)
	defer c.DeleteFile(fname)

	c.WriteBytes(fname, offset, &data)

	orders := [n]int32{0, 2, 1, 4, 3, 10, 9, 8, 13, 4, 2, 8, 1, 15, 14, 4, 5, 10, 3}

	for j, i := range orders {
		localOffset := i * int32(recordSize)
		to := localOffset + int32(recordSize)
		dt, ok := c.ReadBytes(fname, localOffset, int32(recordSize))
		if !ok {
			t.Fatalf("failed to read range of bytes %d: %d", localOffset, to)
		}

		rdt := data[localOffset:to]

		if !bytes.Equal(dt, rdt) {
			t.Errorf("step %v(%v) expected %v but get %v", j, i, rdt, dt)
		}

	}

}

func TestLocalIOReadWrite(t *testing.T) {
	readWriteFile(LocalIO{}, t)
}

func TestLocalIOCreateFile(t *testing.T) {
	createFile(LocalIO{}, t)
}

func TestLocalIOComplete(t *testing.T) {
	completeTest(LocalIO{}, t)
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

func TestCacheComplete(t *testing.T) {
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
	completeTest(&cache, t)
}

func TestDistributeFSReadWrite(t *testing.T) {
	var mapa = map[string]int32{
		"nodes.store": 13,
	}

	client, ok := utils.GetRemoteClient(fmt.Sprintf("%s:5001", utils.GetIPAddress()))
	if !ok {
		t.Error("failed to connect to remote client")
	}
	dfs := utils.DFSClient{Client: client}
	dfs.InitRecordMappings(&mapa)

	readWriteFile(&dfs, t)
}

func TestDistributeFSComplete(t *testing.T) {
	var mapa = map[string]int32{
		"nodes.store": 13,
	}

	client, ok := utils.GetRemoteClient(fmt.Sprintf("%s:5001", utils.GetIPAddress()))
	if !ok {
		t.Error("failed to connect to remote client")
	}
	dfs := utils.DFSClient{Client: client}
	dfs.InitRecordMappings(&mapa)

	completeTest(&dfs, t)
}

func TestCacheWithDistributeFSReadWrite(t *testing.T) {
	var mapa = map[string]int32{
		"nodes.store": 13,
	}

	client, ok := utils.GetRemoteClient(fmt.Sprintf("%s:5001", utils.GetIPAddress()))
	if !ok {
		t.Error("failed to connect to remote client")
	}
	dfs := utils.DFSClient{Client: client}
	dfs.InitRecordMappings(&mapa)

	cache := LRUCache{}
	cache.Init(&dfs, &mapa, 5)

	readWriteFile(&cache, t)
}

func TestCacheWithDistributeFSComplete(t *testing.T) {
	var mapa = map[string]int32{
		"nodes.store": 13,
	}

	client, ok := utils.GetRemoteClient(fmt.Sprintf("%s:5001", utils.GetIPAddress()))
	if !ok {
		t.Error("failed to connect to remote client")
	}
	dfs := utils.DFSClient{Client: client}
	dfs.InitRecordMappings(&mapa)

	cache := LRUCache{}
	cache.Init(&dfs, &mapa, 5)

	completeTest(&cache, t)
}
