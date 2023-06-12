package main

import (
	"flag"
	"fmt"
	. "jRPC/cs"
	"jRPC/logger"
	"net"
)

var (
	clientIp   string
	clientPort string
)

func main() {
	flag.StringVar(&clientIp, "i", "", logger.WarnMsg("client connect ip   usage:-i <ipv4/ipv6>"))
	flag.StringVar(&clientPort, "p", "", logger.WarnMsg("client connect port   usage:-p <port>"))
	flag.Parse()

	if clientIp == "" {
		logger.Fatalln("client's connect ip is required")
	}
	if clientPort == "" {
		logger.Fatalln("client's connect port is required")
	}

	addrStr := fmt.Sprintf("%s:%s", clientIp, clientPort)
	var conn net.Conn
	var err error

	conn, err = net.Dial("tcp6", fmt.Sprintf("[%s]", clientIp)+":"+clientPort)
	if err != nil {
		conn, err = net.Dial("tcp4", addrStr)
	}
	if err != nil {
		logger.Fatalln(err.Error())
	}

	client := NewClient(conn)
	res1 := client.Call("Sum", 2, 2)
	res2 := client.Call("Product", 3, 3)
	res3 := client.Call("Revert", "HELLO")
	fmt.Printf("2 + 2  = %v\n", res1[0])
	fmt.Printf("3 * 3  = %v\n", res2[0])
	fmt.Printf("Revert %s to %s", "HELLO", res3[0])
	client.Call("NoFunc", "HELLO")
}
