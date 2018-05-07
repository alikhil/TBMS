package api

import (
	en "github.com/alikhil/TBMS/internals/engine"
	io "github.com/alikhil/TBMS/internals/io"
	"testing"
)

/// Cache Tests

func getCache() io.IO {
	var mapa = en.GetFileToBytesMap()

	cache := io.LRUCache{}
	cache.Init(io.LocalIO{}, mapa, 5)
	return &cache
}

func TestCreateRelationshipWithCache(t *testing.T) {
	testCreateRelationship(t, getCache())
}
func TestCreateRelationshipWithPropertiesWithCache(t *testing.T) {
	testCreateRelationshipWithProperties(t, getCache())
}
func TestGetLabelsWithCache(t *testing.T) {
	testGetLabels(t, getCache())
}
func TestGetNodePropertyWithCache(t *testing.T) {
	testGetNodeProperty(t, getCache())
}
func TestGetFromRelationshipWithCache(t *testing.T) {
	testGetFromRelationship(t, getCache())
}
func TestGetToRelationshipWithCache(t *testing.T) {
	testGetToRelationship(t, getCache())
}
func TestGetRelationshipTypeWithCache(t *testing.T) {
	testGetRelationshipType(t, getCache())
}
func TestGetRelationshipPropertyWithCache(t *testing.T) {
	testGetRelationshipProperty(t, getCache())
}
func TestSelectNodesWhereWithCache(t *testing.T) {
	testSelectNodesWhere(t, getCache())
}
func TestSelectRelationshipsWhereWithCache(t *testing.T) {
	testSelectRelationshipsWhere(t, getCache())
}
func TestGetRelationshipsWithCache(t *testing.T) {
	testGetRelationships(t, getCache())
}
