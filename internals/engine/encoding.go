package engine

import (
	"bytes"
	"encoding/binary"
	"github.com/alikhil/TBMS/internals/logger"
)

func encodeInt(val int) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, ConventionByteOrder, int32(val))
	if err != nil {
		logger.Error.Printf("Error on encoding int: %v", err)
	}
	return buf.Bytes()
}

func encodeNode(node *ENode) *[]byte {
	buffer := []byte{1} // Encode InUse

	buffer = append(buffer, encodeInt(node.NextLabelID)...)    // Add NextLabelID
	buffer = append(buffer, encodeInt(node.NextPropertyID)...) // Add NextProprertyID
	buffer = append(buffer, encodeInt(node.NextRelID)...)      // Add NextRelID

	return &buffer
}

func encodeRelationship(rel *ERelationship) *[]byte {
	// TODO: use encodeNode as base and look and SPEC.md
	panic("not implemented")
}

func encodeProperty(prop *EProperty) *[]byte {
	// TODO: use encodeNode as base and look and SPEC.md
	panic("not implemented")
}

func encodeInUseRecord(record *EInUseRecord) *[]byte {
	var isHead byte = 0
	if record.IsHead {
		isHead = 1
	}
	buffer := []byte{1, byte(record.StoreType), isHead}

	buffer = append(buffer, encodeInt(record.ObjID)...)
	buffer = append(buffer, encodeInt(record.NextRecordID)...)

	return &buffer
}
