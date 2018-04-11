package engine

import (
	"testing"

	io "github.com/alikhil/TBMS/internals/io"
)

func TestGetLabelID(t *testing.T) {
	rw := io.LocalIO{}
	filename := FNLabelsStrings
	labelStr := "testLabelStr"

	label := []byte{1} // In Use byte

	label = append(label, ([]byte(labelStr))...)
	label = append(label, make([]byte, BytesPerLabelString-len(label))...) // fill left part with zeros

	rw.CreateFile(filename)
	defer rw.DeleteFile(filename)

	rw.WriteBytes(filename, 0, &label)

	re := RealEngine{IO: rw}
	id, _ := re.GetLabelID(labelStr)

	if id != 0 {
		t.Errorf("Expected %v but got %v", 0, id)
	}

}

func TestGetSaveNode(t *testing.T) {
	var en = RealEngine{IO: io.LocalIO{}}
	var node = ENode{
		ID:             0,
		NextLabelID:    -1,
		NextPropertyID: -1,
		NextRelID:      -1}

	en.SaveNode(&node)
	defer en.IO.DeleteFile(FNNodes)

	var parsedNode, ok = en.GetNodeByID(0)
	if !ok {
		t.Fatalf("Can not read saved node!")
	}
	if *parsedNode != node {
		t.Fatalf("Expected %+v but get %+v", node, parsedNode)
	}
}
