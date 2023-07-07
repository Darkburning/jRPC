package main

import (
	"fmt"
	. "jRPC/cs"
	"jRPC/logger"
	. "jRPC/service"
	"log"
	"net"
	"testing"
)

func startServer(ch chan string) {
	server := NewServer()
	server.Register("Add", Add)
	server.Register("Multi", Multi)

	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal("server listen failed")
	}

	ch <- lis.Addr().String()
	server.Accept(lis)
}

func TestCS(t *testing.T) {
	log.SetFlags(0)
	ch := make(chan string)
	go startServer(ch)
	addr := <-ch

	conn, err := Dial(addr)
	if err != nil {
		log.Fatal("Client Dial Failed")
	}
	client := NewClient(conn)

	if client.Discover("NotExists") {
		logger.Infoln("Discover Success!")
	} else {
		logger.Infoln("Discover failed!")
	}
	res1 := client.Call("Add", 2, 2)
	res2 := client.Call("Multi", 3, 3)
	fmt.Printf("2 + 2 = %v\n", res1)
	fmt.Printf("3 * 3  = %v\n", res2)
}
