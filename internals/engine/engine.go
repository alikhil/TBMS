package engine

import (
	tuple "github.com/kmanley/golang-tuple"
)

// EType is used for representing property types in database

// types starting with E are used only within Engine
type Engine interface {
	GetObjectIterator() func() []byte

	GetNodesIterator() func() ENode
	GetNodesByLabelIterator(labelID string) func() ENode

	IterateAllRelatioships() func() ERelationship
	IterateRelationships(relTypeID int) func() ERelationship

	IterateObjectProperties(objID int) func() EProperty
	IterateRelationshipProperties(relID int) func() EProperty
	IterateObjectLabels(objID int) func() string

	CreateRelationship(firstObjID, secondObjectID int, relTypeID int) int
	DeleteRelationship(firstObjID, secondObjectID int, relTypeID int) (int, error) // error if no such relationship

	AddObjectLabel(objID, labelID int)
	DeleteObjectLabel(objID, labelID int)
	CreateObject(labels []string, properties []tuple.Tuple) // array of key-values
	AddObjectProperty(objID int, property tuple.Tuple)
	UpdateObjectProperty(objID int, property tuple.Tuple)

	GetLabelByID(int) (string, error) // error if no such labelid
	GetLabelID(string) (int, error)   // error if not exists
	CreateLabel(string) int           // by name, return id

	GetRelationshipTypeID(string) (int, error) // error if no such relationship type
	CreateRelationshipType(string) int         // by type, return id

	// Getting free ids for each store using InUse
}
