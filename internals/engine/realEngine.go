package engine

// TODO: pass all arrays and slices by reference

func (re *RealEngine) GetLabelIteratorFromId(labelID int) func() (int, bool) {
	return func() (int, bool) {
		return 0, false
	}
}

func (re *RealEngine) GetNodesByLabelIterator(label string) func() (*ENode, bool) {
	nextNode := re.GetNodesIterator()
	neededlabelID, ok := re.GetLabelID(label)
	if !ok {

		return func() (*ENode, bool) {
			return nil, false
		}
	}

	return func() (*ENode, bool) {
		node, ok := nextNode()
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

func (re *RealEngine) GetNodesIterator() func() (*ENode, bool) {
	next := re.GetObjectIterator(FNNodes, BytesPerNode)
	i := 0
	return func() (*ENode, bool) {
		data, ok := next()
		if ok {
			node, nodeInUse := parseNode(&data, i)
			if !nodeInUse {
				// Trying to access deleted or nonexisting node
				return nil, false
			}
			i++
			return node, ok
		}
		return nil, false
	}
}

func (re *RealEngine) GetLabelStringIterator() func() (*ELabelString, bool) {
	next := re.GetObjectIterator(FNLabelsStrings, BytesPerLabelString)
	i := 0
	return func() (*ELabelString, bool) {
		data, ok := next()
		if ok {
			label, labelStringInUse := parseLabelString(&data, i)
			if !labelStringInUse {
				return nil, false
			}
			i++
			return label, ok
		}
		return nil, false
	}
}

func (re *RealEngine) GetRelationshiptIterator() func() (*ERelationship, bool) {
	next := re.GetObjectIterator(FNRelationships, BytesPerRelationship)
	i := 0
	return func() (*ERelationship, bool) {
		data, ok := next()
		if ok {
			rel, relInUse := parseRelationship(&data, i)
			if !relInUse {
				return nil, false
			}
			i++
			return rel, ok
		}
		return nil, false
	}
}

func (re *RealEngine) GetObjectIterator(filename string, recordLength int) func() ([]byte, bool) {
	curOffset := 0
	return func() (data []byte, ok bool) {
		data, ok = re.IO.ReadBytes(filename, curOffset, recordLength)
		if ok {
			curOffset += BytesPerNode
		}
		return
	}
}
func (re *RealEngine) GetLabelID(label string) (int, bool) {
	next := re.GetLabelStringIterator()
	i := 0
	for l, ok := next(); ok; {
		if ok && label == l.String {
			return i, true
		}
		i++
	}
	return -1, false
}
