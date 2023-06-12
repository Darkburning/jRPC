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
	server.Register("Sum", Sum)
	server.Register("Product", Product)
	server.Register("Sleep", Sleep)

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

	client, err := Dial("tcp", addr)
	if err != nil {
		log.Fatal("Client Dial Failed")
	}
	res1 := client.Call("Sum", 2, 2)
	res2 := client.Call("Product", 3, 3)
	fmt.Printf("2 + 2 = %v\n", res1[0])
	fmt.Printf("3 * 3  = %v\n", res2[0])
}
