package goNBT

import (
	"fmt"
	//	"os"
	"io"
	"strings"
)

const (
	TAG_END byte = iota
	TAG_BYTE
	TAG_SHORT
	TAG_INT
	TAG_LONG
	TAG_FLOAT
	TAG_DOUBLE
	TAG_BYTE_ARRAY
	TAG_STRING
	TAG_LIST
	TAG_COMPOUND
	TAG_INT_ARRAY
)

const (
	NONE   byte   = 0
	GZIP   byte   = 1
	ZLIB   byte   = 2
	FLATE  byte   = 3
	LZW    byte   = 128 //not officially supported. Do not use
	PREFIX string = "  "
)

var debug bool = false

func getPrefix(indent int) string {
	var output string = ""
	for i := 0; i < indent; i++ {
		string2 := []string{output, PREFIX}
		output = strings.Join(string2, "")
	}
	return output
}

type TAG_Unnamed interface {
}

type TAG_Printable interface {
	string(prefix int) string
}
type TAG_Name interface {
	GetName() string
	SetName(name string)
}

type TAG_Writeable interface {
	Write(writer io.WriteCloser) (err error)
	WritePayload(writer io.WriteCloser) (err error)
}

type TAG_End byte

func toString(value interface{}, prefix int) string {
	switch v := value.(type) {
	case TAG_Printable:
		return v.string(prefix)
	}
	return ""
}

func ReadInt64(reader *io.Reader) (result int64) {
	var realreader io.Reader = *reader
	tmp := make([]byte, 8)
	realreader.Read(tmp)
	for i := 0; i < 8; i++ {
		result |= int64(tmp[i]) << (8 * uint(7-i))
	}
	return result
}

func ReadInt32(reader *io.Reader) (result int32) {
	var realreader io.Reader = *reader
	tmp := make([]byte, 4)
	realreader.Read(tmp)
	for i := 0; i < 4; i++ {
		result |= int32(tmp[i]) << (8 * uint(3-i))
	}
	return result
}

func ReadInt16(reader *io.Reader) (result int16) {
	var realreader io.Reader = *reader
	tmp := make([]byte, 2)
	realreader.Read(tmp)
	for i := 0; i < 2; i++ {
		result |= int16(tmp[i]) << (8 * uint(1-i))
	}
	return result
}

func ReadByteArray(reader *io.Reader, length uint32) (result []byte) {
	var realreader io.Reader = *reader
	if length > 32768 {
		fmt.Printf("length: %d but should be %d at max\n", length, 32768)
		return make([]byte, 0)
	}
	tmp := make([]byte, length)
	n, err := realreader.Read(tmp)
	if err != nil {
		fmt.Println(err)
		return make([]byte, 0)
		//		os.Exit(-1)
	}
	if n < int(length) {
		if debug {
			fmt.Printf("only %d out of %d bytes read, trying to read again ", n, length)
		}
		m, _ := realreader.Read(tmp[n:])
		if debug && n+m == int(length) {
			fmt.Printf("OK now\n")
		}
	}
	result = tmp
	return result
}
