package api

import (
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

func fillPropertyValue(prop *en.EProperty, val interface{}) (ok bool) {
	ok = true
	switch t := val.(type) {
	case int32:
		prop.Typename = en.Eint
		prop.ValueOrStringPtr = t
	case float32:
		prop.Typename = en.Efloat
		prop.ValueOrStringPtr = t
	case bool:
		prop.Typename = en.Ebool
		prop.ValueOrStringPtr = t
	case string:
		prop.Typename = en.Estring
		prop.ValueOrStringPtr, ok = engine.FindOrCreateObject(en.StoreString,
			func(ob en.EObject) bool { return ob.(*en.EString).LoadString(&engine) == t },
			func(id int32) en.EObject { return engine.CreateStringAndReturnFirstChunk(t) })
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
		fillPropertyValue(property, (*values)[i])
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
	// for each property key check and obtain id

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

type Node struct {
	*en.ENode
}

type Relationship struct {
	*en.ERelationship
	// GetProperties() ?
}
