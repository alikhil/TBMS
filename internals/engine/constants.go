package engine

import (
	"encoding/binary"
)

const (
	BytesPerNode         = 13
	BytesPerRelationship = 34
	BytesPerProperty     = 10
	BytesPerString       = 64
	BytesPerLabel        = 9
	BytesPerLabelString  = 21
	BytesPerPropertyKey  = 21
	BytesPerRelType      = 21
	BytesPerInUse        = 10
)

/**
File names
*/
const (
	FNNodes         = "nodes"
	FNLabels        = "labels"
	FNLabelsStrings = "labelsStrings"
	FNRelationships = "relationships"
	FNProperties    = "properties"
	FNStrings       = "strings"
	// TODO: add others
)

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
