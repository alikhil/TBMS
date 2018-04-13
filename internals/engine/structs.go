package engine

import "github.com/alikhil/TBMS/internals/io"

type RealEngine struct {
	io.IO
}

type EIterator struct {
	ID int32
}

type ELabel struct {
	ID int32
}

type ENode struct {
	ID             int32
	NextLabelID    int32
	NextPropertyID int32
	NextRelID      int32
}

type EParser struct {
}

type EProperty struct {
	ID               int32
	ValueOrStringPtr int32
	Typename         EType
	KeyStringID      int32
}

type ERelationship struct {
	ID           int32
	FirstInChain bool

	FirstNodeID int32
	SecondNode  int32

	FirstNodeNxtRelID  int32
	SecondNodeNxtRelID int32
	FirstNodePrvRelID  int32
	SecondNodePrvRelID int32
	NextPropertyID     int32
	TypeID             int32
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
