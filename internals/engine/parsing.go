package engine

import (
	"bytes"
	"encoding/binary"
	"github.com/alikhil/TBMS/internals/logger"
)

func parseInt(data []byte) int {
	var ret int32
	buf := bytes.NewBuffer(data)
	err := binary.Read(buf, ConventionByteOrder, &ret)
	if err != nil {
		logger.Error.Printf("Can not parse int %v", err)
	}
	return int(ret)
}

func parseBool(b byte) bool {
	return b > 0
}

func parseNode(data *[]byte, nodeID int) (*ENode, bool) {

	var inUse = parseBool((*data)[0])
	if !inUse {
		return nil, false
	}
	var node = ENode{
		ID:             nodeID,
		NextLabelID:    parseInt((*data)[1:5]),
		NextPropertyID: parseInt((*data)[5:9]),
		NextRelID:      parseInt((*data)[9:13])}
	return &node, true
}

func parseProperty(data *[]byte) (*EProperty, bool) {
	var inUse = parseBool((*data)[0])
	if !inUse {
		return nil, false
	}
	return &EProperty{
		Typename:         EType((*data)[1]),
		KeyStringID:      parseInt((*data)[2:6]),
		ValueOrStringPtr: parseInt((*data)[6:10])}, true
}

func parseLabelString(data *[]byte, id int) (*ELabelString, bool) {
	var inUse = parseBool((*data)[0])
	if !inUse {
		return nil, false
	}

	end := bytes.IndexByte(*data, 0)
	s := string((*data)[1:end])
	return &ELabelString{ID: id, String: s}, true
}

func parseRelationship(data *[]byte, id int) (*ERelationship, bool) {
	var inUse = parseBool((*data)[0])
	if !inUse {
		return nil, false
	}

	return &ERelationship{
		ID:                 id,
		FirstInChain:       parseBool((*data)[1]),
		SecondNode:         parseInt((*data)[2:6]),
		FirstNodeID:        parseInt((*data)[6:10]),
		FirstNodeNxtRelID:  parseInt((*data)[10:14]),
		SecondNodeNxtRelID: parseInt((*data)[14:18]),
		FirstNodePrvRelID:  parseInt((*data)[18:22]),
		SecondNodePrvRelID: parseInt((*data)[22:26]),
		NextPropertyID:     parseInt((*data)[26:30]),
		TypeID:             parseInt((*data)[30:34]),
	}, true
}
