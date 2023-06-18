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
	server.Register("Add", Add)
	server.Register("Substract", Substract)
	server.Register("Consub", Consub)
	server.Register("Multi", Multi)
	server.Register("Divide", Divide)
	server.Register("Condiv", Condiv)
	server.Register("Power", Power)
	server.Register("Mod", Mod)
	server.Register("Sqrtmul", Sqrtmul)
	server.Register("Triangle", Triangle)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Fatalln(err.Error())
	}

	server.Accept(lis)
}
