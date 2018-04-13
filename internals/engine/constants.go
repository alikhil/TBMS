package engine

import (
	"encoding/binary"
)

/**
Record length for each store
*/
const (
	BytesPerNode         = 13
	BytesPerRelationship = 34
	BytesPerProperty     = 10
	BytesPerString       = 64
	BytesPerLabel        = 9
	BytesPerLabelString  = 21
	BytesPerPropertyKey  = 21
	BytesPerRelType      = 21
	BytesPerInUse        = 11
)

/**
File names
*/
const (
	FNNodes             = "nodes.store"
	FNLabels            = "labels.store"
	FNLabelsStrings     = "labelsStrings.store"
	FNRelationships     = "relationships.store"
	FNProperties        = "properties.store"
	FNStrings           = "strings.store"
	FNInUse             = "inuse.store"
	FNPropertyKeys      = "propertykeys.store"
	FNRelationshipTypes = "relationshiptypes.store"
	// TODO: add others
)

type EType byte

/**
TypeNames
*/
const (
	EInt EType = iota + 1
	EString
	EFloat
	EBool
)

/**
Byte order to use in encoding/decoding
*/
var (
	ConventionByteOrder = binary.LittleEndian
)

type EStore byte

const (
	StoreInUse EStore = iota
	StoreNode
	StoreRelationship
	StoreProperty
	StoreLabel
	StoreLabelString
	StorePropertyKey
	StoreString
	StoreRelationshipType
)

// FilenameStore maps StoreId with filename where it stores
var FilenameStore = map[EStore]string{
	StoreInUse:            FNInUse,
	StoreNode:             FNNodes,
	StoreProperty:         FNProperties,
	StoreLabel:            FNLabels,
	StoreRelationship:     FNRelationships,
	StoreLabelString:      FNLabelsStrings,
	StorePropertyKey:      FNPropertyKeys,
	StoreString:           FNStrings,
	StoreRelationshipType: FNRelationshipTypes,
}

var BytesPerStore = map[EStore]int32{
	StoreNode:             BytesPerNode,
	StoreRelationship:     BytesPerRelationship,
	StoreProperty:         BytesPerProperty,
	StoreString:           BytesPerString,
	StoreLabel:            BytesPerLabel,
	StoreLabelString:      BytesPerLabelString,
	StorePropertyKey:      BytesPerPropertyKey,
	StoreRelationshipType: BytesPerRelType,
	StoreInUse:            BytesPerInUse,
}
