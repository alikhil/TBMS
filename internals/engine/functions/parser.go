package engine

import (
	"encoding/binary"
	"bytes"
	"github/alikhil/TBMS/internals/engine"
)

type ENode engine.ENode
type EProperty engine.EProperty

func parseInt(data []byte) (ret int) {
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.LittleEndian, &ret)
	return
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
		NextLabelID:    parseInt((*data)[1:4]),
		NextPropertyID: parseInt((*data)[5:8]),
		NextRelID:      parseInt((*data)[9:12])}
	return &node, true
}

func parseProperty(data *[]byte) (*EProperty, bool) {
	var inUse = parseBool((*data)[0])
	if !inUse {
		return nil, false
	}
	return &EProperty{
		Typename:         engine.EType((*data)[1]),
		KeyStringID:      parseInt((*data)[2:5]),
		ValueOrStringPtr: parseInt((*data)[6:9])}, true
}
