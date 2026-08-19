package main

import (
	"io"
	"log"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", ":443")
	if err != nil {
		log.Fatal(err)
	}
	log.Println("gateway listening on :443")

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handle(conn)
	}
}

func handle(client net.Conn) {
	defer client.Close()

	unix, err := net.Dial("unix", "/tmp/reality.sock")
	if err != nil {
		log.Println("unix dial error:", err)
		return
	}
	defer unix.Close()

	go io.Copy(unix, client)
	io.Copy(client, unix)
}
