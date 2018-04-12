package engine

import (
	io "github.com/alikhil/TBMS/internals/io"
	"testing"
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

// for debuge purposes
func getAllRecords(re *RealEngine) []EInUseRecord {
	var list []EInUseRecord
	var next = re.GetInUseRecordIterator()
	for el, ok := next(); ok; el, ok = next() {
		list = append(list, *el)
	}
	return list
}

func TestEncodeParseInt(t *testing.T) {

	for a := -(1 << 21); a < (1 << 22); a++ {

		var real = parseInt(encodeInt(a))
		if real != a {
			t.Errorf("expected %v but recieved %v", a, real)
		}
	}
}

func TestInitDatabase(t *testing.T) {
	var en = RealEngine{IO: io.LocalIO{}}

	en.InitDatabase()
	defer en.DeleteFile(FNInUse)

	for i := 1; i < 9; i++ {
		id, ok := en.GetAndLockFreeIDForStore(EStore(i))
		if !ok && id != 0 {
			t.Fatalf("Expected id = 0 but get %d", id)
		}
	}
}

func TestGetAndLockFreeID(t *testing.T) {
	var en = RealEngine{IO: io.LocalIO{}}

	en.InitDatabase()
	defer en.DeleteFile(FNInUse)

	id, ok := en.GetAndLockFreeIDForStore(StoreNode)
	if !ok {
		t.Fatalf("can not lock id for %s", FilenameStore[StoreNode])
	}

	deleted := en.DeleteObject(id, StoreNode)
	if !deleted {
		t.Fatalf("object is not deleted")
	}

	newID, okNew := en.GetAndLockFreeIDForStore(StoreNode)

	if !okNew {
		t.Fatalf("can not lock id for %s after deletion", FilenameStore[StoreNode])
	}

	if id != newID {
		t.Fatalf("expected %v but get %v", id, newID)
	}
}
