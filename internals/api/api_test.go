package api

import (
	"log"
	"testing"

	"github.com/alikhil/TBMS/internals/logger"

	en "github.com/alikhil/TBMS/internals/engine"
	io "github.com/alikhil/TBMS/internals/io"
	"github.com/kmanley/golang-tuple"
)

func testCreateRelationship(t *testing.T, IO io.IO) {
	var re = &en.RealEngine{IO}

	re.InitDatabase()
	defer re.DropDatabase()
	Init(re)

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
func testCreateRelationshipWithProperties(t *testing.T, IO io.IO) {
	var re = &en.RealEngine{IO}

	re.InitDatabase()
	defer re.DropDatabase()
	Init(re)

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
func testGetLabels(t *testing.T, IO io.IO) {
	var re = &en.RealEngine{IO}
	re.InitDatabase()
	defer re.DropDatabase()
	Init(re)

	node, _ := CreateNode("test_label")

	labels := *node.GetLabels()
	if len(labels) > 1 || len(labels) < 1 {
		t.Errorf("Wrong number of labels")
	}

	if labels[0] != "test_label" {
		t.Errorf("wrong label")
	}

}

func testGetNodeProperty(t *testing.T, IO io.IO) {
	var re = &en.RealEngine{IO}
	re.InitDatabase()
	defer re.DropDatabase()
	Init(re)

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
func testGetFromRelationship(t *testing.T, IO io.IO) {
	var re = &en.RealEngine{IO}
	re.InitDatabase()
	defer re.DropDatabase()
	Init(re)

	node1, _ := CreateNode("from", tuple.NewTupleFromItems("id", 1))
	node2, _ := CreateNode("to", tuple.NewTupleFromItems("id", 2))

	rel, _ := CreateRelationship(node1, node2, "link")

	retrivedNode := *rel.GetFrom()

	oldID, _ := node1.GetProperty("id")
	newID, _ := retrivedNode.GetProperty("id")

	if oldID != newID {
		t.Errorf("expected to be equal, but get %v and %v", oldID, newID)
	}
}

func testGetToRelationship(t *testing.T, IO io.IO) {
	var re = &en.RealEngine{IO}
	re.InitDatabase()
	defer re.DropDatabase()
	Init(re)

	node1, _ := CreateNode("from", tuple.NewTupleFromItems("id", 1))
	node2, _ := CreateNode("to", tuple.NewTupleFromItems("id", 2))

	rel, _ := CreateRelationship(node1, node2, "link")

	retrivedNode := *rel.GetTo()

	oldID, _ := node2.GetProperty("id")
	newID, _ := retrivedNode.GetProperty("id")

	if oldID != newID {
		t.Errorf("expected to be equal, but get %v and %v", oldID, newID)
	}
}

func testGetRelationshipType(t *testing.T, IO io.IO) {
	var re = &en.RealEngine{IO}
	re.InitDatabase()
	defer re.DropDatabase()

	Init(re)

	node1, _ := CreateNode("from", tuple.NewTupleFromItems("id", 1))
	node2, _ := CreateNode("to", tuple.NewTupleFromItems("id", 2))

	rel, _ := CreateRelationship(node1, node2, "link")

	relType := rel.GetType()

	if relType != "link" {
		t.Errorf("Wrong relationship type")
	}
}

func testGetRelationshipProperty(t *testing.T, IO io.IO) {
	var re = &en.RealEngine{IO}
	re.InitDatabase()
	defer re.DropDatabase()

	Init(re)

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

func testSelectNodesWhere(t *testing.T, IO io.IO) {
	var re = &en.RealEngine{IO}
	re.InitDatabase()
	defer re.DropDatabase()

	Init(re)

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

func testSelectRelationshipsWhere(t *testing.T, IO io.IO) {
	var re = &en.RealEngine{IO}
	re.InitDatabase()
	defer re.DropDatabase()

	Init(re)

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

func testGetRelationships(t *testing.T, IO io.IO) {
	var re = &en.RealEngine{IO}
	re.InitDatabase()
	defer re.DropDatabase()

	Init(re)

	node1, _ := CreateNode("a", tuple.NewTupleFromItems("id", 1))
	node2, _ := CreateNode("a", tuple.NewTupleFromItems("id", 2))

	_, ok := CreateRelationship(node1, node2, "show1")
	if !ok {
		t.Errorf("Could not create show 1")
	}
	_, ok = CreateRelationship(node1, node2, "show2")
	if !ok {
		t.Errorf("Could not create show 2")
	}
	_, ok = CreateRelationship(node1, node2, "show3")
	if !ok {
		t.Errorf("Could not create show 3")
	}
	_, ok = CreateRelationship(node1, node2, "show4")
	if !ok {
		t.Errorf("Could not create show 4")
	}

	relationships := *node1.GetRelationships()

	if len(relationships) != 4 {
		t.Errorf("Expected lenght: 4, get: %v", len(relationships))
	}

}

/// LocalIO tests

func TestCreateRelationshipWithLocalIO(t *testing.T) {
	testCreateRelationship(t, &io.LocalIO{})
}
func TestCreateRelationshipWithPropertiesWithLocalIO(t *testing.T) {
	testCreateRelationshipWithProperties(t, &io.LocalIO{})
}
func TestGetLabelsWithLocalIO(t *testing.T) {
	testGetLabels(t, &io.LocalIO{})
}
func TestGetNodePropertyWithLocalIO(t *testing.T) {
	testGetNodeProperty(t, &io.LocalIO{})
}
func TestGetFromRelationshipWithLocalIO(t *testing.T) {
	testGetFromRelationship(t, &io.LocalIO{})
}
func TestGetToRelationshipWithLocalIO(t *testing.T) {
	testGetToRelationship(t, &io.LocalIO{})
}
func TestGetRelationshipTypeWithLocalIO(t *testing.T) {
	testGetRelationshipType(t, &io.LocalIO{})
}
func TestGetRelationshipPropertyWithLocalIO(t *testing.T) {
	testGetRelationshipProperty(t, &io.LocalIO{})
}
func TestSelectNodesWhereWithLocalIO(t *testing.T) {
	testSelectNodesWhere(t, &io.LocalIO{})
}
func TestSelectRelationshipsWhereWithLocalIO(t *testing.T) {
	testSelectRelationshipsWhere(t, &io.LocalIO{})
}
func TestGetRelationshipsWithLocalIO(t *testing.T) {
	testGetRelationships(t, &io.LocalIO{})
}
