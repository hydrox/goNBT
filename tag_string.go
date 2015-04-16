package goNBT

import (
	"encoding/binary"
	"io"
	"strconv"
	"strings"
)

type TAG_String struct {
	*TAG_Named
	Length uint16
	Array  []byte
}

func (tag TAG_String) string(prefix int) string {
	var name string
	if tag.TAG_Named.GetName() == "" {
		name = "TAG_String: "
	} else {
		name = "TAG_String(\"" + tag.TAG_Named.GetName() + "\"): "
	}
	string2 := []string{getPrefix(prefix), name, string(tag.Array), ": [", strconv.Itoa(len(tag.Array)), " bytes]\n"}
	output := strings.Join(string2, "")
	return output
}

func (tag TAG_String) GetString() string {
	return string(tag.Array)
}

func (tag TAG_String) Write(writer io.WriteCloser) (err error) {
	binary.Write(writer, binary.BigEndian, TAG_STRING)
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

func (tag TAG_String) WritePayload(writer io.WriteCloser) (err error) {
	err = binary.Write(writer, binary.BigEndian, int16(tag.Length))
	if err != nil {
		return err
	}
	err = binary.Write(writer, binary.BigEndian, []byte(tag.Array))
	return err
}
