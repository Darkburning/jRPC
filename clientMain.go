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

	ok := client.Discover("NotExists")
	if ok {
		logger.Infoln("Discover Success!")
	} else {
		logger.Infoln("Discover failed!")
	}
	res1 := client.Call("Add", 2, 2)
	res2 := client.Call("Substract", 2, 2)
	fmt.Printf("远程调用的响应消息：%v\n", res1[0])
	fmt.Printf("远程调用的响应消息：%v\n", res2[0])
	client.Call("NoFunc", "HELLO")
}
