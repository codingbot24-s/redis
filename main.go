package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func main() {
	fmt.Println("Listening on port :6379")
	// create a server
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		log.Fatal("port binding error", err)
		return
	}
	// listen for a connection
	conn, err := l.Accept()
	if err != nil {
		log.Fatal("connot accept the connection", err)
		return
	}
	defer conn.Close()

	for {
		resp := NewResp(conn)
		val, err := resp.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal("error in reading", err)
			return
		}

		if val.typ != "array" {
			fmt.Println("Invalid request, expected array")
			continue
		}
		if len(val.array) == 0 {
			fmt.Println("Invalid request, expected array length > 0")
			continue
		}
		cmd := strings.ToUpper(val.array[0].bulk)
		args := val.array[1:]
		writer := newWritter(conn)
		handler, ok := Handlers[cmd]
		if !ok {
			fmt.Println("Invalid command")
			writer.Write(Value{typ: "error", str: "invalid command"})
			continue
		}

		result := handler(args)
		writer.Write(result)
	}
}
