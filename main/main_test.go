package main

import (
	"log"
	"testing"
)

func TestServerHandleTimeOut(t *testing.T) {
	t.Parallel()
	ch := make(chan string)
	go startServer(ch)

	addr := <-ch
	client, err := Dial(network, addr)
	if err != nil {
		log.Fatal(err)
	}
	client.Call("Sleep", 4)
}

func TestNewClientTimeOut(t *testing.T) {
	t.Parallel()
	ch := make(chan string)
	go startServer(ch)

	addr := <-ch
	client, err := Dial(network, addr)
	if err != nil {
		log.Fatal(err)
	}
	client.Call("Sleep", 4)
}
