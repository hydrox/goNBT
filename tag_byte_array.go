package goNBT

import (
	"encoding/binary"
	"io"
	"strconv"
	"strings"
)

type TAG_Byte_Array struct {
	*TAG_Named
	Length uint32
	Array  []byte
}

func (tag TAG_Byte_Array) string(prefix int) string {
	var name string
	if tag.TAG_Named.GetName() == "" {
		name = "TAG_Byte_Array: "
	} else {
		name = "TAG_Byte_Array(\"" + tag.TAG_Named.GetName() + "\"): "
	}
	//	string2 := []string{getPrefix(prefix), name, string(tag.Array), ": [", strconv.Itoa(len(tag.Array)), " bytes]\n"}
	string2 := []string{getPrefix(prefix), name, "[", strconv.Itoa(len(tag.Array)), " bytes]\n"}
	output := strings.Join(string2, "")
	return output
}

func (tag TAG_Byte_Array) Write(writer io.WriteCloser) (err error) {
	binary.Write(writer, binary.BigEndian, TAG_BYTE_ARRAY)
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

func (tag TAG_Byte_Array) WritePayload(writer io.WriteCloser) (err error) {
	err = binary.Write(writer, binary.BigEndian, int32(tag.Length))
	if err != nil {
		return err
	}
	err = binary.Write(writer, binary.BigEndian, tag.Array)
	return err
}
