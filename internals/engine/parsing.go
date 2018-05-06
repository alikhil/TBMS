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

func parse(data []byte, typ EType) (ret interface{}) {
	buf := bytes.NewBuffer(data)
	var err error
	if typ == Estring || typ == Eint {
		var i32 int32 = 4
		err = binary.Read(buf, ConventionByteOrder, &i32)
		ret = i32
	} else if typ == Ebool {
		var iret int32 = 0
		err = binary.Read(buf, ConventionByteOrder, &iret)
		ret = true
		if iret == 0 {
			ret = false
		}
	} else if typ == Efloat {
		var il float32 = 4.0
		err = binary.Read(buf, ConventionByteOrder, &il)
		ret = il
	}
	if err != nil {
		logger.Error.Printf("Can not parse interface{} %v", err)
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

func (l *ELabelString) fill(data *[]byte, id int32) {

	end := bytes.IndexByte(*data, 0)
	if end == -1 {
		end = len(*data)
	}
	l.ID = id
	l.String = string((*data)[1:end])
}

func (rel *ERelationship) fill(data *[]byte, id int32) {
	rel.ID = id
	rel.FirstInChain = parseBool((*data)[1])
	rel.SecondNodeID = parseInt((*data)[2:6])
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

func (r *ERelationshipType) fill(data *[]byte, id int32) {
	r.ID = id
	end := bytes.IndexByte(*data, 0)
	if end == -1 {
		end = len(*data)
	}
	r.TypeString = string((*data)[1:end])
}

func (r *EString) fill(data *[]byte, id int32) {
	r.ID = id
	r.Extra = (*data)[1]
	r.NextPartID = parseInt((*data)[2:6])
	buf := (*data)[6:64]
	r.Value = &buf
}

func (r *EProperty) fill(data *[]byte, id int32) {
	r.ID = id
	r.Typename = EType((*data)[1])
	r.KeyStringID = parseInt((*data)[2:6])
	r.NextPropertyID = parseInt((*data)[6:10])
	r.ValueOrStringPtr = parse((*data)[10:], r.Typename)
}

func (r *EPropertyKey) fill(data *[]byte, id int32) {
	end := bytes.IndexByte(*data, 0)
	if end == -1 {
		end = len(*data)
	}
	r.ID = id
	r.KeyString = string((*data)[1:end])
}

func (r *ELabel) fill(data *[]byte, id int32) {
	r.ID = id

	r.LabelStringID = parseInt((*data)[1:5])
	r.NextLabelID = parseInt((*data)[5:])
}
