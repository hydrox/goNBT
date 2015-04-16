package goNBT

import (
	"encoding/binary"
	"io"
	"strconv"
	"strings"
)

type TAG_Byte struct {
	*TAG_Named
	Payload byte
}

func (tag TAG_Byte) string(prefix int) string {
	var name string
	if tag.TAG_Named.GetName() == "" {
		name = "TAG_Byte: "
	} else {
		name = "TAG_Byte(\"" + tag.TAG_Named.GetName() + "\"): "
	}
	string2 := []string{getPrefix(prefix), name, strconv.FormatUint(uint64(tag.Payload), 10), "\n"}
	output := strings.Join(string2, "")
	return output
}

func (tag TAG_Byte) Write(writer io.WriteCloser) (err error) {
	binary.Write(writer, binary.BigEndian, TAG_BYTE)
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
	err = tag.WritePayload(writer)
	return err
}

func (tag TAG_Byte) WritePayload(writer io.WriteCloser) (err error) {
	err = binary.Write(writer, binary.BigEndian, tag.Payload)
	return err
}
