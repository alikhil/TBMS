package engine

import (
	"bytes"
	"fmt"
	"github.com/alikhil/TBMS/internals/io"
)

// TODO: not sure about storing it here
type RealEngine struct {
	io.IO
}

// EObject - Common interface for all database objects
type EObject interface {
	getID() int32
	getStore() EStore
	encode() *[]byte
	fill(*[]byte, int32)
}

type EString struct {
	ID         int32
	Extra      byte
	Value      *[]byte
	NextPartID int32
}

type ELabel struct {
	ID            int32
	LabelStringID int32
	NextLabelID   int32
}

type EPropertyKey struct {
	ID        int32
	KeyString string
}

type ERelationshipType struct {
	ID         int32
	TypeString string
}

type ENode struct {
	ID             int32
	NextLabelID    int32
	NextPropertyID int32
	NextRelID      int32
}

type EProperty struct {
	ID               int32
	ValueOrStringPtr interface{}
	Typename         EType
	KeyStringID      int32
	NextPropertyID   int32
}

type ERelationship struct {
	ID           int32
	FirstInChain bool //TODO: how it's used

	FirstNodeID  int32
	SecondNodeID int32

	FirstNodeNxtRelID  int32
	SecondNodeNxtRelID int32
	FirstNodePrvRelID  int32
	SecondNodePrvRelID int32
	NextPropertyID     int32
	TypeID             int32
}

type ERelPart struct {
	NodeNxtRelID  int32
	NodePrevRelID int32
}

type ELabelString struct {
	ID     int32
	String string
}

type EInUseRecord struct {
	ID           int32
	StoreType    EStore
	IsHead       bool
	ObjID        int32
	NextRecordID int32
}

// GetPart returns needed part of relationship for node
func (rel *ERelationship) GetPart(nodeID int32) *ERelPart {
	if rel.FirstNodeID == nodeID {
		return &ERelPart{
			NodeNxtRelID:  rel.FirstNodeNxtRelID,
			NodePrevRelID: rel.FirstNodePrvRelID,
		}
	}
	return &ERelPart{
		NodeNxtRelID:  rel.SecondNodeNxtRelID,
		NodePrevRelID: rel.SecondNodePrvRelID,
	}
}

func (rel *ERelationship) SetNextRelationshipID(nodeID, relID int32) {
	if rel.FirstNodeID == nodeID {
		rel.FirstNodeNxtRelID = relID
	} else if rel.SecondNodeID == nodeID {
		rel.SecondNodeNxtRelID = relID
	} else {
		panic(fmt.Sprintf("there is no such nodeID(%v) in relationship(%v)", nodeID, *rel))
	}
}

// Store getters

func (*ERelationship) getStore() EStore {
	return StoreRelationship
}

func (*ENode) getStore() EStore {
	return StoreNode
}

func (*EInUseRecord) getStore() EStore {
	return StoreInUse
}

func (*ELabel) getStore() EStore {
	return StoreLabel
}

func (*ELabelString) getStore() EStore {
	return StoreLabelString
}

func (*EProperty) getStore() EStore {
	return StoreProperty
}

func (*EPropertyKey) getStore() EStore {
	return StorePropertyKey
}

func (*ERelationshipType) getStore() EStore {
	return StoreRelationshipType
}

func (*EString) getStore() EStore {
	return StoreString
}

// ID getters

func (o *ERelationship) getID() int32 {
	return o.ID
}

func (o *ENode) getID() int32 {
	return o.ID
}

func (o *EInUseRecord) getID() int32 {
	return o.ID
}

func (o *ELabel) getID() int32 {
	return o.ID
}

func (o *ELabelString) getID() int32 {
	return o.ID
}

func (o *EProperty) getID() int32 {
	return o.ID
}

func (o *EPropertyKey) getID() int32 {
	return o.ID
}

func (o *ERelationshipType) getID() int32 {
	return o.ID
}

func (o *EString) getID() int32 {
	return o.ID
}

// String functions
func (s *EString) LoadString(re *RealEngine) string {
	end := bytes.IndexByte(*s.Value, 0)
	if end == -1 {
		end = len(*s.Value)
	}
	str := string((*s.Value)[:end])

	for s.NextPartID != -1 {
		s.ID = s.NextPartID
		ok := re.GetObject(s)
		if !ok {
			panic("can not load str " + string(s.ID))
		}
		end = bytes.IndexByte(*s.Value, 0)
		if end == -1 {
			end = len(*s.Value)
		}
		str += string((*s.Value)[:end])
	}
	return str
}

func (re *RealEngine) CreateStringAndReturnFirstChunk(s string) *EString {
	data := []byte(s)
	lb := 58
	var nxtPart *EString = &EString{ID: -1}
	for len(data) > lb {
		id, ok := re.GetAndLockFreeIDForStore(StoreString)
		if !ok {
			panic("can not create string: " + s)
		}
		var part = data[len(data)-lb:]
		data = data[:len(data)-lb]
		nxtPart = &EString{ID: id, Value: &part, NextPartID: nxtPart.ID}
		saved := re.SaveObject(nxtPart)
		if !saved {
			panic("failed to save parts of str - " + s)
		}
	}
	id, ok := re.GetAndLockFreeIDForStore(StoreString)
	if !ok {
		panic("can not create string: " + s)
	}
	head := &EString{ID: id, Value: &data, NextPartID: nxtPart.ID}
	re.SaveObject(head)
	return head
}
