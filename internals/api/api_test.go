package api

import (
	"log"
	"testing"

	"github.com/alikhil/TBMS/internals/logger"

	en "github.com/alikhil/TBMS/internals/engine"
	io "github.com/alikhil/TBMS/internals/io"
	"github.com/kmanley/golang-tuple"
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
	defer re.DeleteFile(en.FNProperties)
	defer re.DeleteFile(en.FNPropertyKeys)
	// TODO: check bug
	defer re.DeleteFile(en.FNStrings)

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

	fnext := engine.GetEObjectIterator(en.StoreProperty)
	prop := &en.EProperty{}
	for ok := fnext(prop); ok; ok = fnext(prop) {
		log.Printf("%v", prop)
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
func TestGetLabels(t *testing.T) {
	var re = &en.RealEngine{IO: io.LocalIO{}}
	re.InitDatabase()
	Init(re)
	defer re.DeleteFile(en.FNInUse)
	defer re.DeleteFile(en.FNNodes)
	defer re.DeleteFile(en.FNRelationships)
	defer re.DeleteFile(en.FNRelationshipTypes)
	defer re.DeleteFile(en.FNLabels)
	defer re.DeleteFile(en.FNLabelsStrings)
	defer re.DeleteFile(en.FNProperties)
	defer re.DeleteFile(en.FNPropertyKeys)

	node, _ := CreateNode("test_label")

	labels := *node.GetLabels()
	if len(labels) > 1 || len(labels) < 1 {
		t.Errorf("Wrong number of labels")
	}

	if labels[0] != "test_label" {
		t.Errorf("wrong label")
	}

}

func TestGetNodeProperty(t *testing.T) {
	var re = &en.RealEngine{IO: io.LocalIO{}}
	re.InitDatabase()
	Init(re)
	defer re.DeleteFile(en.FNInUse)
	defer re.DeleteFile(en.FNNodes)
	defer re.DeleteFile(en.FNRelationships)
	defer re.DeleteFile(en.FNRelationshipTypes)
	defer re.DeleteFile(en.FNLabels)
	defer re.DeleteFile(en.FNLabelsStrings)
	defer re.DeleteFile(en.FNProperties)
	defer re.DeleteFile(en.FNStrings)
	defer re.DeleteFile(en.FNPropertyKeys)

	node, _ := CreateNode("test_label",
		tuple.NewTupleFromItems("test_prop1", 1),
		tuple.NewTupleFromItems("test_prop2", "two"))

	pr, ok := node.GetProperties()
	if !ok {
		t.Errorf("Can't get properties")
	}
	properties := *pr

	if !ok {
		t.Errorf("can't get properties")
	}

	if properties["test_prop1"] != int32(1) {
		t.Errorf("wrong 1st property: expected 1, get %v", properties["test_prop1"])
	}
	if properties["test_prop2"] != "two" {
		t.Errorf("wrong 2nd property: expected 'two', get %v", properties["test_prop2"])
	}

}
func TestGetFromRelationship(t *testing.T) {
	var re = &en.RealEngine{IO: io.LocalIO{}}
	re.InitDatabase()
	Init(re)
	defer re.DeleteFile(en.FNInUse)
	defer re.DeleteFile(en.FNNodes)
	defer re.DeleteFile(en.FNRelationships)
	defer re.DeleteFile(en.FNRelationshipTypes)
	defer re.DeleteFile(en.FNLabels)
	defer re.DeleteFile(en.FNLabelsStrings)
	defer re.DeleteFile(en.FNProperties)
	defer re.DeleteFile(en.FNPropertyKeys)

	node1, _ := CreateNode("from", tuple.NewTupleFromItems("id", 1))
	node2, _ := CreateNode("to", tuple.NewTupleFromItems("id", 2))

	rel, _ := CreateRelationship(node1, node2, "link")

	retrived_node := *rel.GetFrom()

	old_id, _ := node1.GetProperty("id")
	new_id, _ := retrived_node.GetProperty("id")

	if old_id != new_id {
		t.Errorf("expected to be equal, but get %v and %v", old_id, new_id)
	}
}

func TestGetToRelationship(t *testing.T) {
	var re = &en.RealEngine{IO: io.LocalIO{}}
	re.InitDatabase()
	Init(re)
	defer re.DeleteFile(en.FNInUse)
	defer re.DeleteFile(en.FNNodes)
	defer re.DeleteFile(en.FNRelationships)
	defer re.DeleteFile(en.FNRelationshipTypes)
	defer re.DeleteFile(en.FNLabels)
	defer re.DeleteFile(en.FNLabelsStrings)
	defer re.DeleteFile(en.FNProperties)
	defer re.DeleteFile(en.FNPropertyKeys)

	node1, _ := CreateNode("from", tuple.NewTupleFromItems("id", 1))
	node2, _ := CreateNode("to", tuple.NewTupleFromItems("id", 2))

	rel, _ := CreateRelationship(node1, node2, "link")

	retrived_node := *rel.GetTo()

	old_id, _ := node2.GetProperty("id")
	new_id, _ := retrived_node.GetProperty("id")

	if old_id != new_id {
		t.Errorf("expected to be equal, but get %v and %v", old_id, new_id)
	}
}

func TestGetRelationshipType(t *testing.T) {
	var re = &en.RealEngine{IO: io.LocalIO{}}
	re.InitDatabase()
	Init(re)
	defer re.DeleteFile(en.FNInUse)
	defer re.DeleteFile(en.FNNodes)
	defer re.DeleteFile(en.FNRelationships)
	defer re.DeleteFile(en.FNRelationshipTypes)
	defer re.DeleteFile(en.FNLabels)
	defer re.DeleteFile(en.FNLabelsStrings)
	defer re.DeleteFile(en.FNProperties)
	defer re.DeleteFile(en.FNPropertyKeys)

	node1, _ := CreateNode("from", tuple.NewTupleFromItems("id", 1))
	node2, _ := CreateNode("to", tuple.NewTupleFromItems("id", 2))

	rel, _ := CreateRelationship(node1, node2, "link")

	relType := rel.GetType()

	if relType != "link" {
		t.Errorf("Wrong relationship type")
	}
}

func TestGetRelationshipProperty(t *testing.T) {
	var re = &en.RealEngine{IO: io.LocalIO{}}
	re.InitDatabase()
	Init(re)
	defer re.DeleteFile(en.FNInUse)
	defer re.DeleteFile(en.FNNodes)
	defer re.DeleteFile(en.FNRelationships)
	defer re.DeleteFile(en.FNRelationshipTypes)
	defer re.DeleteFile(en.FNLabels)
	defer re.DeleteFile(en.FNLabelsStrings)
	defer re.DeleteFile(en.FNProperties)
	defer re.DeleteFile(en.FNPropertyKeys)
	defer re.DeleteFile(en.FNStrings)

	node1, _ := CreateNode("from", tuple.NewTupleFromItems("id", 1))
	node2, _ := CreateNode("to", tuple.NewTupleFromItems("id", 2))

	rel, _ := CreateRelationship(node1, node2, "link", tuple.NewTupleFromItems("test-prop", "hello"))

	pr, ok := rel.GetProperties()
	if !ok {
		t.Errorf("Can't get properties")
	}
	properties := *pr

	if !ok {
		t.Errorf("can't get properties")
	}

	if properties["test-prop"] != "hello" {
		t.Errorf("wrong 1st property: expected 'hello' , get %v", properties["test-prop"])
	}
}

func TestSelectNodesWhere(t *testing.T) {
	var re = &en.RealEngine{IO: io.LocalIO{}}
	re.InitDatabase()
	Init(re)
	defer re.DeleteFile(en.FNInUse)
	defer re.DeleteFile(en.FNNodes)
	defer re.DeleteFile(en.FNRelationships)
	defer re.DeleteFile(en.FNRelationshipTypes)
	defer re.DeleteFile(en.FNLabels)
	defer re.DeleteFile(en.FNLabelsStrings)
	defer re.DeleteFile(en.FNProperties)
	defer re.DeleteFile(en.FNPropertyKeys)
	defer re.DeleteFile(en.FNStrings)

	CreateNode("show", tuple.NewTupleFromItems("id", 1))
	CreateNode("show", tuple.NewTupleFromItems("id", 2))
	CreateNode("dont_show", tuple.NewTupleFromItems("id", 3))
	CreateNode("dont_show", tuple.NewTupleFromItems("id", 4))

	nodes, _ := SelectNodesWhere(func(node *Node) bool {
		l := *node.GetLabels()
		label := l[0]
		return label == "show"
	})

	logger.Info.Printf("len(nodes) = %v", len(nodes))
	for _, elem := range nodes {
		l := *elem.GetLabels()
		p, _ := elem.GetProperty("id")

		label := l[0]
		if label != "show" {
			t.Errorf("Returns wrong node: %v with id %v", label, p)
		}
	}

}

func TestSelectRelationshipsWhere(t *testing.T) {
	var re = &en.RealEngine{IO: io.LocalIO{}}
	re.InitDatabase()
	Init(re)
	defer re.DeleteFile(en.FNInUse)
	defer re.DeleteFile(en.FNNodes)
	defer re.DeleteFile(en.FNRelationships)
	defer re.DeleteFile(en.FNRelationshipTypes)
	defer re.DeleteFile(en.FNLabels)
	defer re.DeleteFile(en.FNLabelsStrings)
	defer re.DeleteFile(en.FNProperties)
	defer re.DeleteFile(en.FNPropertyKeys)

	node1, _ := CreateNode("a", tuple.NewTupleFromItems("id", 1))
	node2, _ := CreateNode("a", tuple.NewTupleFromItems("id", 2))

	CreateRelationship(node1, node2, "show")
	CreateRelationship(node1, node2, "show")
	CreateRelationship(node1, node2, "dont_show")
	CreateRelationship(node1, node2, "dont_show")

	nodes, _ := SelectNodesWhere(func(node *Node) bool {
		l := *node.GetLabels()
		label := l[0]
		return label == "show"
	})

	logger.Info.Printf("len(nodes) = %v", len(nodes))
	for _, elem := range nodes {
		l := *elem.GetLabels()
		p, _ := elem.GetProperty("id")

		label := l[0]
		if label != "show" {
			t.Errorf("Returns wrong node: %v with id %v", label, p)
		}
	}

}

func TestGetRelationships(t *testing.T) {
	var re = &en.RealEngine{IO: io.LocalIO{}}
	re.InitDatabase()
	Init(re)
	defer re.DeleteFile(en.FNInUse)
	defer re.DeleteFile(en.FNNodes)
	defer re.DeleteFile(en.FNRelationships)
	defer re.DeleteFile(en.FNRelationshipTypes)
	defer re.DeleteFile(en.FNLabels)
	defer re.DeleteFile(en.FNLabelsStrings)
	defer re.DeleteFile(en.FNProperties)
	defer re.DeleteFile(en.FNPropertyKeys)

	node1, _ := CreateNode("a", tuple.NewTupleFromItems("id", 1))
	node2, _ := CreateNode("a", tuple.NewTupleFromItems("id", 2))

	CreateRelationship(node1, node2, "show")
	CreateRelationship(node1, node2, "show")
	CreateRelationship(node1, node2, "dont_show")
	CreateRelationship(node1, node2, "dont_show")

	relationships := *node1.GetRelationships()

	if len(relationships) != 4 {
		t.Errorf("Expected lenght: 4, get: %v", len(relationships))
	}

}
