package goNBT

import (
	"encoding/binary"
	"io"
	"strconv"
	"strings"
)

type TAG_Long struct {
	*TAG_Named
	Payload int64
}

func (tag TAG_Long) string(prefix int) string {
	var name string
	if tag.TAG_Named.GetName() == "" {
		name = "TAG_Long: "
	} else {
		name = "TAG_Long(\"" + tag.TAG_Named.GetName() + "\"): "
	}
	string2 := []string{getPrefix(prefix), name, strconv.FormatInt(int64(tag.Payload), 10), "\n"}
	output := strings.Join(string2, "")
	return output
}

func (tag TAG_Long) Write(writer io.WriteCloser) (err error) {
	binary.Write(writer, binary.BigEndian, TAG_LONG)
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

func (tag TAG_Long) WritePayload(writer io.WriteCloser) (err error) {
	err = binary.Write(writer, binary.BigEndian, tag.Payload)
	return err
}
