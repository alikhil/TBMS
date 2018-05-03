package api

import (
	en "github.com/alikhil/TBMS/internals/engine"
	io "github.com/alikhil/TBMS/internals/io"
	"github.com/kmanley/golang-tuple"
	"testing"
)

func TestCreateRelationship(t *testing.T) {
	var re = &en.RealEngine{IO: io.LocalIO{}}

	re.InitDatabase()
	Init(re)
	defer re.DeleteFile(en.FNInUse)
	defer re.DeleteFile(en.FNNodes)
	defer re.DeleteFile(en.FNRelationships)
	defer re.DeleteFile(en.FNRelationshipTypes)

	nodeA := &en.ENode{
		ID:             en.FirstID,
		NextLabelID:    -1,
		NextPropertyID: -1,
		NextRelID:      -1}

	nodeB := &en.ENode{
		ID:             en.FirstID + 1,
		NextLabelID:    -1,
		NextPropertyID: -1,
		NextRelID:      -1}

	if !re.SaveObject(nodeA) || !re.SaveObject(nodeB) {
		t.Errorf("failed to save")
	}

	rel, ok := CreateRelationship(&Node{nodeA}, &Node{nodeB}, "likes")
	if !ok {
		t.Errorf("failed to create relationship")
	}

	if rel.ID != en.FirstID {
		t.Errorf("created invalid relationship")
	}

	// check if ids are same
	relType := &en.ERelationshipType{}
	re.FindObject(en.StoreRelationshipType, func(obj en.EObject) bool {
		return obj.(*en.ERelationshipType).TypeString == "likes"
	}, relType)

	if rel.TypeID != relType.ID {
		t.Errorf("expected %d but get %d", relType.ID, rel.TypeID)
	}

	// check if we can find our newly created relationship
	fillNext := re.GetEObjectIterator(en.StoreRelationship)
	fRel := &en.ERelationship{}
	foundRel := false
	for ok := fillNext(fRel); ok; ok = fillNext(fRel) {
		if fRel.TypeID == relType.ID {
			foundRel = true
		}
	}

	if !foundRel {
		t.Errorf("can not found newly created realationship")
	}

}

func TestCreateRelationshipWithProperties(t *testing.T) {
	var re = &en.RealEngine{IO: io.LocalIO{}}

	re.InitDatabase()
	Init(re)
	defer re.DeleteFile(en.FNInUse)
	defer re.DeleteFile(en.FNNodes)
	defer re.DeleteFile(en.FNRelationships)
	defer re.DeleteFile(en.FNRelationshipTypes)

	nodeA := &en.ENode{
		ID:             en.FirstID,
		NextLabelID:    -1,
		NextPropertyID: -1,
		NextRelID:      -1}

	nodeB := &en.ENode{
		ID:             en.FirstID + 1,
		NextLabelID:    -1,
		NextPropertyID: -1,
		NextRelID:      -1}

	if !re.SaveObject(nodeA) || !re.SaveObject(nodeB) {
		t.Errorf("failed to save")
	}

	rel, ok := CreateRelationship(&Node{nodeA}, &Node{nodeB}, "knows",
		tuple.NewTupleFromItems("since", 2015),
		tuple.NewTupleFromItems("met in", "kazan"),
		tuple.NewTupleFromItems("met_after", true))
	if !ok {
		t.Errorf("failed to create relationship")
	}

	if rel.ID != en.FirstID {
		t.Errorf("created invalid relationship")
	}

}
