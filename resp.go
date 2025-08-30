package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strconv"
)

const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

type Value struct {
	typ   string
	str   string
	num   int
	bulk  string
	array []Value
}

func (v Value) Marshal() []byte {
	switch v.typ {
	case "string":
		return v.MarshalString()
	case "bulk":
		return v.MarshalBulk()
	case "array":
		return v.MarshalArray()
	case "error":
		return v.MarshalError()
	case "null":
		return v.MarshalNull()
	default:
		return []byte{}
		
	}
}

func (v Value) MarshalString () []byte { 
	var byte []byte
	byte = append(byte,STRING)
	byte = append(byte,v.str...)
	byte = append(byte, '\r', '\n')
	return byte
}

func (v Value) MarshalBulk() []byte {
	var byte []byte
	byte = append(byte, BULK)
	byte = append(byte, strconv.Itoa(len(v.bulk))...)
	byte = append(byte, '\r', '\n')
	byte = append(byte, v.bulk...)
	byte = append(byte, '\r', '\n')
	return byte
}

func (v Value) MarshalArray() []byte {
	len := len(v.array)
	var byte []byte
	byte = append(byte, ARRAY)
	byte = append(byte, strconv.Itoa(len)...)
	byte = append(byte, '\r', '\n')
	for i := 0; i < len; i++ {
		byte = append(byte, v.array[i].Marshal()...)		
	}
	return byte
}

func (v Value) MarshalError() []byte {
	var byte []byte
	byte = append(byte, ERROR)
	byte = append(byte, v.str...)
	byte = append(byte, '\r', '\n')
	return byte
}

func (v Value) MarshalNull() []byte {
	return []byte("$-1\r]\n")
}

type Resp struct {
	reader bufio.Reader
}

func NewResp(rd io.Reader) *Resp {
	return &Resp{reader: *bufio.NewReader(rd)}
}

func (r *Resp) readLine() (line []byte, n int, err error) {
	for {
		b, err := r.reader.ReadByte()
		if err != nil {
			return nil, 0, err
		}
		n += 1
		line = append(line, b)
		if len(line) >= 2 && line[len(line)-2] == '\r' {
			break
		}
	}
	return line[:len(line)-2], n, nil
}

func (r *Resp) readInteger() (x int, n int, err error) {
	line, n, err := r.readLine()
	if err != nil {
		return 0, 0, err
	}
	i64, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, n, err
	}
	return int(i64), n, nil
}
func (r *Resp) Read() (Value, error) {
	_type, err := r.reader.ReadByte()
	if err != nil {
		log.Fatal("error reading type")
	}

	switch _type {
	case ARRAY:
		return r.ReadArray()
	case BULK:
		return r.ReadBulk()
	default:
		fmt.Printf("unknown type %c\n", _type)
		return Value{}, err
	}
}

func (r *Resp) ReadArray() (Value, error) {
	v := Value{}
	v.typ = "array"

	len, _, err := r.readInteger()
	if err != nil {
		log.Fatal("error reading length")
	}

	for range len {
		val, err := r.Read()
		if err != nil {
			log.Fatal("error reading value")
		}
		v.array = append(v.array, val)
	}
	return v, err
}

func (r *Resp) ReadBulk() (Value, error) {
	v := Value{}
	v.typ = "bulk"

	len, _, err := r.readInteger()
	if err != nil {
		log.Fatal("error reading length")
	}
	bulk := make([]byte, len)
	r.reader.Read(bulk)
	v.bulk = string(bulk)
	r.readLine()

	return v, nil
}

type writer struct {
	writter io.Writer
}

func newWritter (w io.Writer) *writer {
	return &writer{writter:w}
}

func (w *writer) Write(v Value) error {
	byte := v.Marshal()
	_, err := w.writter.Write(byte)
	if err != nil {
		return err
	}
	return nil
}	