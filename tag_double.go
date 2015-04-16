package goNBT

import (
	"encoding/binary"
	"io"
	"strconv"
	"strings"
)

type TAG_Double struct {
	*TAG_Named
	Payload float64
}

func (tag TAG_Double) string(prefix int) string {
	var name string
	if tag.TAG_Named.GetName() == "" {
		name = "TAG_Double: "
	} else {
		name = "TAG_Double(\"" + tag.TAG_Named.GetName() + "\"): "
	}
	string2 := []string{getPrefix(prefix), name, strconv.FormatFloat(float64(tag.Payload), 'G', -1, 64), "\n"}
	output := strings.Join(string2, "")
	return output
}

func (tag TAG_Double) Write(writer io.WriteCloser) (err error) {
	binary.Write(writer, binary.BigEndian, TAG_DOUBLE)
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

func (tag TAG_Double) WritePayload(writer io.WriteCloser) (err error) {
	err = binary.Write(writer, binary.BigEndian, tag.Payload)
	return err
}
