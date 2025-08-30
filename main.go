package main

import (
	"fmt"
	"io"
	"log"
	"net"
)


func main() {
	fmt.Println("Listening on port :6379")
	// create a server 
	l,err := net.Listen("tcp",":6379")
	if err !=  nil {
		log.Fatal("port binding error",err)
		return
	}
	// listen for a connection
	conn,err := l.Accept()	
	if err != nil {
		log.Fatal("connot accept the connection",err)
		return 
	}
	defer conn.Close()	

	for {
		resp := NewResp(conn)
		_,err := resp.Read()	
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal("error in reading",err)
			return
		}
		writer := newWritter(conn)		
		err = writer.Write(Value{typ:"string",str:"hello"})
		if err != nil {
			log.Fatal("error in writing",err)
			return
		}
	}
}
