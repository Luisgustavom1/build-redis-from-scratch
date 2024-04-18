package resp

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

const (
	STRING  = "+"
	ERROR   = "-"
	INTEGER = ":"
	BULK    = "$"
	ARRAY   = "*"
)

type ValueType uint8

const (
	Array ValueType = iota
	Bulk
	String
	Null
	Error
)

func (v ValueType) String() string {
	switch v {
	case Array:
		return "Array"
	case Bulk:
		return "Bulk"
	case String:
		return "String"
	case Null:
		return "Null"
	case Error:
		return "Error"
	default:
		return "Unknown"
	}
}

type Value struct {
	Typ   ValueType
	Str   string
	Num   int
	Bulk  string
	Array []Value
}

type Resp struct {
	reader *bufio.Reader
}

func NewResp(rd io.Reader) *Resp {
	return &Resp{reader: bufio.NewReader(rd)}
}

func (r *Resp) readInteger() (x int, err error) {
	line, _, err := r.reader.ReadLine()
	if err != nil {
		return 0, err
	}
	i64, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, err
	}
	return int(i64), nil
}

func (r *Resp) Read() (Value, error) {
	t, err := r.reader.ReadByte()
	if err != nil {
		return Value{}, err
	}

	switch string(t) {
	case ARRAY:
		return r.readArray()
	case BULK:
		return r.readBulk()
	default:
		fmt.Printf("Unknown type: %v", string(t))
		return Value{}, nil
	}
}

func (r *Resp) readArray() (Value, error) {
	v := Value{}
	v.Typ = Array

	arrLen, err := r.readInteger()
	if err != nil {
		return v, err
	}

	v.Array = make([]Value, arrLen)
	for i := 0; i < arrLen; i++ {
		val, err := r.Read()
		if err != nil {
			return v, err
		}
		v.Array[i] = val
	}

	return v, nil
}

func (r *Resp) readBulk() (Value, error) {
	v := Value{}
	v.Typ = Bulk

	bulkLen, err := r.readInteger()
	if err != nil {
		return v, err
	}

	bulkBytes := make([]byte, bulkLen)
	r.reader.Read(bulkBytes)
	v.Bulk = string(bulkBytes)

	// Read the trailing \r\n
	r.reader.ReadLine()

	return v, nil
}
