package engine

import (
	"testing"

	io "github.com/alikhil/TBMS/internals/io"
)

func TestGetLabelID(t *testing.T) {
	rw := io.LocalIO{}

	filename := "labelsStrings"

	label := []byte("fakeLabel")

	// label = append(label, make([]byte, 21-len(label))...)

	rw.CreateFile(filename)
	defer rw.DeleteFile(filename)

	rw.WriteBytes(filename, 0, label)

	re := RealEngine{}
	id, _ := re.GetLabelID("fakeLabel")

	if id != 0 {
		t.Errorf("Expected %v but got %v", 0, id)
	}

}
