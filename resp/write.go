package resp

import (
	"fmt"
	"io"
	"strconv"
)

type Writer struct {
	writer io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{writer: w}
}

func (v Value) Marshal() []byte {
	switch v.typ {
	case Array:
		return v.marshalArray()
	case Bulk:
		return v.marshalBulk()
	case String:
		return v.marshalString()
	case Null:
		return v.marshallNull()
	case Error:
		return v.marshallError()
	default:
		return []byte{}
	}
}

func (v Value) marshalString() []byte {
	str := fmt.Sprintf("%s%s\r\n", STRING, v.str)
	return []byte(str)
}

func (v Value) marshalBulk() []byte {
	bulkLen := len(v.bulk)
	str := fmt.Sprintf("%s%s\r\n%s\r\n", BULK, strconv.Itoa(bulkLen), v.bulk)
	return []byte(str)
}

func (v Value) marshalArray() []byte {
	arrLen := len(v.array)
	str := fmt.Sprintf("%s%s\r\n", ARRAY, strconv.Itoa(arrLen))
	bytes := []byte(str)
	for i := 0; i < arrLen; i++ {
		bytes = append(bytes, v.array[i].Marshal()...)
	}
	return bytes
}

func (v Value) marshallError() []byte {
	bytes := fmt.Sprintf("%s%s\r\n", ERROR, v.str)
	return []byte(bytes)
}

func (v Value) marshallNull() []byte {
	return []byte("$-1\r\n")
}
