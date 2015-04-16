package goNBT

import (
	"compress/flate"
	"compress/gzip"
	"compress/lzw"
	"compress/zlib"
	"fmt"
	"io"
	"os"
)

func WriteCompressedStream(compressedWriter io.Writer, tag TAG_Writeable, compression byte, compressionRate int) (err error) {
	var writer io.WriteCloser
	switch compression {
	case FLATE:
		writer, err = flate.NewWriter(compressedWriter, compressionRate)
	case GZIP:
		writer, err = gzip.NewWriterLevel(compressedWriter, compressionRate)
	case ZLIB:
		writer, err = zlib.NewWriterLevel(compressedWriter, compressionRate)
	case LZW:
		writer = lzw.NewWriter(compressedWriter, lzw.MSB, 8)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "error %s\n", err)
		return err
	}
	defer writer.Close()
	return Write(writer, tag)
}

func Write(writer io.WriteCloser, tag TAG_Writeable) (err error) {
	return tag.Write(writer)
}
