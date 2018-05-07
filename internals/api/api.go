package api

import (
	"reflect"

	en "github.com/alikhil/TBMS/internals/engine"
	"github.com/alikhil/TBMS/internals/logger"
	"github.com/kmanley/golang-tuple"
)

var engine en.RealEngine

// Init should be called first
func Init(eng *en.RealEngine) {
	engine = *eng
}

func findOrCreateRelationType(relTypeString string) (relTypeID int32, ok bool) {
	return engine.FindOrCreateObject(en.StoreRelationshipType,
		func(ob en.EObject) bool { return ob.(*en.ERelationshipType).TypeString == relTypeString },
		func(id int32) en.EObject { return &en.ERelationshipType{ID: id, TypeString: relTypeString} })
}

func CreateNodeLabel(nodeLabel string) (nodeLabelID int32, ok bool) {
	// label string can already exists
	labelStringID, ok := findOrCreateLabelString(nodeLabel)
	if !ok {
		return -1, false
	}

	// but label for the new node must be new
	return engine.CreateObject(en.StoreLabel,
		func(id int32) en.EObject { return &en.ELabel{ID: id, LabelStringID: labelStringID, NextLabelID: -1} })
}

func findOrCreateLabelString(labelString string) (labelStringID int32, ok bool) {
	return engine.FindOrCreateObject(en.StoreLabelString,
		func(ob en.EObject) bool { return ob.(*en.ELabelString).String == labelString },
		func(id int32) en.EObject { return &en.ELabelString{ID: id, String: labelString} })
}

func validateProps(props []*tuple.Tuple) (*[]string, *[]interface{}, bool) {
	propsCnt := len(props)
	keys := make([]string, propsCnt)
	vals := make([]interface{}, propsCnt)
	// validate props

	for i := 0; i < propsCnt; i++ {
		if props[i].Len() != 2 {
			logger.Error.Printf("property tuples should have 2 values, but have %v", props[i].Len())
			return nil, nil, false
		}

		key, typeMatch := props[i].Get(0).(string)
		if !typeMatch {
			logger.Error.Printf("%v-th property has invalid key", i)
			return nil, nil, false
		}
		keys[i] = key
		vals[i] = props[i].Get(1)
	}
	return &keys, &vals, true
}

func getProperties(nextPropertyID int32) (*map[string]interface{}, bool) {
	dict := make(map[string]interface{})

	if nextPropertyID == -1 {
		return &dict, true
	}

	property := &en.EProperty{ID: nextPropertyID}
	for ok := engine.GetObject(property); ok; ok = engine.GetObject(property) {
		pkey := &en.EPropertyKey{ID: property.KeyStringID}
		pkeyFound := engine.GetObject(pkey)

		if !pkeyFound {
			logger.Error.Printf("Can not load property key with id %v", property.KeyStringID)
			return nil, false
		}
		dict[pkey.KeyString] = property.ValueOrStringPtr

		if property.Typename == en.Estring {
			str := &en.EString{ID: property.ValueOrStringPtr.(int32)}
			valStrFound := engine.GetObject(str)

			if !valStrFound {
				logger.Error.Printf("Can not found value of string with id %v", str.ID)
				return nil, false
			}

			dict[pkey.KeyString] = str.LoadString(&engine)
		}

		property.ID = property.NextPropertyID
		if property.ID == -1 {
			break
		}
	}

	return &dict, true
}

func fillPropertyValue(prop *en.EProperty, val interface{}) (ok bool) {
	ok = true
	switch t := val.(type) {
	case int:
		prop.Typename = en.Eint
		prop.ValueOrStringPtr = int32(t)
	case int32:
		prop.Typename = en.Eint
		prop.ValueOrStringPtr = t
	case float32:
		prop.Typename = en.Efloat
		prop.ValueOrStringPtr = t
	case bool:
		prop.Typename = en.Ebool
		prop.ValueOrStringPtr = int32(1)
		if !t {
			prop.ValueOrStringPtr = int32(0)
		}
	case string:
		prop.Typename = en.Estring
		str := &en.EString{ID: -1}
		found := engine.FindObject(en.StoreString,
			func(ob en.EObject) bool {
				return ob.(*en.EString).LoadString(&engine) == t
			}, str)
		if !found {
			str = engine.CreateStringAndReturnFirstChunk(t)
		}
		prop.ValueOrStringPtr = str.ID
	}
	return ok
}

func fillProperties(props []*tuple.Tuple) (int32, bool) {
	propsCnt := len(props)
	if propsCnt == 0 {
		// empty list
		return -1, true
	}

	keys, values, isValid := validateProps(props)
	if !isValid {
		return -1, false
	}

	propKeyIDs := make([]int32, propsCnt)
	for i := 0; i < propsCnt; i++ {
		propertyKeyID, ok := engine.FindOrCreateObject(en.StorePropertyKey,
			func(ob en.EObject) bool { return ob.(*en.EPropertyKey).KeyString == (*keys)[i] },
			func(id int32) en.EObject { return &en.EPropertyKey{ID: id, KeyString: (*keys)[i]} })

		if !ok {
			return -1, false
		}
		propKeyIDs[i] = propertyKeyID
	}

	var property = &en.EProperty{ID: -1}
	for i := 0; i < propsCnt; i++ {
		id, ok := engine.GetAndLockFreeIDForStore(en.StoreProperty)
		if !ok {
			return -1, false
		}
		property = &en.EProperty{ID: id, KeyStringID: propKeyIDs[i], NextPropertyID: property.ID}
		filled := fillPropertyValue(property, (*values)[i])
		if !filled {
			logger.Error.Printf("some shit happened!!! failed to fill property!")
			return -1, false
		}
		engine.SaveObject(property)
	}

	return property.ID, true
}

func getLastRelationship(node *en.ENode) (*en.ERelationship, bool) {
	if node.NextRelID == -1 {
		return nil, false
	}
	fillNext := engine.GetNodeRelationshipsIteratorStartingFrom(node.ID, node.NextRelID)

	var prev *en.ERelationship
	for res, ok := fillNext(); ok; res, ok = fillNext() {
		prev = res
	}

	if prev != nil {
		if prev.GetPart(node.ID).NodeNxtRelID != -1 {
			// Then why we stopped here, if there is one more next relation? ¡PANIC!
			panic("iterator stopped before last relationship")
		}
		return prev, true
	}
	return nil, false

}

func CreateNode(nodeLabel string, properties ...*tuple.Tuple) (*Node, bool) {
	nodeLabelID, ok := CreateNodeLabel(nodeLabel)
	if !ok {
		return nil, false
	}

	nextPropertyID, ok := fillProperties(properties)
	if !ok {
		return nil, false
	}

	nodeID, ok := engine.GetAndLockFreeIDForStore(en.StoreNode)
	if !ok {
		return nil, false
	}

	node := &en.ENode{
		ID:             nodeID,
		NextRelID:      -1,
		NextPropertyID: nextPropertyID,
		NextLabelID:    nodeLabelID,
	}

	engine.SaveObject(node)
	return &Node{node}, true
}

func CreateRelationship(a, b *Node, relType string, properties ...*tuple.Tuple) (*Relationship, bool) {

	if a == nil || b == nil {
		return nil, false
	}

	// check if such relType exist, obtain id
	relTypeID, ok := findOrCreateRelationType(relType)

	if !ok {
		return nil, false
	}

	nextPropertyID, ok := fillProperties(properties)
	if !ok {
		return nil, false
	}

	relID, ok := engine.GetAndLockFreeIDForStore(en.StoreRelationship)
	if !ok {
		return nil, false
	}

	relationship := &en.ERelationship{
		ID:                 relID,
		FirstNodeID:        a.ID,
		SecondNodeID:       b.ID,
		TypeID:             relTypeID,
		NextPropertyID:     nextPropertyID,
		FirstNodeNxtRelID:  -1,
		SecondNodeNxtRelID: -1,
	}

	// constructing links of two double linked lists

	aLastRel, foundA := getLastRelationship(a.ENode)
	bLastRel, foundB := getLastRelationship(b.ENode)
	if foundA {
		relationship.FirstNodePrvRelID = aLastRel.ID
		aLastRel.SetNextRelationshipID(a.ID, relID)
		if !engine.SaveObject(aLastRel) {
			panic("failed to save object")
		}
	}

	if foundB {
		relationship.SecondNodePrvRelID = bLastRel.ID
		bLastRel.SetNextRelationshipID(b.ID, relID)
		if !engine.SaveObject(bLastRel) {
			panic("failed to save object")
		}
	}

	engine.SaveObject(relationship)
	return &Relationship{relationship}, true
}

func SelectNodesWhere(condition func(*Node) bool) ([]*Node, bool) {
	fillNextObj := engine.GetEObjectIterator(en.StoreNode)

	result := make([]*Node, 0)

	res := &en.ENode{}

	for ok := fillNextObj(res); ok; ok = fillNextObj(res) {
		a := &Node{res}
		if condition(a) {
			result = append(result, a)
		}
		// problem with rewriting links
		res = &en.ENode{}
	}
	return result, true

}

func SelectRelationshipWhere(condition func(*Relationship) bool) ([]*Relationship, bool) {
	fillNextObj := engine.GetEObjectIterator(en.StoreRelationship)

	result := make([]*Relationship, 0)

	res := &en.ERelationship{}

	for ok := fillNextObj(res); ok; ok = fillNextObj(res) {
		a := Relationship{res}
		if condition(&a) {
			result = append(result, &Relationship{res})
		}
		// problem with rewritnig links
		res = &en.ERelationship{}
	}
	return result, true
}

type Node struct {
	*en.ENode
}

// GetLabels returns all labels of the given node
func (node *Node) GetLabels() *[]string {

	res := make([]string, 0)
	label := &en.ELabel{ID: node.NextLabelID}
	for engine.GetObject(label) {
		labelStr := &en.ELabelString{ID: label.LabelStringID}
		engine.GetObject(labelStr)
		res = append(res, labelStr.String)
		if label.NextLabelID != -1 {
			label = &en.ELabel{ID: label.NextLabelID}
		} else {
			break
		}
	}

	return &res
}

// GetProperty returns property value by key, if exists, else return false in bool
func (node *Node) GetProperty(key string) (interface{}, bool) {
	p, ok := node.GetProperties()
	prop := *p
	if !ok {
		return nil, false
	}
	obj, ok := prop[key]
	return obj, ok
}

// GetProperties returns all properties as a map
func (node *Node) GetProperties() (*map[string]interface{}, bool) {
	return getProperties(node.NextPropertyID)
}

// GetRelationships returns pointer to a slice of all relationships, this node belonges to
func (node *Node) GetRelationships() *[]Relationship {

	relationships := make([]Relationship, 0)
	if node.NextRelID == -1 {
		return &relationships
	}

	fillNext := engine.GetNodeRelationshipsIteratorStartingFrom(node.ID, node.NextRelID)

	for res, ok := fillNext(); ok; res, ok = fillNext() {
		relationships = append(relationships, Relationship{res})
	}

	return &relationships
}

type Relationship struct {
	*en.ERelationship
	// GetProperties() ?
}

// GetFrom returns node from wich relationship starts
func (rel *Relationship) GetFrom() *Node {
	node := &en.ENode{ID: rel.FirstNodeID}
	if engine.GetObject(node) {
		return &Node{node}
	}
	logger.Error.Print("Node doesn't exists")
	return nil
}

// GetTo returns node with which relationship ends
func (rel *Relationship) GetTo() *Node {
	node := &en.ENode{ID: rel.SecondNodeID}
	if engine.GetObject(node) {
		return &Node{node}
	}
	logger.Error.Print("Node doesn't exists")
	return nil
}

// GetType returns type of the relationship
func (rel *Relationship) GetType() string {
	relType := &en.ERelationshipType{ID: rel.TypeID}
	if engine.GetObject(relType) {
		return relType.TypeString
	}
	logger.Error.Print("Node doesn't exists")
	return ""
}

// GetProperty returns property value by key, if exists, else return false in bool
func (rel *Relationship) GetProperty(key string) (interface{}, bool) {
	p, ok := rel.GetProperties()
	prop := *p
	if !ok {
		return nil, false
	}
	obj, ok := prop[key]
	return obj, ok
}

func (rel *Relationship) GetProperties() (*map[string]interface{}, bool) {
	return getProperties(rel.NextPropertyID)
}

func Contains(list *[]interface{}, obj interface{}) bool {
	for _, a := range *list {
		if reflect.DeepEqual(a, obj) {
			return true
		}
	}
	return false
}
