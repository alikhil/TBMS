package engine

import (
	"fmt"
	io "github.com/alikhil/TBMS/internals/io"
	"github.com/alikhil/distributed-fs/utils"
	"testing"
)

func getDFS() io.IO {
	var mapa = GetFileToBytesMap()

	client, ok := utils.GetRemoteClient(fmt.Sprintf("%s:5001", utils.GetIPAddress()))
	if !ok {
		panic("failed to connect to remote client")
	}
	dfs := utils.DFSClient{Client: client}
	dfs.InitRecordMappings(mapa)
	return &dfs
}

func TestFindObjectWithDFS(t *testing.T) {
	testFindObject(t, getDFS())
}
func TestCreateAndLoadStringWithDFS(t *testing.T) {
	testCreateAndLoadString(t, getDFS())
}
func TestFindAndCreateWithDFS(t *testing.T) {
	testFindAndCreate(t, getDFS())
}
func TestGetAndLockFreeIDWithDFS(t *testing.T) {
	testGetAndLockFreeID(t, getDFS())
}
func TestInitDatabaseWithDFS(t *testing.T) {
	testInitDatabase(t, getDFS())
}
func TestGetSaveNodeWithDFS(t *testing.T) {
	testGetSaveNode(t, getDFS())
}
func TestGetLabelIDWithDFS(t *testing.T) {
	testGetLabelID(t, getDFS())
}

func TestEngineWithDFS(t *testing.T) {

	var mapa = GetFileToBytesMap()

	client, ok := utils.GetRemoteClient(fmt.Sprintf("%s:5001", utils.GetIPAddress()))
	if !ok {
		t.Error("failed to connect to remote client")
	}
	dfs := utils.DFSClient{Client: client}
	dfs.InitRecordMappings(mapa)

	testPack(t, func() io.IO { return &dfs })
}
