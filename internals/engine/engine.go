package engine

import (
	tuple "github.com/kmanley/golang-tuple"
)

// EType is used for representing property types in database

// types starting with E are used only within Engine
type Engine interface {
	GetObjectIterator() func() ([]byte, bool)

	GetNodesIterator() func() (*ENode, bool)
	GetNodesByLabelIterator(labelID string) func() (*ENode, bool)

	GetRelationshiptIterator() func() (*ERelationship, bool)
	GetRelationshiptIteratorByType(relTypeID int) func() (*ERelationship, bool)

	GetObjectPropertiesIterator(objID int) func() (*EProperty, bool)
	GetRelationshipPropertiesIterator(relID int) func() (*EProperty, bool)
	GetObjectLabelsIterator(objID int) func() (labelID int, ok bool)

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
