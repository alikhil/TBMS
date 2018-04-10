package engine

import "github.com/alikhil/TBMS/internals/io"

type EType byte

type RealEngine struct {
	io.IO
}

type EIterator struct {
	ID int
}

type ELabel struct {
	ID int
}

type ENode struct {
	ID             int
	NextLabelID    int
	NextPropertyID int
	NextRelID      int
}

type EParser struct {
}

type EProperty struct {
	ID               int
	ValueOrStringPtr int
	Typename         EType
	KeyStringID      int
}

type ERelationship struct {
	ID           int
	FirstInChain bool

	FirstNodeID int
	SecondNode  int

	FirstNodeNxtRelID  int
	SecondNodeNxtRelID int
	FirstNodePrvRelID  int
	SecondNodePrvRelID int
	NextPropertyID     int
	TypeID             int
}

type ELabelString struct {
	ID     int
	String string
}
