package goNBT

import (
	"container/list"
	"encoding/binary"
	"io"
	"strconv"
	"strings"
)

type TAG_List struct {
	*TAG_Named
	TagId  byte
	Length int32
	List   *list.List
}

func (tag TAG_List) String() string {
	return tag.string(0)
}

func (tag TAG_List) string(prefix int) string {
	var name string
	if tag.TAG_Named.GetName() == "" {
		name = "TAG_List: "
	} else {
		name = "TAG_List(\"" + tag.TAG_Named.GetName() + "\"): "
	}
	string2 := []string{getPrefix(prefix), name, strconv.Itoa(tag.List.Len()), " entries of Type ", strconv.Itoa(int(tag.TagId)), "\n"}
	output := strings.Join(string2, "")
	for e := tag.List.Front(); e != nil; e = e.Next() {
		tmp := toString(e.Value, prefix+1)
		string2 := []string{output, tmp}
		output = strings.Join(string2, "")
	}
	return output
}

func (tag TAG_List) GetByteEntries() (entry []*TAG_Byte, success bool) {
	if tag.TagId != TAG_BYTE {
		return entry, false
	}
	entry = make([]*TAG_Byte, tag.List.Len())
	i := 0
	for e := tag.List.Front(); e != nil; e = e.Next() {
		switch t := e.Value.(type) {
		case TAG_Byte:
			entry[i] = &t
		}
		i++
	}
	return entry, true
}

func (tag TAG_List) GetShortEntries() (entry []TAG_Short, success bool) {
	if tag.TagId != TAG_SHORT {
		return entry, false
	}
	entry = make([]TAG_Short, tag.List.Len())
	i := 0
	for e := tag.List.Front(); e != nil; e = e.Next() {
		switch t := e.Value.(type) {
		case TAG_Short:
			entry[i] = t
		}
		i++
	}
	return entry, true
}

func (tag TAG_List) GetIntEntries() (entry []*TAG_Int, success bool) {
	if tag.TagId != TAG_INT {
		return entry, false
	}
	entry = make([]*TAG_Int, tag.List.Len())
	i := 0
	for e := tag.List.Front(); e != nil; e = e.Next() {
		switch t := e.Value.(type) {
		case TAG_Int:
			entry[i] = &t
		}
		i++
	}
	return entry, true
}

func (tag TAG_List) GetLongEntries() (entry []*TAG_Long, success bool) {
	if tag.TagId != TAG_LONG {
		return entry, false
	}
	entry = make([]*TAG_Long, tag.List.Len())
	i := 0
	for e := tag.List.Front(); e != nil; e = e.Next() {
		switch t := e.Value.(type) {
		case TAG_Long:
			entry[i] = &t
		}
		i++
	}
	return entry, true
}

func (tag TAG_List) GetFloatEntries() (entry []*TAG_Float, success bool) {
	if tag.TagId != TAG_FLOAT {
		return entry, false
	}
	entry = make([]*TAG_Float, tag.List.Len())
	i := 0
	for e := tag.List.Front(); e != nil; e = e.Next() {
		switch t := e.Value.(type) {
		case TAG_Float:
			entry[i] = &t
		}
		i++
	}
	return entry, true
}

func (tag TAG_List) GetDoubleEntries() (entry []*TAG_Double, success bool) {
	if tag.TagId != TAG_DOUBLE {
		return entry, false
	}
	entry = make([]*TAG_Double, tag.List.Len())
	i := 0
	for e := tag.List.Front(); e != nil; e = e.Next() {
		switch t := e.Value.(type) {
		case TAG_Double:
			entry[i] = &t
		}
		i++
	}
	return entry, true
}

func (tag TAG_List) GetByteArrayEntries() (entry []*TAG_Byte_Array, success bool) {
	if tag.TagId != TAG_BYTE_ARRAY {
		return entry, false
	}
	entry = make([]*TAG_Byte_Array, tag.List.Len())
	i := 0
	for e := tag.List.Front(); e != nil; e = e.Next() {
		switch t := e.Value.(type) {
		case TAG_Byte_Array:
			entry[i] = &t
		}
		i++
	}
	return entry, true
}

func (tag TAG_List) GetStringEntries() (entry []*TAG_String, success bool) {
	if tag.TagId != TAG_STRING {
		return entry, false
	}
	entry = make([]*TAG_String, tag.List.Len())
	i := 0
	for e := tag.List.Front(); e != nil; e = e.Next() {
		switch t := e.Value.(type) {
		case TAG_String:
			entry[i] = &t
		}
		i++
	}
	return entry, true
}

func (tag TAG_List) GetListEntries() (entry []*TAG_List, success bool) {
	if tag.TagId != TAG_LIST {
		return entry, false
	}
	entry = make([]*TAG_List, tag.List.Len())
	i := 0
	for e := tag.List.Front(); e != nil; e = e.Next() {
		switch t := e.Value.(type) {
		case TAG_List:
			entry[i] = &t
		}
		i++
	}
	return entry, true
}

func (tag TAG_List) GetCompoundEntries() (entry []*TAG_Compound, success bool) {
	if tag.TagId != TAG_COMPOUND {
		return entry, false
	}
	entry = make([]*TAG_Compound, tag.List.Len())
	i := 0
	for e := tag.List.Front(); e != nil; e = e.Next() {
		switch t := e.Value.(type) {
		case TAG_Compound:
			entry[i] = &t
		}
		i++
	}
	return entry, true
}

func (tag TAG_List) GetIntArrayEntries() (entry []*TAG_Int_Array, success bool) {
	if tag.TagId != TAG_INT_ARRAY {
		return entry, false
	}
	entry = make([]*TAG_Int_Array, tag.List.Len())
	i := 0
	for e := tag.List.Front(); e != nil; e = e.Next() {
		switch t := e.Value.(type) {
		case TAG_Int_Array:
			entry[i] = &t
		}
		i++
	}
	return entry, true
}

func (tag TAG_List) Write(writer io.WriteCloser) (err error) {
	binary.Write(writer, binary.BigEndian, TAG_LIST)
	if tag.TAG_Named.GetName() != "" {
		err = binary.Write(writer, binary.BigEndian, int16(len(tag.GetName())))
		if err != nil {
			return err
		}
		_, err = writer.Write([]byte(tag.GetName()))
		if err != nil {
			return err
		}
	}
	return tag.WritePayload(writer)
}

func (tag TAG_List) WritePayload(writer io.WriteCloser) (err error) {
	err = binary.Write(writer, binary.BigEndian, byte(tag.TagId))
	err = binary.Write(writer, binary.BigEndian, int32(tag.Length))
	for e := tag.List.Front(); e != nil; e = e.Next() {
		value := e.Value.(TAG_Writeable)
		err = value.WritePayload(writer)
		if err != nil {
			return err
		}
	}
	return err
}
