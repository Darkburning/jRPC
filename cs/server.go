package cs

import (
	"fmt"
	"io"
	"jRPC/codec"
	"jRPC/logger"
	"jRPC/protocol"
	"net"
	"reflect"
	"strings"
	"sync"
	"time"
)

const serverHandleReqTimeOut = time.Second * 3
const serverReadReqTimeOut = time.Second * 5

type Server struct {
	serviceList map[string]reflect.Value // 维护的service列表
}

func NewServer() *Server {
	return &Server{
		serviceList: make(map[string]reflect.Value),
	}
}

// serveCodec 流程：读取请求，处理请求，发送响应
// 每个serveCodec一把锁
func (s *Server) serveCodec(sc *codec.ServerCodec) {
	mutex := new(sync.Mutex)  // 保证写回的有序
	wg := new(sync.WaitGroup) // 等待直到所有的请求被处理完
	defer func() {
		wg.Wait()
		_ = sc.Close()
	}()

	for {
		// 处理读取客户端请求数据时，读数据导致的异常/超时
		var req *protocol.Request
		var err error
		read := make(chan struct{})
		go func() {
			//time.Sleep(serverTimeOut + time.Second) // 测试处理读取客户端请求数据时，读数据导致的异常/超时
			req, err = sc.ReadRequest()
			read <- struct{}{}
		}()

		select {
		case <-time.After(serverReadReqTimeOut):
			logger.Warnln(fmt.Sprintf("rpc server: ReadRequest timeout: expect within %v", serverReadReqTimeOut))
			return
		case <-read:
			// 继续往后执行
		}

		if err != nil {
			if req == nil {
				logger.Warnln("rpc server: serveCodec ReadRequest failed: Request empty")
				break
			}
			// 发送读取错误的报文
			sc.WriteResponse(err, nil, mutex)
			continue
		}
		wg.Add(1) // 需等待的协程+1
		// 利用协程并发处理请求
		if strings.HasPrefix(req.Method, "Discover:") {
			go s.handleDiscover(sc, req, mutex, wg)
		} else {
			go s.handleRequest(sc, req, mutex, wg)
		}
	}
	// 等待所有请求的处理结束
}

func (s *Server) serveConn(conn io.ReadWriteCloser) {
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

func (s *Server) handleDiscover(sc *codec.ServerCodec, req *protocol.Request, mutex *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	// 获得期待调用的方法
	method := strings.TrimLeft(req.Method, "Discover:")
	sent := make(chan struct{}) // 处理发送响应数据时，写数据导致的异常/超时
	var replies string

	if s.isMethodExists(method) {
		// 返回方法存在
		replies = "The function has been registered!"
	} else {
		// 返回方法不存在
		replies = "The function has not been registered!"
	}
	go func() {
		ret := make([]interface{}, 0, 1)
		ret = append(ret, replies)
		sc.WriteResponse(nil, ret, mutex)
		sent <- struct{}{}
	}()

	select {
	case <-time.After(serverHandleReqTimeOut):
		errMsg := fmt.Errorf("rpc server: handleDiscover timeout: expect within %v", serverHandleReqTimeOut)
		logger.Warnln(errMsg.Error())
		sc.WriteResponse(errMsg, nil, mutex)
	case <-sent:
		// end
	}

}

// handleRequest 根据请求读取入参并调用方法得到出参然后返回数据
func (s *Server) handleRequest(sc *codec.ServerCodec, req *protocol.Request, mutex *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	called := make(chan struct{}, 1) // 处理调用映射服务的方法时，处理数据导致的异常/超时
	sent := make(chan struct{})      // 处理发送响应数据时，写数据导致的异常/超时

	go func() {
		//time.Sleep(serverTimeOut + time.Second) //  测试处理调用映射服务的方法时，处理数据导致的异常/超时
		outArgs, err := s.call(req)
		called <- struct{}{}
		if err != nil {
			sc.WriteResponse(err, nil, mutex)
			sent <- struct{}{}
			return
		}
		//time.Sleep(serverTimeOut + time.Second) //   处理发送响应数据时，写数据导致的异常/超时
		sc.WriteResponse(nil, outArgs, mutex)
		sent <- struct{}{}
	}()

	select {
	case <-time.After(serverHandleReqTimeOut):
		errMsg := fmt.Errorf("rpc server: handleRequest timeout: expect within %v", serverHandleReqTimeOut)
		logger.Warnln(errMsg.Error())
		sc.WriteResponse(errMsg, nil, mutex)
	case <-called:
		{
			select {
			case <-time.After(serverHandleReqTimeOut):
				errMsg := fmt.Errorf("rpc server: handleRequest timeout: expect within %v", serverHandleReqTimeOut)
				logger.Warnln(errMsg.Error())
				sc.WriteResponse(errMsg, nil, mutex)
			case <-sent:
				// end
			}

		}
	}
}

// Register serviceList[string]reflect.Value
func (s *Server) Register(serviceName string, f interface{}) {
	if _, ok := s.serviceList[serviceName]; ok {
		logger.Warnln("rpc server: service already registered")
		return
	}

	fVal := reflect.ValueOf(f)
	if !fVal.IsValid() {
		logger.Warnln("rpc server: service registered failed - invalid service")
		return
	}

	s.serviceList[serviceName] = fVal
	logger.Infoln(fmt.Sprintf("rpc server: service registered: %s", serviceName))
}

func (s *Server) isMethodExists(method string) bool {
	if _, ok := s.serviceList[method]; ok {
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
		//return nil, errors.New("The function has not been registered!")
	} else {
		// 根据函数原型构造入参切片
		inArgs := make([]reflect.Value, 0, len(req.Args))
		funcType := s.serviceList[req.Method].Type()
		for i := 0; i < funcType.NumIn(); i++ {
			argType := funcType.In(i)
			argValue := req.Args[i]
			// 将输入参数转换为反射值类型并添加到列表
			inArg := reflect.ValueOf(argValue).Convert(argType)
			inArgs = append(inArgs, inArg)
		}

		// 调用函数
		logger.Debugln(fmt.Sprintf("%v\n", inArgs))
		returnValues := s.serviceList[req.Method].Call(inArgs)
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
		go s.serveConn(conn)
	}
}
