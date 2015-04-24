package goNBT

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"compress/lzw"
	"compress/zlib"
	"container/list"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

func ParseCompressed(data []byte, compression byte) (compound TAG_Compound, ok bool) {
	var buffer = bytes.NewBuffer(data)
	return ParseCompressedStream(buffer, compression)
}

func ParseCompressedStream(compressedReader io.Reader, compression byte) (compound TAG_Compound, ok bool) {
	var reader io.ReadCloser
	var okz error
	switch compression {
	case FLATE:
		reader = flate.NewReader(compressedReader)
	case GZIP:
		reader, okz = gzip.NewReader(compressedReader)
	case ZLIB:
		reader, okz = zlib.NewReader(compressedReader)
	case LZW:
		reader = lzw.NewReader(compressedReader, lzw.MSB, 8)
	default:
		fmt.Fprintf(os.Stderr, "invalid compression %d", compression)
	}
	if okz != nil {
		if debug {
			fmt.Fprintf(os.Stderr, "error %s\n", okz)
		}
		return compound, false
	}

	defer reader.Close()
	return Parse(reader)
}

func Parse(reader io.Reader) (compound TAG_Compound, fail bool) {
	var realreader io.Reader = reader

	var tagid [1]byte
	var err error
	var n int
	n, err = realreader.Read(tagid[:])
	var length uint16
	var tmp *uint16 = &length
	binary.Read(realreader, binary.BigEndian, tmp)
	if debug {
		fmt.Printf("length :%d\n", length)
	}
	var name string
	tmp2 := make([]byte, length)
	tmp3 := &tmp2
	binary.Read(realreader, binary.BigEndian, tmp3)
	name = string(tmp2)
	if debug {
		fmt.Printf("name :%s\n", name)
	}
	switch tagid[0] {
	case TAG_COMPOUND:
		compound, fail = parseCompoundTag(&realreader)
	default:
		fmt.Printf("Unexpected tagid: %d n: %d\n", tagid[0], n)
		return compound, false
	}
	if err != nil {
		return
	}
	if n == 0 {
		return
	}
	compound.TAG_Named.Name = name
	return compound, fail
}

func parseCompoundTag(reader *io.Reader) (tagcompound TAG_Compound, fail bool) {
	var realreader io.Reader = *reader
	if debug {
		fmt.Println("Parsing Compound")
	}
	tagcompound = TAG_Compound{
		TAG_Named: &TAG_Named{""},
		List:      list.New(),
	}
	tagcompound.List.Init()
	var tagid [1]byte
	var err error
	var n int
	var count int
	var tagnames string

	n, err = realreader.Read(tagid[:])
	for ; n >= 0; n, err = realreader.Read(tagid[:]) {
		count++
		switch tagid[0] {
		case TAG_END:
			if debug {
				fmt.Println("TAG_END: Close Compound")
			}
			return tagcompound, fail
		default:
			if debug {
				fmt.Printf("tagid: %d\n", tagid[0])
			}
			tag, fail := parseNamedTag(reader, tagid[0])
			if fail {
				fmt.Printf("Tagnames: %s\n", tagnames)
				fmt.Printf("FAILED to parse Compound-Entry %d\n", count)
				tagcompound.List.PushBack(tag)
				return tagcompound, fail
			}
			tagnames = fmt.Sprintf("%s %s", tagnames, tag.GetName())
			tagcompound.List.PushBack(tag)
		}
		if err != nil {
			return
		}
		if n == 0 {
			return
		}
	}
	return tagcompound, fail
}

func parseListTag(reader *io.Reader) (taglist TAG_List, fail bool) {
	var realreader io.Reader = *reader
	if debug {
		fmt.Println("Parsing List")
	}
	taglist = TAG_List{
		TAG_Named: &TAG_Named{""},
		List:      list.New(),
	}
	taglist.List.Init()
	var tagid [1]byte
	realreader.Read(tagid[:])
	taglist.TagId = tagid[0]
	if debug {
		fmt.Printf("tagId :%d\n", tagid[0])
	}
	var length int32
	var tmp *int32 = &length
	binary.Read(realreader, binary.BigEndian, tmp)
	if debug {
		fmt.Printf("length :%d\n", length)
	}
	taglist.Length = length
	var i int32 = 0
	for ; i < taglist.Length; i++ {
		if debug {
			fmt.Printf("Entry %d out of %d\n", i+1, taglist.Length)
		}
		tag, fail := parseUnnamedTag(reader, tagid[0])
		if fail {
			fmt.Printf("FAILED to parse Entry %d of List (expected Tag-Id %d)\n", i, tagid[0])
			taglist.List.PushBack(tag)
			return taglist, fail
		}
		taglist.List.PushBack(tag)
	}
	return taglist, fail
}

func parseUnnamedTag(reader *io.Reader, tagid byte) (tag TAG_Name, fail bool) {
	var realreader io.Reader = *reader
	switch tagid {
	case TAG_BYTE:
		if debug {
			fmt.Println("Parsing Byte")
		}
		tmp2 := make([]byte, 1)
		realreader.Read(tmp2)
		if debug {
			fmt.Printf("value :%b\n", tmp2[0])
		}
		tag = TAG_Byte{
			TAG_Named: &TAG_Named{""},
			Payload:   tmp2[0],
		}
	case TAG_SHORT:
		if debug {
			fmt.Println("Parsing Short")
		}
		number := ReadInt16(&realreader)
		if debug {
			fmt.Printf("value :%d\n", number)
		}
		tag = TAG_Short{
			TAG_Named: &TAG_Named{""},
			Payload:   number,
		}
	case TAG_INT:
		if debug {
			fmt.Println("Parsing Int")
		}
		number := ReadInt32(&realreader)
		if debug {
			fmt.Printf("value :%d\n", number)
		}
		tag = TAG_Int{
			TAG_Named: &TAG_Named{""},
			Payload:   number,
		}
	case TAG_LONG:
		if debug {
			fmt.Println("Parsing Long")
		}
		number := ReadInt64(&realreader)
		if debug {
			fmt.Printf("value :%d\n", number)
		}
		tag = TAG_Long{
			TAG_Named: &TAG_Named{""},
			Payload:   number,
		}
	case TAG_FLOAT:
		if debug {
			fmt.Println("Parsing Float")
		}
		tmp2 := make([]float32, 1)
		tmp3 := &tmp2
		binary.Read(realreader, binary.BigEndian, tmp3)
		if debug {
			fmt.Printf("value :%f\n", tmp2[0])
		}
		tag = TAG_Float{
			TAG_Named: &TAG_Named{""},
			Payload:   tmp2[0],
		}
	case TAG_DOUBLE:
		if debug {
			fmt.Println("Parsing Double")
		}
		tmp2 := make([]float64, 1)
		tmp3 := &tmp2
		binary.Read(realreader, binary.BigEndian, tmp3)
		if debug {
			fmt.Printf("value :%f\n", tmp2[0])
		}
		tag = TAG_Double{
			TAG_Named: &TAG_Named{""},
			Payload:   tmp2[0],
		}
	case TAG_BYTE_ARRAY:
		if debug {
			fmt.Println("Parsing Byte_Array")
		}
		var length uint32
		var tmp *uint32 = &length
		binary.Read(realreader, binary.BigEndian, tmp)
		if debug {
			fmt.Printf("length :%d\n", length)
		}
		tmp2 := ReadByteArray(&realreader, uint32(length))
		if debug {
			//			fmt.Printf("value :%s\n", tmp2)
		}
		tag = TAG_Byte_Array{
			TAG_Named: &TAG_Named{""},
			Length:    length,
			Array:     tmp2,
		}
	case TAG_STRING:
		if debug {
			fmt.Println("Parsing String")
		}
		var length uint16
		var tmp *uint16 = &length
		binary.Read(realreader, binary.BigEndian, tmp)
		if debug {
			fmt.Printf("length :%d\n", length)
		}
		name := ReadByteArray(&realreader, uint32(length))
		if debug {
			fmt.Printf("value :%s\n", string(name))
		}
		tag = TAG_String{
			TAG_Named: &TAG_Named{""},
			Length:    length,
			Array:     name,
		}
	case TAG_LIST:
		tag, fail = parseListTag(reader)
	case TAG_COMPOUND:
		tag, fail = parseCompoundTag(reader)
	case TAG_INT_ARRAY:
		if debug {
			fmt.Println("Parsing Int_Array")
		}
		var length uint32
		var tmp *uint32 = &length
		binary.Read(realreader, binary.BigEndian, tmp)
		if debug {
			fmt.Printf("length :%d\n", length)
		}
		intArray := make([]int32, length)
		for i := 0; i < int(length); i++ {
			intArray[i] = ReadInt32(&realreader)
		}
		if debug {
			//			fmt.Printf("value :%s\n", tmp2)
		}
		tag = TAG_Int_Array{
			TAG_Named: &TAG_Named{""},
			Length:    length,
			Array:     intArray,
		}
	default:
		fmt.Printf("tagID %d doesn't exist\n", tagid)
		return tag, true
	}
	return tag, fail
}

func parseNamedTag(reader *io.Reader, tagid byte) (tag TAG_Name, fail bool) {
	var realreader io.Reader = *reader
	if debug {
		fmt.Println("Parsing Name")
	}
	length := ReadInt16(&realreader)
	if debug {
		fmt.Printf("length :%d\n", length)
	}
	tmp2 := ReadByteArray(&realreader, uint32(length))
	name := string(tmp2)
	if debug {
		fmt.Printf("name :%s\n", name)
	}
	tag, fail = parseUnnamedTag(reader, tagid)
	if fail {
		if len(name) >= 50 {
			fmt.Printf("FAILED to parse Tag with Name DERPED (length: %d) and Tag-ID %d\n", len(name), tagid)
		} else {
			fmt.Printf("FAILED to parse Tag with Name %s and Tag-ID %d\n", name, tagid)
		}
	} else {
		tag.SetName(name)
	}
	return tag, fail
}
