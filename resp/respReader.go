package resp

import (
	"bufio"
	"io"
	"strconv"
)

const (
	Integer        = ':'
	String         = '+'
	Bulk_String    = '$'
	Array          = '*'
	Error          = '-'
	IntegerType    = "Integer"
	StringType     = "String"
	BulkStringType = "Bulk"
	ArrayType      = "Array"
	NULLType       = "NULL"
	ErrorType      = "Error"
)

type Value struct {
	Err   string
	Typ   string
	Str   string
	Num   int64
	Bulk  string
	Array []Value
}

type Resp struct {
	Reader *bufio.Reader
}

func NewRespHandler(rd io.Reader) *Resp {
	return &Resp{Reader: bufio.NewReader(rd)}
}

func (r *Resp) readUntilCLRF() (line []byte, n int, err error) {
	for {
		b, err := r.Reader.ReadByte()
		// fmt.Println(string(b))
		if err != nil {
			return nil, 0, err
		}
		n += 1
		line = append(line, b)
		// fmt.Println(len(line))
		if len(line) >= 2 && line[len(line)-2] == '\r' {
			// fmt.Printf("%v", string(b))
			break
		}

	}
	return line[:len(line)-2], n, err
}

func (r *Resp) ParseAny() (Value, error) {
	Typ, err := r.Reader.ReadByte()
	// fmt.Println(Typ)
	if err != nil {
		return Value{}, err
	}

	switch Typ {
	case Array:
		// fmt.Println("Array It is")
		return r.parseArrays()
	case Bulk_String:
		// fmt.Println("Bulk it is")
		return r.parseBulk()
	case Integer:
		// fmt.Println("Integer it is")
		return r.parseInteger()
	case String:
		// fmt.Println("String it is")
		return r.parseString()
	case Error:
		return r.parseError()
	default:
		// p, err := r.Reader.Peek(1)
		// fmt.Print("Problem: ", p[0])
		return Value{}, err
	}

}

func (r *Resp) parseArrays() (Value, error) {
	v := Value{
		Typ: ArrayType,
	}
	line, _, err := r.readUntilCLRF()
	// fmt.Println(string(line))
	if err != nil {
		return v, err
	}

	length, err := strconv.ParseInt(string(line[0]), 10, 32)
	// fmt.Println("length: %d", length)
	if err != nil {
		return v, err
	}
	v.Array = make([]Value, 0)
	for i := 0; i < int(length); i++ {
		val, err := r.ParseAny()
		// fmt.Println("Parsed in the Loop")
		if err != nil {
			return v, err
		}
		v.Array = append(v.Array, val)
	}
	// fmt.Printf("%v", v.Typ)
	return v, nil
}

func (r *Resp) parseBulk() (Value, error) {
	// fmt.Println("Reached at Bulk")
	v := Value{
		Typ: BulkStringType,
	}
	line, _, err := r.readUntilCLRF()
	if string(line) == "-1" {
		v.Bulk = string(line)
		return v, nil
	}
	if err != nil {
		return v, err
	}
	length, err := strconv.ParseInt(string(line[0]), 10, 32)
	if err != nil {
		return v, err
	}

	Bulk := make([]byte, length)
	r.Reader.Read(Bulk)
	v.Bulk = string(Bulk)
	r.readUntilCLRF()
	return v, nil
}

func (r *Resp) parseInteger() (Value, error) {
	v := Value{
		Typ: IntegerType,
	}
	line, _, err := r.readUntilCLRF()

	if err != nil {
		return v, err
	}

	v.Num, err = strconv.ParseInt(string(line), 10, 64)

	if err != nil {
		return v, err
	}
	r.readUntilCLRF()
	return v, nil

}

func (r *Resp) parseString() (Value, error) {
	v := Value{
		Typ: StringType,
	}
	line, _, err := r.readUntilCLRF()

	if err != nil {
		return v, err
	}
	v.Str = string(line)
	// fmt.Println("String Parsed")
	return v, nil
}

func (r *Resp) parseError() (Value, error) {
	v := Value{
		Typ: ErrorType,
	}

	line, _, err := r.readUntilCLRF()

	if err != nil {
		return v, err
	}

	v.Err = string(line)
	return v, nil
}
