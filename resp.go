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
