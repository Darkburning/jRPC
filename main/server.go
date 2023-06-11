package main

import (
	"errors"
	"fmt"
	"io"
	"jRPC/codec"
	"jRPC/protocol"
	"log"
	"net"
	"reflect"
	"sync"
	"time"
)

const timeOutLimit = time.Second * 3

type Server struct {
	sending *sync.Mutex              // 保证线程安全
	pending map[string]reflect.Value // 维护的service列表
}

func NewServer() *Server {
	return &Server{
		sending: new(sync.Mutex),
		pending: make(map[string]reflect.Value),
	}
}

// serveCodec 流程：读取请求，处理请求，发送响应
func (s *Server) serveCodec(sc *codec.ServerCodec) {
	s.sending.Lock()
	defer s.sending.Unlock()
	wg := new(sync.WaitGroup) // 等待直到所有的请求被处理完
	for {
		req, err := sc.ReadRequest()
		if err != nil {
			if req == nil {
				break
			}
			// 发送读取错误的报文
			sc.WriteResponse(err, nil)
			continue
		}
		wg.Add(1)                       // 需等待的协程+1
		go s.handleRequest(sc, req, wg) // 利用协程并发处理请求
	}
	wg.Wait() // 等待所有请求的处理结束
	_ = sc.Close()
}

func (s *Server) ServeConn(conn io.ReadWriteCloser) {
	defer func() {
		err := conn.Close()
		if err != nil {
			return
		}
	}()
	// 每个连接一个，负责编解码并读取数据
	serverCodec := codec.NewServerCodec(conn)

	s.serveCodec(serverCodec)
}

// handleRequest 根据请求读取入参并调用方法得到出参然后返回数据
func (s *Server) handleRequest(sc *codec.ServerCodec, req *protocol.Request, wg *sync.WaitGroup) {
	defer wg.Done()
	called := make(chan struct{}) // 调用超时
	sent := make(chan struct{})   // 回复超时

	go func() {
		inArgs := make([]reflect.Value, 0, len(req.Args))
		for _, arg := range req.Args {
			argVal := reflect.ValueOf(arg)
			inArgs = append(inArgs, argVal)
			fmt.Printf("%v\t", argVal.Interface())
		}

		outArgs, err := s.call(req.Method, inArgs)
		called <- struct{}{}
		if err != nil {
			sc.WriteResponse(err, nil)
			sent <- struct{}{}
			return
		}
		// 发送响应报文
		sc.WriteResponse(nil, outArgs)
		sent <- struct{}{}
	}()

	select {
	case <-time.After(timeOutLimit):
		errMsg := fmt.Errorf("rpc server: handleRequest timeout: expect within %v", timeOutLimit)
		sc.WriteResponse(errMsg, nil)
	case <-called:
		<-sent
	}
}

// Register 将方法注册到pending[string]reflect.Value
func (s *Server) Register(serviceName string, f interface{}) {
	if _, ok := s.pending[serviceName]; ok {
		log.Println("rpc server: Already Registered!")
		return
	}

	fVal := reflect.ValueOf(f)
	if !fVal.IsValid() {
		log.Println("rpc server: Unable to add function to pending - invalid value")
		return
	}

	s.pending[serviceName] = fVal
	log.Printf("rpc server: Function added to pending: %s", serviceName)
}

func (s *Server) isMethodExists(method string) bool {
	if _, ok := s.pending[method]; ok {
		log.Printf("rpc server: Method Exists\n")
		return true
	} else {
		log.Printf("rpc server: Method Not Exists\n")
		return false
	}
}

func (s *Server) call(serviceName string, inArgs []reflect.Value) ([]interface{}, error) {
	if !s.isMethodExists(serviceName) {
		return nil, errors.New("rpc server: no Func Error")
	} else {
		log.Printf("%v\n", inArgs)
		returnValues := s.pending[serviceName].Call(inArgs)
		log.Printf("rpc server: Called Success!\n")

		outArgs := make([]interface{}, 0, len(returnValues))
		for _, ret := range returnValues {
			outArgs = append(outArgs, ret.Interface())
		}
		log.Printf("rpc server: Make outArgs Success!\n")
		return outArgs, nil
	}
}

// Accept 方法实现接收监听者的连接 开启协程处理每个到来的连接
// 若想启动服务只需传入listener，TCP或UNIX协议均可
func (s *Server) Accept(lis net.Listener) {
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Println("rpc server: accept error:", err)
		}
		go s.ServeConn(conn)
	}
}
