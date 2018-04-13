package engine

import (
	"bytes"
	"encoding/binary"
	"github.com/alikhil/TBMS/internals/logger"
)

func parseInt(data []byte) (ret int32) {
	buf := bytes.NewBuffer(data)
	err := binary.Read(buf, ConventionByteOrder, &ret)
	if err != nil {
		logger.Error.Printf("Can not parse int %v", err)
	}
	return ret
}

func parseBool(b byte) bool {
	return b > 0
}

func (node *ENode) fill(data *[]byte, id int32) {
	node.ID = id
	node.NextLabelID = parseInt((*data)[1:5])
	node.NextPropertyID = parseInt((*data)[5:9])
	node.NextRelID = parseInt((*data)[9:13])
}

func parseProperty(data *[]byte) (*EProperty, bool) {
	var inUse = parseBool((*data)[0])
	if !inUse {
		return nil, false
	}
	return &EProperty{
		Typename:         EType((*data)[1]),
		KeyStringID:      parseInt((*data)[2:6]),
		ValueOrStringPtr: parseInt((*data)[6:10]),
	}, true
}

func (l *ELabelString) fill(data *[]byte, id int32) {

	end := bytes.IndexByte(*data, 0)
	l.ID = id
	l.String = string((*data)[1:end])
}

func (rel *ERelationship) fill(data *[]byte, id int32) {
	rel.ID = id
	rel.FirstInChain = parseBool((*data)[1])
	rel.SecondNode = parseInt((*data)[2:6])
	rel.FirstNodeID = parseInt((*data)[6:10])
	rel.FirstNodeNxtRelID = parseInt((*data)[10:14])
	rel.SecondNodeNxtRelID = parseInt((*data)[14:18])
	rel.FirstNodePrvRelID = parseInt((*data)[18:22])
	rel.SecondNodePrvRelID = parseInt((*data)[22:26])
	rel.NextPropertyID = parseInt((*data)[26:30])
	rel.TypeID = parseInt((*data)[30:34])
}

func (r *EInUseRecord) fill(data *[]byte, id int32) {
	r.ID = id
	r.StoreType = EStore((*data)[1])
	r.IsHead = parseBool((*data)[2])
	r.ObjID = parseInt((*data)[3:7])
	r.NextRecordID = parseInt((*data)[7:11])
}
