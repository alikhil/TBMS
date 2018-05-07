package api

import (
	"fmt"
	io "github.com/alikhil/TBMS/internals/io"
	"github.com/alikhil/distributed-fs/utils"
	"testing"
)

/// Dfs Tests

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
func TestCreateRelationshipWithDfs(t *testing.T) {
	testCreateRelationship(t, getDFS())
}
func TestCreateRelationshipWithPropertiesWithDfs(t *testing.T) {
	testCreateRelationshipWithProperties(t, getDFS())
}
func TestGetLabelsWithDfs(t *testing.T) {
	testGetLabels(t, getDFS())
}
func TestGetNodePropertyWithDfs(t *testing.T) {
	testGetNodeProperty(t, getDFS())
}
func TestGetFromRelationshipWithDfs(t *testing.T) {
	testGetFromRelationship(t, getDFS())
}
func TestGetToRelationshipWithDfs(t *testing.T) {
	testGetToRelationship(t, getDFS())
}
func TestGetRelationshipTypeWithDfs(t *testing.T) {
	testGetRelationshipType(t, getDFS())
}
func TestGetRelationshipPropertyWithDfs(t *testing.T) {
	testGetRelationshipProperty(t, getDFS())
}
func TestSelectNodesWhereWithDfs(t *testing.T) {
	testSelectNodesWhere(t, getDFS())
}
func TestSelectRelationshipsWhereWithDfs(t *testing.T) {
	testSelectRelationshipsWhere(t, getDFS())
}
func TestGetRelationshipsWithDfs(t *testing.T) {
	testGetRelationships(t, getDFS())
}
