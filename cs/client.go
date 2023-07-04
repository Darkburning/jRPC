package cs

import (
	"fmt"
	"jRPC/codec"
	"jRPC/logger"
	"jRPC/protocol"
	"net"
	"sync"
	"time"
)

const clientConnectTimeOut = time.Second * 3
const clientCallTimeOut = time.Second * 4

type Client struct {
	clientCodec *codec.ClientCodec
	sending     *sync.Mutex // 保证客户端并发安全
}

func NewClient(conn net.Conn) *Client {
	return &Client{
		clientCodec: codec.NewClientCodec(conn),
		sending:     new(sync.Mutex),
	}
}

func (c *Client) Close() error {
	return c.clientCodec.Close()
}

// Discover 发送并等待回复
func (c *Client) Discover(method string) bool {
	c.sending.Lock()
	// 处理发送请求到服务端,写数据导致的异常/超时
	sent := make(chan struct{})
	go func() {
		req := &protocol.Request{
			Method: "Discover:" + method,
			Args:   nil,
		}
		//time.Sleep(clientTimeOut + time.Second) // 测试写数据导致的异常/超时
		c.clientCodec.WriteRequest(req)
		sent <- struct{}{}
	}()

	select {
	//  超时直接返回
	case <-time.After(clientCallTimeOut):
		logger.Warnln(fmt.Sprintf("rpc client: WriteRequest timeout: expect within %v", clientCallTimeOut))
		c.sending.Unlock()
		return false
	case <-sent:
		c.sending.Unlock()
	}

	// 处理等待服务器处理导致的异常/超时和从服务端接收响应时，读数据导致的异常/超时
	read := make(chan struct{}, 1)
	var resp *string
	go func() {
		//time.Sleep(clientTimeOut + time.Second) // 测试从服务端接收响应时，读数据导致的异常/超时
		resp, _ = c.clientCodec.ReadRes()
		read <- struct{}{}
	}()

	select {
	//  超时直接返回
	case <-time.After(clientCallTimeOut):
		logger.Warnln(fmt.Sprintf("rpc client: ReadResponse timeout: expect within %v", clientCallTimeOut))
		return false
	case <-read:
		// 继续往后执行
	}

	logger.Debugln("rpc client: client discover success\n")
	if *resp == "The function has been registered!" {
		logger.Infoln(fmt.Sprintf("The function %s has been registered!", method))
		return true
	} else {
		logger.Infoln(fmt.Sprintf("The function %s has not been registered!", method))
		return false
	}

}

func (c *Client) Call(method string, args ...interface{}) string {

	c.sending.Lock()
	// 处理发送请求到服务端,写数据导致的异常/超时
	sent := make(chan struct{})
	go func() {
		req := &protocol.Request{
			Method: method,
			Args:   args,
		}
		//time.Sleep(clientTimeOut + time.Second) // 测试写数据导致的异常/超时
		c.clientCodec.WriteRequest(req)
		sent <- struct{}{}
	}()

	select {
	//  超时直接返回
	case <-time.After(clientCallTimeOut):
		logger.Warnln(fmt.Sprintf("rpc client: WriteRequest timeout: expect within %v", clientCallTimeOut))
		c.sending.Unlock()
		return fmt.Sprintf("rpc client: WriteRequest timeout: expect within %v", clientCallTimeOut)
	case <-sent:
	}

	// 处理等待服务器处理导致的异常/超时和从服务端接收响应时，读数据导致的异常/超时
	read := make(chan struct{}, 1)
	var res *string
	go func() {
		//time.Sleep(clientTimeOut + time.Second) // 测试从服务端接收响应时，读数据导致的异常/超时
		res, _ = c.clientCodec.ReadRes()
		read <- struct{}{}
		c.sending.Unlock()
	}()

	select {
	//  超时直接返回
	case <-time.After(clientCallTimeOut):
		logger.Warnln(fmt.Sprintf("rpc client: ReadResponse timeout: expect within %v", clientCallTimeOut))
		return fmt.Sprintf("rpc client: ReadResponse timeout: expect within %v", clientCallTimeOut)
	case <-read:
		logger.Debugln("rpc client: client call success\n")
		return *res
	}

}

// Dial 采用TCP连接，兼容Ipv6和Ipv4
// 可处理：与服务端建立连接,导致的异常/超时
func Dial(ip string, port string) (net.Conn, error) {
	var conn net.Conn
	var err error
	connected := make(chan struct{}, 1)

	go func() {
		conn, err = net.DialTimeout("tcp6", fmt.Sprintf("[%s]", ip)+":"+port, clientConnectTimeOut)
		if err != nil {
			logger.Debugln("rpc client: new client error: " + err.Error())
			conn, err = net.DialTimeout("tcp4", fmt.Sprintf("%s:%s", ip, port), clientConnectTimeOut)
			if err != nil {
				logger.Warnln("rpc client: new client error: " + err.Error())
			}
		}
		//time.Sleep(clientTimeOut + time.Second)  // 测试与服务端建立连接超时
		connected <- struct{}{}
	}()

	select {
	case <-time.After(clientConnectTimeOut):
		return nil, fmt.Errorf("rpc client: new client timeout: expect within %v", clientConnectTimeOut)
	case <-connected:
		return conn, nil
	}
}
