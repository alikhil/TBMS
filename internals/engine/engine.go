package engine

// types starting with E are used only within Engine

type Engine interface {
	IterateObjects(labelID int) func() EObject
	IterateRelationships(relTypeID int) func() ERelationship
	IterateObjectProperties(objID int) func() EProperty
	IterateRelationshipProperties(relID int) func() EProperty
	IterateObjectLabels(objID int) func() string

	AddObjectLabel(objID, labelID int)
	DeleteObjectLabel(objID, labelID int)
	// AddObjectProperty(objID int, key string, val bool) how to use without Generics?

	GetLabelByID(int) (string, error) // error if no such labelid
	GetLabelId(string) (int, error)   // error if not exists
	CreateLabel(string) int           // by name, return id

	GetRelationshipID(string) (int, error) // error if no such relationship
	CreateRelationship(string) int         // by type, return id
}

type EObject struct {
	ID int
	// nextLabelId?
	// labels []string
	// properties []EProperty
}

func (o EObject) GetId() int {
	return o.ID
}

func (o EObject) SetId(newId int) {
	o.ID = newId
}

type EProperty struct {
	id  int
	key string
	// propertyId?
	// getInt? getBool? getString?
}

type ERelationship struct {
	id             int
	firstObjectID  int
	secondObjectID int
	typeID         int
}
