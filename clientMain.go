package main

import (
	"flag"
	"fmt"
	. "jRPC/cs"
	"jRPC/logger"
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

	conn, err := Dial(clientIp, clientPort)
	if err != nil {
		logger.Fatalln(err.Error())
	}
	client := NewClient(conn)

	res1 := client.Call("Sum", 2, 2)
	res2 := client.Call("Product", 3, 3)
	res3 := client.Call("Revert", "HELLO")
	fmt.Printf("2 + 2  = %v\n", res1[0])
	fmt.Printf("3 * 3  = %v\n", res2[0])
	fmt.Printf("Revert %s to %s\n", "HELLO", res3[0])
	client.Call("NoFunc", "HELLO")
	//client.Call("Sleep", 4)
}
