package engine

import (
	"bytes"
	"encoding/binary"
)

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
