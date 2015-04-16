package goNBT

import (
	"encoding/binary"
	"io"
	"strconv"
	"strings"
)

type TAG_Float struct {
	*TAG_Named
	Payload float32
}

func (tag TAG_Float) string(prefix int) string {
	var name string
	if tag.TAG_Named.GetName() == "" {
		name = "TAG_Float: "
	} else {
		name = "TAG_Float(\"" + tag.TAG_Named.GetName() + "\"): "
	}
	string2 := []string{getPrefix(prefix), name, strconv.FormatFloat(float64(tag.Payload), 'G', -1, 32), "\n"}
	output := strings.Join(string2, "")
	return output
}

func (tag TAG_Float) Write(writer io.WriteCloser) (err error) {
	binary.Write(writer, binary.BigEndian, TAG_FLOAT)
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

func (tag TAG_Float) WritePayload(writer io.WriteCloser) (err error) {
	err = binary.Write(writer, binary.BigEndian, tag.Payload)
	return err
}
