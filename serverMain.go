package main

import (
	"flag"
	"fmt"
	. "jRPC/cs"
	"jRPC/logger"
	. "jRPC/service"
	"net"
)

var (
	serverIp   string
	serverPort string
)

func main() {
	flag.StringVar(&serverIp, "l", "0.0.0.0", logger.WarnMsg("server listen ip   usage:-l <ipv4/ipv6>"))
	flag.StringVar(&serverPort, "p", "", logger.WarnMsg("server listen port   usage:-p <port>"))
	flag.Parse()

	if serverPort == "" {
		logger.Fatalln("server port is required")
	}

	addr := fmt.Sprintf("[%s]:%s", serverIp, serverPort)

	server := NewServer()
	server.Register("Sum", Sum)
	server.Register("Subtract", Subtract)
	server.Register("Division", Division)
	server.Register("Product", Product)
	server.Register("Square", Square)
	server.Register("Cube", Cube)
	server.Register("Sleep", Sleep)
	server.Register("Upper", Upper)
	server.Register("Lower", Lower)
	server.Register("Revert", Revert)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Fatalln(err.Error())
	}

	server.Accept(lis)
}
