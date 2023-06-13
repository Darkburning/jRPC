package cs

import (
	"fmt"
	"io"
	"jRPC/codec"
	"jRPC/logger"
	"jRPC/protocol"
	"net"
	"reflect"
	"sync"
	"time"
)

const timeOutLimit = time.Second * 3

type Server struct {
	pending map[string]reflect.Value // 维护的service列表
}

func NewServer() *Server {
	return &Server{
		pending: make(map[string]reflect.Value),
	}
}

// serveCodec 流程：读取请求，处理请求，发送响应
// 每个serveCodec一把锁
func (s *Server) serveCodec(sc *codec.ServerCodec) {
	mutex := new(sync.Mutex)  // 保证写回的有序
	wg := new(sync.WaitGroup) // 等待直到所有的请求被处理完

	for {
		// 处理读取客户端请求数据时，读数据导致的异常/超时
		var req *protocol.Request
		var err error
		read := make(chan struct{})
		go func() {
			req, err = sc.ReadRequest()
			read <- struct{}{}
		}()

		select {
		case <-time.After(timeOutLimit):
			logger.Warnln(fmt.Sprintf("rpc server: ReadRequest timeout: expect within %v", timeOutLimit))
		case <-read:
			// 继续往后执行
		}
		if err != nil {
			if req == nil {
				logger.Warnln("rpc server: serveCodec readReq failed: Request empty")
				break
			}
			// 发送读取错误的报文
			logger.Warnln("rpc server: serveCodec readReq failed")
			sc.WriteResponse(err, nil, mutex)
			continue
		}
		wg.Add(1)                              // 需等待的协程+1
		go s.handleRequest(sc, req, mutex, wg) // 利用协程并发处理请求
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
func (s *Server) handleRequest(sc *codec.ServerCodec, req *protocol.Request, mutex *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	called := make(chan struct{}) // 处理调用映射服务的方法时，处理数据导致的异常/超时
	sent := make(chan struct{})   // 处理发送响应数据时，写数据导致的异常/超时

	go func() {
		{
			outArgs, err := s.call(req)
			called <- struct{}{}
			if err != nil {
				sc.WriteResponse(err, nil, mutex)
				sent <- struct{}{}
				return
			}
			// 发送响应报文
			sc.WriteResponse(nil, outArgs, mutex)
			sent <- struct{}{}
		}
	}()

	select {
	case <-time.After(timeOutLimit):
		errMsg := fmt.Errorf("rpc server: handleRequest timeout: expect within %v", timeOutLimit)
		sc.WriteResponse(errMsg, nil, mutex)
	case <-called:
		<-sent
	}
}

// Register 将方法注册到pending[string]reflect.Value
func (s *Server) Register(serviceName string, f interface{}) {
	if _, ok := s.pending[serviceName]; ok {
		logger.Warnln("rpc server: service already registered")
		return
	}

	fVal := reflect.ValueOf(f)
	if !fVal.IsValid() {
		logger.Warnln("rpc server: service registered failed - invalid service")
		return
	}

	s.pending[serviceName] = fVal
	logger.Infoln(fmt.Sprintf("rpc server: service registered: %s", serviceName))
}

func (s *Server) isMethodExists(method string) bool {
	if _, ok := s.pending[method]; ok {
		logger.Debugln(fmt.Sprintf("rpc server: method %s found", method))
		return true
	} else {
		logger.Warnln(fmt.Sprintf("rpc server: method %s not found ", method))
		return false
	}
}

func (s *Server) call(req *protocol.Request) ([]interface{}, error) {
	if !s.isMethodExists(req.Method) {
		return nil, fmt.Errorf("rpc server: method %s not found", req.Method)
	} else {
		// 根据函数原型构造入参切片
		inArgs := make([]reflect.Value, 0, len(req.Args))
		funcType := s.pending[req.Method].Type()
		for i := 0; i < funcType.NumIn(); i++ {
			argType := funcType.In(i)
			argValue := req.Args[i]
			// 将输入参数转换为反射值类型并添加到列表
			inArg := reflect.ValueOf(argValue).Convert(argType)
			inArgs = append(inArgs, inArg)
		}

		// 调用函数
		logger.Debugln(fmt.Sprintf("%v\n", inArgs))
		returnValues := s.pending[req.Method].Call(inArgs)
		logger.Infoln(fmt.Sprintf("rpc server: call %s success", req.Method))

		// 构造出参切片
		outArgs := make([]interface{}, 0, len(returnValues))
		for _, ret := range returnValues {
			outArgs = append(outArgs, ret.Interface())
		}
		logger.Debugln("rpc server: make outArgs success")
		return outArgs, nil
	}
}

// Accept 方法实现接收监听者的连接 开启协程处理每个到来的连接
func (s *Server) Accept(lis net.Listener) {
	logger.Infoln("rpc server: listen and serve......")
	for {
		conn, err := lis.Accept()
		if err != nil {
			logger.Warnln("rpc server: accept error:" + err.Error())
		}
		go s.ServeConn(conn)
	}
}
