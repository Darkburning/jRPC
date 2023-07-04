package main

import (
	"fmt"
	. "jRPC/cs"
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
	t.Parallel()
	log.SetFlags(0)
	ch := make(chan string)
	go startServer(ch)
	addr := <-ch

	ip, port, err := net.SplitHostPort(addr)
	if err != nil {
		log.Fatal("SplitHostPort error")
	}

	conn, err := Dial(ip, port)
	if err != nil {
		log.Fatal("Client Dial Failed")
	}
	client := NewClient(conn)
	res1 := client.Call("Add", 2, 2)
	res2 := client.Call("Multi", 3, 3)
	fmt.Printf("2 + 2 = %v\n", res1[0])
	fmt.Printf("3 * 3  = %v\n", res2[0])
}
