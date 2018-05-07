package engine

import (
	io "github.com/alikhil/TBMS/internals/io"
	"testing"
)

func wrap(fnc func(*testing.T, io.IO), i io.IO) func(*testing.T) {
	return func(t *testing.T) {
		fnc(t, i)
	}
}

func testGetLabelID(t *testing.T, rw io.IO) {
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

func testGetSaveNode(t *testing.T, IO io.IO) {
	var en = RealEngine{IO}
	var node = ENode{
		ID:             FirstID,
		NextLabelID:    -1,
		NextPropertyID: -1,
		NextRelID:      -1}

	en.SaveObject(&node)
	defer en.DropDatabase()

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

func testInitDatabase(t *testing.T, IO io.IO) {
	var en = RealEngine{IO}

	en.InitDatabase()
	defer en.DropDatabase()

	// print32AllRecords(&en)

	for i := 2; i <= 9; i++ {
		id, ok := en.GetAndLockFreeIDForStore(EStore(i))
		if !ok && id != FirstID {
			t.Fatalf("Expected id = %d but get %d", FirstID, id)
		}
	}
}

func testGetAndLockFreeID(t *testing.T, IO io.IO) {
	var en = RealEngine{IO}

	en.InitDatabase()
	defer en.DropDatabase()

	id, ok := en.GetAndLockFreeIDForStore(StoreNode)
	if !ok {
		t.Fatalf("can not lock id for %s", FilenameStore[StoreNode])
	}
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

func testFindAndCreate(t *testing.T, IO io.IO) {
	var en = RealEngine{IO}

	en.InitDatabase()
	defer en.DropDatabase()

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

func testCreateAndLoadString(t *testing.T, IO io.IO) {
	var en = RealEngine{IO}

	en.InitDatabase()
	defer en.DropDatabase()

	var stringsForTest = []string{"shortString", "veryVeryLOOOOOOOOOooooooooooooooooooooooooooooooooooooooooooooooooooooooooooooveryVeryLOOOOOOOOOoooooooooooooooooooooooooooooooooooooooooooooooooooooooooooongString"}
	for _, str := range stringsForTest {
		es := en.CreateStringAndReturnFirstChunk(str)

		laoded := es.LoadString(&en)
		if laoded != str {
			t.Errorf("expected str '%s' but get '%s'", str, laoded)
		}
	}
}

func testFindObject(t *testing.T, IO io.IO) {

	var en = RealEngine{IO}

	en.InitDatabase()
	defer en.DropDatabase()

	var node = &ENode{
		ID:             FirstID,
		NextLabelID:    -1,
		NextPropertyID: 2,
		NextRelID:      -1}

	en.SaveObject(node)

	var res = &ENode{}
	found := en.FindObject(StoreNode, func(obj EObject) bool { return obj.(*ENode).NextPropertyID == 2 }, res)

	if !found {
		t.Errorf("can not found object when it should")
	} else if *res != *node {
		t.Errorf("Expected %v but get %v", node, res)
	}

}

func testPack(t *testing.T, getIO func() io.IO) {
	t.Run("TestFindObject", wrap(testFindObject, getIO()))
	t.Run("TestCreateAndLoadString", wrap(testCreateAndLoadString, getIO()))
	t.Run("TestFindAndCreate", wrap(testFindAndCreate, getIO()))
	t.Run("TestGetAndLockFreeID", wrap(testGetAndLockFreeID, getIO()))
	t.Run("TestInitDatabase", wrap(testInitDatabase, getIO()))
	t.Run("TestGetSaveNode", wrap(testGetSaveNode, getIO()))
	t.Run("TestGetLabelID", wrap(testGetLabelID, getIO()))
}

func TestFindObjectWithCache(t *testing.T) {
	testFindObject(t, getCache())
}
func TestCreateAndLoadStringWithCache(t *testing.T) {
	testCreateAndLoadString(t, getCache())
}
func TestFindAndCreateWithCache(t *testing.T) {
	testFindAndCreate(t, getCache())
}
func TestGetAndLockFreeIDWithCache(t *testing.T) {
	testGetAndLockFreeID(t, getCache())
}
func TestInitDatabaseWithCache(t *testing.T) {
	testInitDatabase(t, getCache())
}
func TestGetSaveNodeWithCache(t *testing.T) {
	testGetSaveNode(t, getCache())
}
func TestGetLabelIDWithCache(t *testing.T) {
	testGetLabelID(t, getCache())
}

func TestEngineWithLocalIO(t *testing.T) {
	testPack(t, func() io.IO { return &io.LocalIO{} })
}

func getCache() io.IO {
	var mapa = GetFileToBytesMap()

	cache := io.LRUCache{}
	cache.Init(io.LocalIO{}, mapa, 5)
	return &cache
}

func TestEngineWithCache(t *testing.T) {

	testPack(t, getCache)
}
