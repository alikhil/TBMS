package engine

import (
	"bytes"
	"encoding/binary"
)

func encodeInt(val int) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, ConventionByteOrder, val)
	return buf.Bytes()
}

func encodeNode(node *ENode) *[]byte {
	buffer := []byte{1} // Encode InUse

	buffer = append(buffer, encodeInt(node.NextLabelID)...)    // Add NextLabelID
	buffer = append(buffer, encodeInt(node.NextPropertyID)...) // Add NextProprertyID
	buffer = append(buffer, encodeInt(node.NextRelID)...)      // Add NextRelID

	return &buffer
}
