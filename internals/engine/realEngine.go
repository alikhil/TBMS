package engine

import (
	"github.com/alikhil/TBMS/internals/logger"
)

// TODO: pass all arrays and slices by reference

/* ************** *
	 ITERATORS
 * ************** */

func (re *RealEngine) GetLabelIteratorFromId(labelID int) func() (int, bool) {
	// TODO: implement it using GetLabelIterator()
	return func() (int, bool) {
		return 0, false
	}
}

func (re *RealEngine) GetNodesByLabelIterator(label string) func() (*ENode, bool) {
	nextNode := re.GetNodesIterator()
	neededlabelID, ok := re.GetLabelID(label)
	if !ok {

		return func() (*ENode, bool) {
			return nil, false
		}
	}

	return func() (*ENode, bool) {
		node, ok := nextNode()
		if ok {
			nextLabel := re.GetLabelIteratorFromId(node.NextLabelID)
			for labelID, ok := nextLabel(); ok; {
				if labelID == neededlabelID {
					return node, true
				}
			}
			return nil, false
		}
		return nil, false
	}
}

func (re *RealEngine) GetNodesIterator() func() (*ENode, bool) {
	next := re.GetObjectIterator(FNNodes, BytesPerNode)
	i := 0
	return func() (*ENode, bool) {

	iterate:
		data, ok := next()
		if ok {
			node, nodeInUse := parseNode(&data, i)
			i++
			if !nodeInUse {
				// jump to another node, since this one is not in use
				goto iterate
			}
			return node, ok
		}
		return nil, false
	}
}

func (re *RealEngine) GetLabelStringIterator() func() (*ELabelString, bool) {
	next := re.GetObjectIterator(FNLabelsStrings, BytesPerLabelString)
	i := 0
	return func() (*ELabelString, bool) {

	iterate:
		data, ok := next()
		if ok {
			label, labelStringInUse := parseLabelString(&data, i)
			i++
			if !labelStringInUse {
				goto iterate
			}
			return label, ok
		}
		return nil, false
	}
}

func (re *RealEngine) GetRelationshiptIterator() func() (*ERelationship, bool) {
	next := re.GetObjectIterator(FNRelationships, BytesPerRelationship)
	i := 0
	return func() (*ERelationship, bool) {

	iterate:
		data, ok := next()
		if ok {
			rel, relInUse := parseRelationship(&data, i)
			i++
			if !relInUse {
				goto iterate
			}
			return rel, ok
		}
		return nil, false
	}
}

func (re *RealEngine) GetObjectIterator(filename string, recordLength int) func() ([]byte, bool) {
	curOffset := 0
	return func() (data []byte, ok bool) {
		data, ok = re.IO.ReadBytes(filename, curOffset, recordLength)
		if ok {
			curOffset += recordLength
		}
		return
	}
}

func (re *RealEngine) GetInUseRecordIterator() func() (*EInUseRecord, bool) {

	next := re.GetObjectIterator(FNInUse, BytesPerInUse)
	i := 0
	return func() (*EInUseRecord, bool) {

	iterate:
		data, ok := next()
		if ok {
			rec, recInUse := parseInUse(&data, i)
			i++
			if !recInUse {
				goto iterate
			}
			return rec, ok
		}
		return nil, false
	}
}

/* ***************** *
	Getters/Setters
 * ***************** */

func (re *RealEngine) GetLabelID(label string) (int, bool) {
	next := re.GetLabelStringIterator()
	i := 0
	for l, ok := next(); ok; l, ok = next() {
		if ok && label == l.String {
			return i, true
		}
		i++
	}
	return -1, false
}

// GetObjectByID returns byte record of any object from certain file
func (re *RealEngine) GetObjectByID(filename string, recordLength, id int) (*[]byte, bool) {
	offset := recordLength * id
	data, ok := re.IO.ReadBytes(filename, offset, recordLength)
	if !ok {
		logger.Trace.Printf("Object with id = %d cannot be read from file %s", id, filename)
	}
	return &data, ok
}

func (re *RealEngine) GetNodeByID(id int) (*ENode, bool) {
	data, ok := re.GetObjectByID(FNNodes, BytesPerNode, id)
	if !ok {
		return nil, false
	}
	return parseNode(data, id)
}

func (re *RealEngine) GetInUseRecord(id int) (*EInUseRecord, bool) {
	data, ok := re.GetObjectByID(FNInUse, BytesPerInUse, id)
	if !ok {
		return nil, false
	}
	return parseInUse(data, id)
}

func (re *RealEngine) saveObject(store EStore, id int, data *[]byte) bool {
	offset := BytesPerStore[store] * id
	ok := re.IO.WriteBytes(FilenameStore[store], offset, data)
	if !ok {
		logger.Warning.Printf("Failed to save object with id = %d to file %s", id, FilenameStore[store])
	}
	return ok
}

func (re *RealEngine) SaveNode(node *ENode) bool {
	data := encodeNode(node)
	return re.saveObject(StoreNode, node.ID, data)
}

func (re *RealEngine) SaveInUseRecord(record *EInUseRecord) bool {
	data := encodeInUseRecord(record)
	return re.saveObject(StoreInUse, record.ID, data)
}

func (re *RealEngine) DeleteObject(objID int, store EStore) bool {
	emptyRecord := make([]byte, BytesPerStore[store])
	saved := re.saveObject(store, objID, &emptyRecord)
	if !saved {
		logger.Error.Printf("Failed to delete object with id = %v in store %s", objID, FilenameStore[store])
		return false
	}

	headRecord, found := re.FindHeadInUseRecord(store)
	if !found {
		return false
	}

	// WARN: recursion problem??
	var newRecordID, ok = re.GetAndLockFreeIDForStore(StoreInUse)
	if !ok {
		logger.Warning.Printf("Can not update inUse table on object(%d - %s) deletion", objID, FilenameStore[store])
	}
	var newRecord = &EInUseRecord{
		ID:           newRecordID,
		StoreType:    StoreInUse,
		IsHead:       false,
		ObjID:        objID,
		NextRecordID: headRecord.NextRecordID,
	}

	headRecord.NextRecordID = newRecordID
	return re.SaveInUseRecord(newRecord) && re.SaveInUseRecord(headRecord)
}

/* ******************** *
	Obtaining free ids
 * ******************** */

func (re *RealEngine) FindHeadInUseRecord(store EStore) (*EInUseRecord, bool) {

	// WARN: there could be problems with concurrency
	var next = re.GetInUseRecordIterator()

	for record, ok := next(); ok; record, ok = next() {
		if record.StoreType == store && record.IsHead {
			return record, true
		}
	}
	logger.Error.Printf("Can not find head record for store %s", FilenameStore[store])
	return nil, false

}

func (re *RealEngine) GetAndLockFreeIDForStore(store EStore) (int, bool) {
	record, found := re.FindHeadInUseRecord(store)
	if !found {
		return -1, false
	}

	if record.NextRecordID != -1 {
		// Use deleted obj places if they exist
		nxtRecord, ok := re.GetInUseRecord(record.NextRecordID)
		if !ok {
			logger.Error.Printf("Can not get in use record with id = %v", nxtRecord.ID)
			return -1, false
		}
		record.NextRecordID = nxtRecord.NextRecordID // Update link to next free id
		if re.SaveInUseRecord(record) && re.DeleteObject(nxtRecord.ID, StoreInUse) {
			return nxtRecord.ObjID, true
		}
		logger.Warning.Printf("Failed to update record(%d) and delete nxtRecord(%d)", record.ID, record.NextRecordID)
		return -1, false

	} else {
		// TODO: write code here
		newID := record.ObjID
		record.ObjID++

		return newID, re.SaveInUseRecord(record)
	}
}

func (re *RealEngine) setupInUseFor(store EStore) {
	var inUseRecord = &EInUseRecord{
		ID:           int(store),
		StoreType:    store,
		IsHead:       true,
		ObjID:        0,
		NextRecordID: -1,
	}
	ok := re.SaveInUseRecord(inUseRecord)
	if !ok {
		logger.Error.Fatalf("Can not init InUseStore for %s", FilenameStore[store])
	}
}

func printAllRecords(re *RealEngine) {
	var list []EInUseRecord
	var next = re.GetInUseRecordIterator()
	for el, ok := next(); ok; el, ok = next() {
		list = append(list, *el)
	}
	logger.Trace.Printf("now: %+v", list)
}

// InitDatabase should be called to initialize data needed to start database first time
func (re *RealEngine) InitDatabase() {
	// Setup InUse Store
	re.setupInUseFor(StoreNode)
	re.setupInUseFor(StoreProperty)
	re.setupInUseFor(StoreRelationship)
	re.setupInUseFor(StoreLabel)
	re.setupInUseFor(StoreRelationshipType)
	re.setupInUseFor(StoreLabelString)
	re.setupInUseFor(StorePropertyKey)
	re.setupInUseFor(StoreString)

	var inUseRecord = &EInUseRecord{
		ID:           int(StoreInUse),
		StoreType:    StoreInUse,
		IsHead:       true,
		ObjID:        9,
		NextRecordID: -1,
	}
	ok := re.SaveInUseRecord(inUseRecord)
	if !ok {
		logger.Error.Fatalf("Can not init InUseStore for %s", FilenameStore[StoreInUse])
	}
}