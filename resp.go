package main

import (
	"bufio"
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

// $5\r\nAhmed\r\n"
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
