package engine

import "github/alikhil/TBMS/internals/io"

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
	firstInChain bool

	firstNodeID int
	secondNode  int

	firstNodeNxtRelID  int
	secondNodeNxtRelID int
	firstNodePrvRelID  int
	secondNodePrvRelID int
	typeID             int
}
