package engine

import (
	"github.com/alikhil/TBMS/internals/logger"
)

// TODO: pass all arrays and slices by reference

/* ************** *
	 ITERATORS
 * ************** */

func (re *RealEngine) GetLabelIteratorFromId(labelID int32) func() (int32, bool) {
	// TODO: implement it using GetLabelIterator()
	return func() (int32, bool) {
		return 0, false
	}
}

func (re *RealEngine) GetNodesByLabelIterator(label string) func() (*ENode, bool) {
	fillNextNode := re.GetEObjectIterator(StoreNode)
	neededlabelID, ok := re.GetLabelID(label)
	if !ok {

		return func() (*ENode, bool) {
			return nil, false
		}
	}

	return func() (*ENode, bool) {
		node := &ENode{}
		ok := fillNextNode(node)
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

func (re *RealEngine) getObjectIterator(store EStore) func() ([]byte, bool) {
	var curOffset int32
	return func() (data []byte, ok bool) {
		data, ok = re.IO.ReadBytes(FilenameStore[store], curOffset, BytesPerStore[store])
		if ok {
			curOffset += BytesPerStore[store]
		}
		return
	}
}

// GetEObjectIterator - returns function for iterating objects of certain store
// Example:
//	nextFill := re.GetEObjectIterator(StoreNodes)
//  node := &ENode{} // empty obj to where next node will be loaded
//  for ok := nextFill(node); ok; ok = nextFill(node) {
//  	// use node somehow
//  }
func (re *RealEngine) GetEObjectIterator(store EStore) func(EObject) bool {
	next := re.getObjectIterator(store)
	var i int32
	return func(ob EObject) bool {
	iterate:
		data, ok := next()
		if ok {
			i++
			if notInUse(&data) {
				goto iterate
			}

			ob.fill(&data, i)
			return ok
		}
		return false
	}
}

func (re *RealEngine) GetNodeRelationshipsIteratorStartingFrom(nodeID int32, nxtRelID int32) func() (*ERelationship, bool) {
	return func() (*ERelationship, bool) {

		if nxtRelID == -1 {
			return nil, false
		}
		cur := &ERelationship{ID: nxtRelID}
		ok := re.GetObject(cur)
		if !ok {
			return nil, false
		}
		nxtRelID = cur.GetPart(nodeID).NodeNxtRelID
		return cur, true
	}
}

// GetNodeRelationshipsIterator return iterator that allows to iterate all rellationship of node
func (re *RealEngine) GetNodeRelationshipsIterator(nodeID int32) func() (*ERelationship, bool) {
	node := &ENode{ID: nodeID}
	ok := re.GetObject(node)
	if !ok {
		return func() (*ERelationship, bool) { return nil, false }
	}

	nxtID := node.NextRelID

	return re.GetNodeRelationshipsIteratorStartingFrom(nodeID, nxtID)
}

/* ***************** *
	Getters/Setters
 * ***************** */

func (re *RealEngine) GetLabelID(label string) (int32, bool) {
	fillNextLabelStr := re.GetEObjectIterator(StoreLabelString)
	l := &ELabelString{}
	for ok := fillNextLabelStr(l); ok; ok = fillNextLabelStr(l) {
		if ok && label == l.String {
			return l.ID, true
		}
	}
	return -1, false
}

func (re *RealEngine) getObjectByID(store EStore, id int32) (*[]byte, bool) {
	offset := BytesPerStore[store] * (id - FirstID)
	data, ok := re.IO.ReadBytes(FilenameStore[store], offset, BytesPerStore[store])
	if !ok {
		logger.Trace.Printf("Object with id = %d cannot be read from file %s", id, FilenameStore[store])
	}
	return &data, ok
}

func notInUse(data *[]byte) bool {
	return !parseBool((*data)[0])
}

// GetObject - loads any object from database by id to passed object
// Example:
//	node := &ENode{ID: 12}
// 	found := re.GetObject(node)
func (re *RealEngine) GetObject(obj EObject) bool {
	if obj.getID() < FirstID {
		logger.Error.Fatalf("Object ID is not set for element of %s ", FilenameStore[obj.getStore()])
	}
	data, ok := re.getObjectByID(obj.getStore(), obj.getID())
	if notInUse(data) {
		return false
	}
	obj.fill(data, obj.getID())
	return ok
}

// SaveObject - saves any EObject to file
func (re *RealEngine) SaveObject(obj EObject) bool {
	return re.saveObject(obj.getStore(), obj.getID(), obj.encode())
}

func (re *RealEngine) saveObject(store EStore, id int32, data *[]byte) bool {
	offset := BytesPerStore[store] * (id - FirstID)
	ok := re.IO.WriteBytes(FilenameStore[store], offset, data)
	if !ok {
		logger.Warning.Printf("Failed to save object with id = %d to file %s", id, FilenameStore[store])
	}
	return ok
}

func (re *RealEngine) DeleteObject(obj EObject) bool {
	store := obj.getStore()
	objID := obj.getID()
	emptyRecord := make([]byte, BytesPerStore[store])
	saved := re.saveObject(store, objID, &emptyRecord)
	if !saved {
		logger.Error.Printf("Failed to delete object with id = %v in store %s", objID, FilenameStore[store])
		return false
	}

	headRecord, found := re.findHeadInUseRecord(store)
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
	return re.SaveObject(newRecord) && re.SaveObject(headRecord)
}

/* ******************** *
	Obtaining free ids
 * ******************** */

func (re *RealEngine) findHeadInUseRecord(store EStore) (*EInUseRecord, bool) {

	// WARN: there could be problems with concurrency
	var fillNext = re.GetEObjectIterator(StoreInUse)

	var record = &EInUseRecord{}
	for ok := fillNext(record); ok; ok = fillNext(record) {
		if record.StoreType == store && record.IsHead {
			return record, true
		}
	}
	logger.Error.Printf("Can not find head record for store %s", FilenameStore[store])
	return nil, false

}

// GetAndLockFreeIDForStore - obtains free id for store
func (re *RealEngine) GetAndLockFreeIDForStore(store EStore) (int32, bool) {
	record, found := re.findHeadInUseRecord(store)
	if !found {
		return -1, false
	}

	if record.NextRecordID != -1 {
		// Use deleted obj places if they exist
		nxtRecord := &EInUseRecord{ID: record.NextRecordID}
		ok := re.GetObject(nxtRecord)
		if !ok {
			logger.Error.Printf("Can not get in use record with id = %v", nxtRecord.ID)
			return -1, false
		}
		record.NextRecordID = nxtRecord.NextRecordID // Update link to next free id
		if re.SaveObject(record) && re.DeleteObject(nxtRecord) {
			return nxtRecord.ObjID, true
		}
		logger.Warning.Printf("Failed to update record(%d) and delete nxtRecord(%d)", record.ID, record.NextRecordID)
		return -1, false

	} else {
		// TODO: write code here
		newID := record.ObjID
		record.ObjID++

		return newID, re.SaveObject(record)
	}
}

func (re *RealEngine) setupInUseFor(store EStore) {
	var inUseRecord = &EInUseRecord{
		ID:           int32(store),
		StoreType:    store,
		IsHead:       true,
		ObjID:        FirstID,
		NextRecordID: -1,
	}
	ok := re.SaveObject(inUseRecord)
	if !ok {
		logger.Error.Fatalf("Can not init InUseStore for %s", FilenameStore[store])
	}
}

func print32AllRecords(re *RealEngine) {
	var list []EInUseRecord
	var nextFill = re.GetEObjectIterator(StoreInUse)
	el := &EInUseRecord{}
	for ok := nextFill(el); ok; ok = nextFill(el) {
		list = append(list, *el)
	}
	logger.Trace.Printf("now: %+v", list)
}

/* ******************** *
		Advanced
 * ******************** */

// FindOrCreateObject ...
func (re *RealEngine) FindOrCreateObject(store EStore, pred func(EObject) bool, create func(int32) EObject) (int32, bool) {
	// fillNextObj := re.GetEObjectIterator(store)
	curObj := create(-1)

	if re.FindObject(store, pred, curObj) {
		return curObj.getID(), true
	}

	return re.CreateObject(store, create)
}

// CreateObject ...
func (re *RealEngine) CreateObject(store EStore, create func(int32) EObject) (int32, bool) {

	newID, ok := re.GetAndLockFreeIDForStore(store)
	if !ok {
		logger.Error.Printf("can not get free id for %s", FilenameStore[store])
		return -1, false
	}

	nObj := create(newID)
	saved := re.SaveObject(nObj)

	if !saved {
		logger.Error.Printf("failed to save new obj wih id %v to %s", newID, FilenameStore[store])
		return -1, false
	}
	return newID, true
}

func (re *RealEngine) FindObject(store EStore, pred func(EObject) bool, res EObject) bool {
	fillNextObj := re.GetEObjectIterator(store)

	for ok := fillNextObj(res); ok; ok = fillNextObj(res) {
		if pred(res) {
			return true
		}
	}
	return false
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
		ID:           int32(StoreInUse),
		StoreType:    StoreInUse,
		IsHead:       true,
		ObjID:        10,
		NextRecordID: -1,
	}
	ok := re.SaveObject(inUseRecord)
	if !ok {
		logger.Error.Fatalf("Can not init InUseStore for %s", FilenameStore[StoreInUse])
	}
}

// DropDatabase - removes all database files
func (re *RealEngine) DropDatabase() {
	re.DeleteFile(FNInUse)
	re.DeleteFile(FNNodes)
	re.DeleteFile(FNRelationships)
	re.DeleteFile(FNRelationshipTypes)

	re.DeleteFile(FNLabels)
	re.DeleteFile(FNLabelsStrings)
	re.DeleteFile(FNStrings)
	re.DeleteFile(FNProperties)
	re.DeleteFile(FNPropertyKeys)
}
