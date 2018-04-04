package engine

import (
	tuple "github.com/kmanley/golang-tuple"
)

// types starting with E are used only within Engine

type Engine interface {
	IterateObjects(labelID int) func() EObject
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
	GetLabelId(string) (int, error)   // error if not exists
	CreateLabel(string) int           // by name, return id

	GetRelationshipTypeID(string) (int, error) // error if no such relationship type
	CreateRelationshipType(string) int         // by type, return id
}

type EObject struct {
	id int
	// nextLabelId?
	// labels []string
	// properties []EProperty
}

func (o EObject) GetId() int {
	return o.id
}

func (o EObject) SetId(newId int) {
	o.ID = newId
}

type EProperty struct {
	ID  int
	key string
	// propertyId?
	// getInt? getBool? getString?
}

type ERelationship struct {
	ID             int
	firstObjectID  int
	secondObjectID int
	typeID         int
}
