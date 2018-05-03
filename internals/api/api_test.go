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
	var ival int32 = 2015
	var fval float32 = 3.3
	rel, ok := CreateRelationship(&Node{nodeA}, &Node{nodeB}, "knows",
		tuple.NewTupleFromItems("since", ival),
		tuple.NewTupleFromItems("met in", "kazan"),
		tuple.NewTupleFromItems("met_after", true),
		tuple.NewTupleFromItems("double_val", fval))
	if !ok {
		t.Errorf("failed to create relationship")
	}

	if rel.ID != en.FirstID {
		t.Errorf("created invalid relationship")
	}

	if rel.NextPropertyID == -1 {
		t.Errorf("No properties found!")
	}

	propsRef, foundProps := getProperties(rel.NextPropertyID)
	if !foundProps {
		t.Errorf("Can not load properties!")
	}

	props := *propsRef

	if props["since"] != ival {
		t.Errorf("Expected %v but got %v", ival, props["since"])
	}

	if props["met in"] != "kazan" {
		t.Errorf("Expected %v but got %v", "kazan", props["met in"])
	}

	if props["met_after"] != true {
		t.Errorf("Expected %v but got %v", true, props["met_after"])
	}

	if props["double_val"] != fval {
		t.Errorf("Expected %v but got %v", fval, props["double_val"])

	}

}
