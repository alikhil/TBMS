package engine

import (
	"bytes"
	"encoding/binary"
	io "github.com/alikhil/TBMS/internals/io"
	tuple "github.com/kmanley/golang-tuple"
)

// Bytes per record in store
const (
	BytesPerNode         = 13
	BytesPerRelationship = 34
	BytesPerProperty     = 10
	BytesPerString       = 64
	BytesPerLabel        = 9
	BytesPerLabelString  = 21
	BytesPerPropertyKey  = 21
	BytesPerRelType      = 21
	BytesPerInUse        = 10
)

// Names of files
const (
	FNNodes         = "nodes"
	FNRelationships = "relationships"
	FNProperties    = "properties"
	FNStrings       = "strings"
	// TODO: add others
)

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

type RealEngine struct {
	io.IO
}

func (re *RealEngine) GetLabelID(label string) (int, bool) {
	return 0, false
}

func (re *RealEngine) GetLabelIteratorFromId(labelID int) func() (int, bool) {
	return func() (int, bool) {
		return 0, false
	}
}

func (re *RealEngine) GetNodesByLabelIterator(label string) func() (*ENode, bool) {
	nextNode := re.GetNodesIterator()
	neededlabelID, ok := re.GetLabelID(label)
	if !ok {

		return func() (*ENode, bool) {
			return nil, false
		}
	}

	return func() (*ENode, bool) {
		node, ok := nextNode()
		if ok {
			nextLabel := re.GetLabelIteratorFromId(node.nextLabelID)
			for labelID, ok := nextLabel(); ok; {
				if labelID == neededlabelID {
					return node, true
				}
			}
			return nil, false
		}
		return nil, false
	}
}

func (re *RealEngine) GetNodesIterator() func() (*ENode, bool) {
	next := re.GetObjectIterator(FNNodes, BytesPerNode)
	i := 0
	return func() (*ENode, bool) {
		data, ok := next()
		if ok {
			// TODO: add normal parsing
			node := ENode{ID: i, nextLabelID: parseInt(data[1:4])}
			i++
			return &node, ok
		}
		return nil, false
	}
}

func parseInt(data []byte) (ret int) {
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.LittleEndian, &ret)
	return
}

func (re *RealEngine) GetObjectIterator(filename string, recordLength int) func() ([]byte, bool) {
	curOffset := 0
	return func() (data []byte, ok bool) {
		data, ok = re.IO.ReadBytes(filename, curOffset, recordLength)
		if ok {
			curOffset += BytesPerNode
		}
		return
	}
}

type ENode struct {
	ID          int
	nextLabelID int
	// labels []string
	// properties []EProperty
}

func (o ENode) GetId() int {
	return o.ID
}

func (o ENode) SetId(newId int) {
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
