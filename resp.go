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

func newResp(rd io.Reader) *Resp {
	return &Resp{reader: *bufio.NewReader(rd)}
}

func (r *Resp) readLine() (line []byte, n int, err error) {

	for {
		b, err := r.reader.ReadByte()
		if err != nil {
			log.Fatal("error reading byte")
		}
		n = len(line)
		line = append(line, b)
		len := len(line)
		var sep = len - 2
		if string(line[sep]) == "\r" {
			break
		}
	}

	return line[:len(line)-2], n, err
}

func (r *Resp) readInteger() (x int,err error) {
	line, _, err := r.readLine()
	if err != nil {
		log.Fatal("Error reading line")
	}
	i,err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		log.Fatal("error in parsing",err)
	}
	return int(i),err
}


func (r *Resp) Read() (Value,error) {
	_type,err := r.reader.ReadByte()
	if err != nil {
		log.Fatal("error reading type")
	}

	switch _type {
		case ARRAY:
			return Value{},err
		case ERROR:
			return Value{},err
		case BULK:
			return Value{},err
		default:
			fmt.Printf("unknown type %c\n",_type)
			return Value{},err
	}
}