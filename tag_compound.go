package goNBT

import (
	"container/list"
	"encoding/binary"
	"io"
	"strconv"
	"strings"
)

type TAG_Compound struct {
	*TAG_Named
	List *list.List
}

func (tag TAG_Compound) string(prefix int) string {
	var name string
	if tag.TAG_Named.GetName() == "" {
		name = "TAG_Compound: "
	} else {
		name = "TAG_Compound(\"" + tag.TAG_Named.GetName() + "\"): "
	}
	string2 := []string{getPrefix(prefix), name, strconv.Itoa(tag.List.Len()), " entries\n"}
	output := strings.Join(string2, "")
	for e := tag.List.Front(); e != nil; e = e.Next() {
		tmp := toString(e.Value, prefix+1)
		string2 := []string{output, tmp}
		output = strings.Join(string2, "")
	}
	return output
}

func (tag TAG_Compound) String() string {
	return tag.string(0)
}

func (tag TAG_Compound) getEntry(name string) (entry TAG_Name, success bool) {
	for e := tag.List.Front(); e != nil; e = e.Next() {
		switch t := e.Value.(type) {
		case TAG_Name:
			entry = t
		}
		if entry.GetName() == name {
			return entry, true
		}
	}
	return entry, false
}

func (tag TAG_Compound) GetByte(name string) (entry TAG_Byte, success bool) {
	possibleEntry, success := tag.getEntry(name)
	switch t := possibleEntry.(type) {
	case TAG_Byte:
		entry = t
		if entry.GetName() == name {
			return entry, true
		}
	}
	return entry, false
}

func (tag TAG_Compound) GetShort(name string) (entry TAG_Short, success bool) {
	possibleEntry, success := tag.getEntry(name)
	switch t := possibleEntry.(type) {
	case TAG_Short:
		entry = t
		if entry.GetName() == name {
			return entry, true
		}
	}
	return entry, false
}

func (tag TAG_Compound) GetInt(name string) (entry TAG_Int, success bool) {
	possibleEntry, success := tag.getEntry(name)
	switch t := possibleEntry.(type) {
	case TAG_Int:
		entry = t
		if entry.GetName() == name {
			return entry, true
		}
	}
	return entry, false
}

func (tag TAG_Compound) GetLong(name string) (entry TAG_Long, success bool) {
	possibleEntry, success := tag.getEntry(name)
	switch t := possibleEntry.(type) {
	case TAG_Long:
		entry = t
		if entry.GetName() == name {
			return entry, true
		}
	}
	return entry, false
}

func (tag TAG_Compound) GetFloat(name string) (entry TAG_Float, success bool) {
	possibleEntry, success := tag.getEntry(name)
	switch t := possibleEntry.(type) {
	case TAG_Float:
		entry = t
		if entry.GetName() == name {
			return entry, true
		}
	}
	return entry, false
}

func (tag TAG_Compound) GetDouble(name string) (entry TAG_Double, success bool) {
	possibleEntry, success := tag.getEntry(name)
	switch t := possibleEntry.(type) {
	case TAG_Double:
		entry = t
		if entry.GetName() == name {
			return entry, true
		}
	}
	return entry, false
}

func (tag TAG_Compound) GetByteArray(name string) (entry TAG_Byte_Array, success bool) {
	possibleEntry, success := tag.getEntry(name)
	switch t := possibleEntry.(type) {
	case TAG_Byte_Array:
		entry = t
		if entry.GetName() == name {
			return entry, true
		}
	}
	return entry, false
}

func (tag TAG_Compound) GetString(name string) (entry TAG_String, success bool) {
	possibleEntry, success := tag.getEntry(name)
	switch t := possibleEntry.(type) {
	case TAG_String:
		entry = t
		if entry.GetName() == name {
			return entry, true
		}
	}
	return entry, false
}

func (tag TAG_Compound) GetList(name string) (entry TAG_List, success bool) {
	possibleEntry, success := tag.getEntry(name)
	switch t := possibleEntry.(type) {
	case TAG_List:
		entry = t
		if entry.GetName() == name {
			return entry, true
		}
	}
	return entry, false
}

func (tag TAG_Compound) GetCompound(name string) (entry TAG_Compound, success bool) {
	possibleEntry, success := tag.getEntry(name)
	switch t := possibleEntry.(type) {
	case TAG_Compound:
		entry = t
		if entry.GetName() == name {
			return entry, true
		}
	}
	return entry, false
}

func (tag TAG_Compound) GetIntArray(name string) (entry TAG_Int_Array, success bool) {
	possibleEntry, success := tag.getEntry(name)
	switch t := possibleEntry.(type) {
	case TAG_Int_Array:
		entry = t
		if entry.GetName() == name {
			return entry, true
		}
	}
	return entry, false
}

func (tag TAG_Compound) Write(writer io.WriteCloser) (err error) {
	binary.Write(writer, binary.BigEndian, TAG_COMPOUND)
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

func (tag TAG_Compound) WritePayload(writer io.WriteCloser) (err error) {
	for e := tag.List.Front(); e != nil; e = e.Next() {
		value := e.Value.(TAG_Writeable)
		err = value.Write(writer)
		if err != nil {
			return err
		}
	}
	return binary.Write(writer, binary.BigEndian, TAG_END)
}
