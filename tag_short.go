package goNBT

import (
	"encoding/binary"
	"io"
	"strconv"
	"strings"
)

type TAG_Short struct {
	*TAG_Named
	Payload int16
}

func (tag TAG_Short) string(prefix int) string {
	var name string
	if tag.TAG_Named.GetName() == "" {
		name = "TAG_Short: "
	} else {
		name = "TAG_Short(\"" + tag.TAG_Named.GetName() + "\"): "
	}
	string2 := []string{getPrefix(prefix), name, strconv.Itoa(int(tag.Payload)), "\n"}
	output := strings.Join(string2, "")
	return output
}

func (tag TAG_Short) Write(writer io.WriteCloser) (err error) {
	binary.Write(writer, binary.BigEndian, TAG_SHORT)
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

func (tag TAG_Short) WritePayload(writer io.WriteCloser) (err error) {
	err = binary.Write(writer, binary.BigEndian, tag.Payload)
	return err
}
