package goNBT_test

import (
	"bytes"
	. "github.com/hydrox/goNBT"
	"io/ioutil"
	"testing"
)

const (
	SHORT_TEST  int16   = 32767
	LONG_TEST   int64   = 9223372036854775807
	FLOAT_TEST  float32 = 0.49823147
	DOUBLE_TEST float64 = 0.4931287132182315
	INT_TEST    int32   = 2147483647
	BYTE_TEST   byte    = 127
	STRING_TEST string  = "HELLO WORLD THIS IS A TEST STRING ÅÄÖ!"
)

func TestParseCompressed(t *testing.T) {
	content, err := ioutil.ReadFile("testdata/bigtest.nbt")
	if err != nil {
		//fmt.Fprintf(os.Stderr, "cat: can't open file: error %s\n", err)
		//os.Exit(1)
		t.Errorf("Can't open testdata/bigtest.nbt: error %s\n", err)
		t.Error(err)
		return
	}
	_, fail := ParseCompressed(content, GZIP)
	if fail {
		t.Errorf("Error while parsing testdata/bigtest.nbt")
	}

}

func TestGetEntry(t *testing.T) {
	LIST_LONG_TEST := [5]int64{11, 12, 13, 14, 15}

	content, err := ioutil.ReadFile("testdata/bigtest.nbt")
	if err != nil {
		t.Errorf("Can't open testdata/bigtest.nbt: error %s\n", err)
		return
	}
	compound, fail := ParseCompressed(content, GZIP)
	if fail {
		t.Errorf("Error while parsing testdata/bigtest.nbt")
	}

	entryShort, success := compound.GetShort("shortTest")
	if !success {
		t.Error("Failed to get entry: shortTest")
	} else if entryShort.Payload != SHORT_TEST {
		t.Errorf("GetShort(\"shortTest\") == <%d> want <%d>", entryShort.Payload, SHORT_TEST)
	}

	entryByte, success := compound.GetByte("byteTest")
	if !success {
		t.Error("Failed to get entry: byteTest")
	} else if entryByte.Payload != BYTE_TEST {
		t.Errorf("GetByte(\"byteTest\") == <%d> want <%d>", entryByte.Payload, BYTE_TEST)
	}

	entryLong, success := compound.GetLong("longTest")
	if !success {
		t.Error("Failed to get entry: longTest")
	} else if entryLong.Payload != LONG_TEST {
		t.Errorf("GetLong(\"longTest\") == <%d> want <%d>", entryLong.Payload, LONG_TEST)
	}

	entryFloat, success := compound.GetFloat("floatTest")
	if !success {
		t.Error("Failed to get entry: floatTest")
	} else if entryFloat.Payload != FLOAT_TEST {
		t.Errorf("GetFloat(\"floatTest\") == <%d> want <%d>", entryFloat.Payload, FLOAT_TEST)
	}

	entryDouble, success := compound.GetDouble("doubleTest")
	if !success {
		t.Error("Failed to get entry: doubleTest")
	} else if entryDouble.Payload != DOUBLE_TEST {
		t.Errorf("GetDouble(\"doubleTest\") == <%d> want <%d>", entryDouble.Payload, DOUBLE_TEST)
	}

	entryString, success := compound.GetString("stringTest")
	if !success {
		t.Error("Failed to get entry: stringTest")
	} else if entryString.GetString() != STRING_TEST {
		t.Errorf("GetString(\"stringTest\") == <%s> want <%s>", entryString.GetString(), STRING_TEST)
	}

	entryInt, success := compound.GetInt("intTest")
	if !success {
		t.Error("Failed to get entry: intTest")
	} else if entryInt.Payload != INT_TEST {
		t.Errorf("GetInt(\"intTest\") == <%d> want <%d>", entryInt.Payload, INT_TEST)
	}

	entryLongList, success := compound.GetList("listTest (long)")
	if !success {
		t.Error("Failed to get entry: listTest (long)")
	} else {
		longArray, success := entryLongList.GetLongEntries()
		if !success {
			t.Error("Failed to get Long-Entries")
		} else {
			for i := int32(0); i < entryLongList.Length; i++ {
				if longArray[i].Payload != LIST_LONG_TEST[i] {
					t.Errorf("Entry i from \"listTest (long)\" == <%d> want <%d>", longArray[i].Payload, LIST_LONG_TEST[i])
				}
			}
		}
	}

	entryCompound, success := compound.GetCompound("nested compound test")
	if !success {
		t.Error("Failed to get entry: nested compound test")
	} else {
		entryCompound2, success := entryCompound.GetCompound("ham")
		if !success {
			t.Error("Failed to get entry: ham")
		} else {
			entryCompoundString, success := entryCompound2.GetString("name")
			if !success {
				t.Error("Failed to get entry: name")
			} else if entryCompoundString.GetString() != "Hampus" {
				t.Errorf("GetString(\"name\") == <%s> want <%s>", entryCompoundString.GetString(), "Hampus")
			}
			entryCompoundFloat, success := entryCompound2.GetFloat("value")
			if !success {
				t.Error("Failed to get entry: value")
			} else if entryCompoundFloat.Payload != 0.75 {
				t.Errorf("GetFloat(\"value\") == <%d> want <%d>", entryCompoundFloat.Payload, 0.75)
			}
		}
		longArray, success := entryLongList.GetLongEntries()
		if !success {
			t.Error("Failed to get Long-Entries")
		} else {
			for i := int32(0); i < entryLongList.Length; i++ {
				if longArray[i].Payload != LIST_LONG_TEST[i] {
					t.Errorf("Entry i from \"listTest (long)\" == <%d> want <%d>", longArray[i].Payload, LIST_LONG_TEST[i])
				}
			}
		}
	}

	entryByteArray, success := compound.GetByteArray("byteArrayTest (the first 1000 values of (n*n*255+n*7)%100, starting with n=0 (0, 62, 34, 16, 8, ...))")
	if !success {
		t.Error("Failed to get entry: byteArray")
	} else {
		for i := uint32(0); i < entryByteArray.Length; i++ {
			if entryByteArray.Array[i] != byte((i*i*255+i*7)%100) {
				t.Errorf("Entry i from \"byteArray\" == <%d> want <%d>", entryByteArray.Array[i], (i*i*255+i*7)%100)
			}
		}
	}
}

func TestWriteNBT(t *testing.T) {
	content, err := ioutil.ReadFile("testdata/bigtest.nbt")
	if err != nil {
		t.Errorf("Can't open testdata/bigtest.nbt: error %s\n", err)
		return
	}
	compound, fail := ParseCompressed(content, GZIP)
	if fail {
		t.Errorf("Error while parsing testdata/bigtest.nbt")
	}

	writer := NewByteWriter()
	err = WriteCompressedStream(writer, compound, GZIP, 0)
	if err != nil {
		t.Errorf("Error while writing NBT-data: {0}", err)
	}
	compound2, fail := ParseCompressed(*writer.Content, GZIP)

	ioutil.WriteFile("testdata/bigtestOutput.nbt", *writer.Content, 0777)
	ioutil.WriteFile("testdata/compound1.txt", []byte(compound.String()), 0777)
	ioutil.WriteFile("testdata/compound2.txt", []byte(compound2.String()), 0777)
	if fail {
		t.Errorf("Error while parsing written NBT-data")
	} else if compound.String() != compound2.String() {
		t.Error("compound != compound2")
	}
}

func BenchmarkParse(b *testing.B) {
	content, err := ioutil.ReadFile("testdata/bigtest.nbt")
	if err != nil {
		//fmt.Fprintf(os.Stderr, "cat: can't open file: error %s\n", err)
		//os.Exit(1)
		b.Errorf("Can't open testdata/bigtest.nbt")
		b.Error(err)
		return
	}

	compound, fail := ParseCompressed(content, GZIP)
	if fail {
		b.Errorf("Error while parsing testdata/bigtest.nbt")
	}

	writer := NewByteWriter()
	Write(writer, compound)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var buffer = bytes.NewBuffer(*writer.Content)
			Parse(buffer)
		}
	})
}

func BenchmarkParseCompressed(b *testing.B) {
	content, err := ioutil.ReadFile("testdata/bigtest.nbt")
	if err != nil {
		//fmt.Fprintf(os.Stderr, "cat: can't open file: error %s\n", err)
		//os.Exit(1)
		b.Errorf("Can't open testdata/bigtest.nbt")
		b.Error(err)
		return
	}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, fail := ParseCompressed(content, GZIP)
			if fail {
				b.Errorf("Error while parsing testdata/bigtest.nbt")
			}

		}
	})
}

func BenchmarkWrite(b *testing.B) {
	content, err := ioutil.ReadFile("testdata/bigtest.nbt")
	if err != nil {
		//fmt.Fprintf(os.Stderr, "cat: can't open file: error %s\n", err)
		//os.Exit(1)
		b.Errorf("Can't open testdata/bigtest.nbt")
		b.Error(err)
		return
	}

	compound, fail := ParseCompressed(content, GZIP)
	if fail {
		b.Errorf("Error while parsing testdata/bigtest.nbt")
	}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			writer := NewByteWriter()
			Write(writer, compound)
		}
	})
}

func BenchmarkWriteCompressed(b *testing.B) {
	content, err := ioutil.ReadFile("testdata/bigtest.nbt")
	if err != nil {
		//fmt.Fprintf(os.Stderr, "cat: can't open file: error %s\n", err)
		//os.Exit(1)
		b.Errorf("Can't open testdata/bigtest.nbt")
		b.Error(err)
		return
	}

	compound, fail := ParseCompressed(content, GZIP)
	if fail {
		b.Errorf("Error while parsing testdata/bigtest.nbt")
	}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			writer := NewByteWriter()
			WriteCompressedStream(writer, compound, GZIP, 9)
		}
	})
}

/*func BenchmarkPrintNBT(b *testing.B) {
	content, err := ioutil.ReadFile("testdata/bigtest.nbt")
	if err != nil {
		//fmt.Fprintf(os.Stderr, "cat: can't open file: error %s\n", err)
		//os.Exit(1)
		b.Errorf("Can't open testdata/bigtest.nbt")
		b.Error(err)
		return
	}
	compound, fail := ParseCompressed(content, GZIP)
	if fail {
		b.Errorf("Error while parsing testdata/bigtest.nbt")
	}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			compound.String()
		}
	})
}*/

type ByteWriter struct {
	Content *[]byte
	pointer int
}

func NewByteWriter() (newWriter *ByteWriter) {
	array := make([]byte, 0)
	newWriter = &ByteWriter{
		Content: &array,
		pointer: 0,
	}
	return newWriter
}

func (b *ByteWriter) Write(p []byte) (n int, err error) {
	*b.Content = append(*b.Content, p...)
	return len(p), nil
}

func (b *ByteWriter) Close() (err error) {
	*b.Content = nil
	return nil
}
