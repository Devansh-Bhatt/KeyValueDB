package resp

import (
	"fmt"
	"io"
	"strconv"
)

type Writer struct {
	writer io.Writer
}

func NewRespWriter(w io.Writer) *Writer {
	return &Writer{writer: w}
}

func (w *Writer) Write(v Value) error {
	var bytes = v.MarshalAny()

	_, err := w.writer.Write(bytes)

	if err != nil {
		fmt.Printf("Err from Here : %s", err)
		// return err
	}

	return nil
}

func (v Value) MarshalAny() []byte {
	Typ := v.Typ

	switch Typ {
	case StringType:
		return v.MarshalString()
	case IntegerType:
		return v.MarshalInteger()
	case BulkStringType:
		// fmt.Println("From Bulk")
		return v.MarshalBulk()
	case ArrayType:
		return v.MarshalArray()
	case NULLType:
		return v.MarshalNull()
	case ErrorType:
		return v.MarshalError()
	case RDBType:
		return v.MarshalRDB()
	default:
		return []byte{}
	}
}

func (v Value) MarshalString() []byte {
	var bytes []byte
	bytes = append(bytes, String)
	bytes = append(bytes, v.Str...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

func (v Value) MarshalNull() []byte {
	return []byte("$-1\r\n")
}

func (v Value) MarshalBulk() []byte {
	// fmt.Println("In bulk")
	var bytes []byte

	bytes = append(bytes, Bulk_String)
	bytes = append(bytes, strconv.Itoa(len(v.Bulk))...)
	bytes = append(bytes, '\r', '\n')
	bytes = append(bytes, v.Bulk...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

func (v Value) MarshalArray() []byte {
	var bytes []byte

	bytes = append(bytes, Array)
	bytes = append(bytes, strconv.Itoa(len(v.Array))...)
	bytes = append(bytes, '\r', '\n')
	for i := 0; i < len(v.Array); i++ {
		bytes = append(bytes, v.Array[i].MarshalAny()...)
	}
	return bytes
}

func (v Value) MarshalInteger() []byte {
	var bytes []byte

	bytes = append(bytes, Integer)
	bytes = append(bytes, byte(v.Num))
	bytes = append(bytes, '\r', '\n')
	return bytes
}

func (v Value) MarshalError() []byte {
	var bytes []byte
	bytes = append(bytes, Error)
	bytes = append(bytes, v.Err...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

func (v Value) MarshalRDB() []byte {
	var bytes []byte
	bytes = append(bytes, Bulk_String)
	bytes = append(bytes, strconv.Itoa(len(v.Bytes))...)
	bytes = append(bytes, '\r', '\n')
	bytes = append(bytes, v.Bytes...)
	return bytes
}
