package engine

import (
	"fmt"
	io "github.com/alikhil/TBMS/internals/io"
	"github.com/alikhil/distributed-fs/utils"
	"testing"
)

func getDFSWithCache() io.IO {

	var mapa = GetFileToBytesMap()

	client, ok := utils.GetRemoteClient(fmt.Sprintf("%s:5001", utils.GetIPAddress()))
	if !ok {
		panic("failed to connect to remote client")
	}
	dfs := utils.DFSClient{Client: client}
	dfs.InitRecordMappings(mapa)

	cache := io.LRUCache{}
	cache.Init(&dfs, mapa, 5)

	return &cache
}

func TestEngineWithDFSAndCache(t *testing.T) {

	testPack(t, getDFSWithCache)
}
