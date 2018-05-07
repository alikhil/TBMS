package engine

import (
	"fmt"
	io "github.com/alikhil/TBMS/internals/io"
	"github.com/alikhil/distributed-fs/utils"
	"testing"
)

func getDFS() io.IO {
	var mapa = map[string]int32{
		"nodes.store":             13,
		"labels.store":            9,
		"labelsStrings.store":     21,
		"relationships.store":     34,
		"properties.store":        14,
		"strings.store":           64,
		"inuse.store":             11,
		"propertykeys.store":      21,
		"relationshiptypes.store": 21,
	}

	client, ok := utils.GetRemoteClient(fmt.Sprintf("%s:5001", utils.GetIPAddress()))
	if !ok {
		panic("failed to connect to remote client")
	}
	dfs := utils.DFSClient{Client: client}
	dfs.InitRecordMappings(&mapa)
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

	var mapa = map[string]int32{
		"nodes.store":             13,
		"labels.store":            9,
		"labelsStrings.store":     21,
		"relationships.store":     34,
		"properties.store":        14,
		"strings.store":           64,
		"inuse.store":             11,
		"propertykeys.store":      21,
		"relationshiptypes.store": 21,
	}

	client, ok := utils.GetRemoteClient(fmt.Sprintf("%s:5001", utils.GetIPAddress()))
	if !ok {
		t.Error("failed to connect to remote client")
	}
	dfs := utils.DFSClient{Client: client}
	dfs.InitRecordMappings(&mapa)

	testPack(t, func() io.IO { return &dfs })
}
