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

	if id != FirstID {
		t.Errorf("Expected %v but got %v", 0, id)
	}

}

func TestGetSaveNode(t *testing.T) {
	var en = RealEngine{IO: io.LocalIO{}}
	var node = ENode{
		ID:             FirstID,
		NextLabelID:    -1,
		NextPropertyID: -1,
		NextRelID:      -1}

	en.SaveObject(&node)
	defer en.IO.DeleteFile(FNNodes)

	var parsedNode = &ENode{ID: FirstID}
	ok := en.GetObject(parsedNode)
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
	var nextFill = re.GetEObjectIterator(StoreInUse)
	el := &EInUseRecord{}
	for ok := nextFill(el); ok; ok = nextFill(el) {
		list = append(list, *el)
	}
	return list
}

func TestEncodeParseInt(t *testing.T) {

	var a int32
	for a = -(1 << 21); a < (1 << 22); a++ {

		var real = parseInt(encodeInt(a))
		if real != a {
			t.Errorf("expected %v but recieved %v", a, real)
		}
	}
}

func TestInitDatabase(t *testing.T) {
	var en = RealEngine{IO: io.LocalIO{}}

	en.InitDatabase()
	print32AllRecords(&en)
	defer en.DeleteFile(FNInUse)

	for i := 2; i <= 9; i++ {
		id, ok := en.GetAndLockFreeIDForStore(EStore(i))
		if !ok && id != FirstID {
			t.Fatalf("Expected id = %d but get %d", FirstID, id)
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

	defer func() {
		if en.FileExists(FNNodes) {
			en.DeleteFile(FNNodes)
		}
	}()
	deleted := en.DeleteObject(&ENode{ID: id})
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

func TestFindAndCreate(t *testing.T) {
	var en = RealEngine{IO: io.LocalIO{}}

	en.InitDatabase()
	defer en.DeleteFile(FNInUse)
	defer en.DeleteFile(FNLabelsStrings)

	var label = "mylabel"

	id, ok := en.FindOrCreateObject(StoreLabelString,
		func(ob EObject) bool { return ob.(*ELabelString).String == label },
		func(id int32) EObject { return &ELabelString{ID: id, String: label} })

	if !ok {
		t.Fatalf("Failed to create object")
	}

	foundID, found := en.FindOrCreateObject(StoreLabelString,
		func(ob EObject) bool {
			return ob.(*ELabelString).String == label
		},
		func(id int32) EObject { return &ELabelString{ID: id, String: label} })

	if !found {
		t.Fatalf("Failed to find object")
	}

	if id != foundID {
		t.Fatalf("expected to get id = %v but get %v", id, foundID)
	}

}
