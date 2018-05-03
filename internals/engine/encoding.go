package engine

import (
	"bytes"
	"encoding/binary"
	"github.com/alikhil/TBMS/internals/logger"
)

func encodeInt(val int32) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, ConventionByteOrder, val)
	if err != nil {
		logger.Error.Printf("Error on encoding int: %v", err)
	}
	return buf.Bytes()
}

func encode(val interface{}) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, ConventionByteOrder, val)
	if err != nil {
		logger.Error.Printf("Error on encoding interface{}: %v", err)
	}
	return buf.Bytes()
}

func boolToByte(b bool) byte {
	if b {
		return 1
	}
	return 0
}

func (node *ENode) encode() *[]byte {
	buffer := []byte{1} // Encode InUse

	buffer = append(buffer, encodeInt(node.NextLabelID)...)    // Add NextLabelID
	buffer = append(buffer, encodeInt(node.NextPropertyID)...) // Add NextProprertyID
	buffer = append(buffer, encodeInt(node.NextRelID)...)      // Add NextRelID

	return &buffer
}

func (rel *ERelationship) encode() *[]byte {
	buffer := []byte{1, boolToByte(rel.FirstInChain)} // Encode InUse

	buffer = append(buffer, encodeInt(rel.SecondNodeID)...)
	buffer = append(buffer, encodeInt(rel.FirstNodeID)...)
	buffer = append(buffer, encodeInt(rel.FirstNodeNxtRelID)...)
	buffer = append(buffer, encodeInt(rel.SecondNodeNxtRelID)...)
	buffer = append(buffer, encodeInt(rel.FirstNodePrvRelID)...)
	buffer = append(buffer, encodeInt(rel.SecondNodePrvRelID)...)
	buffer = append(buffer, encodeInt(rel.NextPropertyID)...)
	buffer = append(buffer, encodeInt(rel.TypeID)...)
	return &buffer
}

func (rel *ELabelString) encode() *[]byte {

	buffer := []byte{1} // In Use byte

	buffer = append(buffer, ([]byte(rel.String))...)
	buffer = append(buffer, make([]byte, BytesPerLabelString-len(buffer))...) // fill remaining part with zeros

	if len(buffer) > BytesPerLabelString {
		logger.Error.Fatalf("Label String length is too big - %s", rel.String)
	}
	return &buffer
}

func (rel *ERelationshipType) encode() *[]byte {
	buffer := []byte{1} // In Use byte

	buffer = append(buffer, ([]byte(rel.TypeString))...)
	buffer = append(buffer, make([]byte, BytesPerRelType-len(buffer))...) // fill remaining part with zeros

	if len(buffer) > BytesPerRelType {
		logger.Error.Fatalf("Label String length is too big - %s", rel.TypeString)
	}
	return &buffer
}

func (rel *EPropertyKey) encode() *[]byte {
	buffer := []byte{1} // In User byte
	buffer = append(buffer, ([]byte(rel.KeyString))...)
	buffer = append(buffer, make([]byte, BytesPerPropertyKey-len(buffer))...) // fill remaining part with zeros

	if len(buffer) > BytesPerPropertyKey {
		logger.Error.Fatalf("Property key length is too big - %s", rel.KeyString)

	}
	return &buffer
}

func (prop *EProperty) encode() *[]byte {
	buffer := []byte{1, byte(prop.Typename)} //In Use

	buffer = append(buffer, encodeInt(prop.KeyStringID)...)
	buffer = append(buffer, encodeInt(prop.NextPropertyID)...)
	buffer = append(buffer, encode(prop.ValueOrStringPtr)...)

	return &buffer
}

func (str *EString) encode() *[]byte {

	buffer := []byte{1, str.Extra}

	buffer = append(buffer, encodeInt(str.NextPartID)...)
	buffer = append(buffer, *str.Value...)

	return &buffer
}

func (record *EInUseRecord) encode() *[]byte {

	buffer := []byte{1, byte(record.StoreType), boolToByte(record.IsHead)}

	buffer = append(buffer, encodeInt(record.ObjID)...)
	buffer = append(buffer, encodeInt(record.NextRecordID)...)

	return &buffer
}
