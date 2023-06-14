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

const clientTimeOut = time.Second * 3

type Client struct {
	clientCodec *codec.ClientCodec
	sending     *sync.Mutex
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

func (c *Client) Call(method string, args ...interface{}) []interface{} {

	c.sending.Lock()
	// 处理发送请求到服务端,写数据导致的异常/超时
	sent := make(chan struct{})
	go func() {
		req := &protocol.Request{
			Method: method,
			Args:   args,
		}
		c.clientCodec.WriteRequest(req)
		sent <- struct{}{}
	}()

	select {
	//  超时直接返回
	case <-time.After(clientTimeOut):
		logger.Warnln(fmt.Sprintf("rpc client: WriteRequest timeout: expect within %v", clientTimeOut))
		c.sending.Unlock()
		return nil
	case <-sent:
		c.sending.Unlock()
	}

	// 处理等待服务器处理导致的异常/超时和从服务端接收响应时，读数据导致的异常/超时
	read := make(chan struct{})
	var err error
	var resp *protocol.Response
	for err == nil {
		go func(response *protocol.Response, err error) {
			resp, err = c.clientCodec.ReadResponse()
			read <- struct{}{}
		}(resp, err)

		select {
		//  超时直接返回
		case <-time.After(clientTimeOut):
			logger.Warnln(fmt.Sprintf("rpc client: ReadResponse timeout: expect within %v", clientTimeOut))
			return nil
		case <-read:
			// 继续往后执行
		}
		// call不存在，服务端照样处理了
		if err != nil {
			logger.Warnln("rpc client: client receive: " + err.Error())
		}
		// 服务端处理出错
		if resp.Err != "" {
			logger.Warnln("rpc client: client receive: " + resp.Err)
			return nil
		} else {
			// call存在，服务端处理正常，读取replies的值
			logger.Debugln("rpc client: client call success\n")
			for idx, reply := range resp.Replies {
				logger.Debugln(fmt.Sprintf("Value %d is : %v", idx, reply))
			}
			return resp.Replies
		}
	}
	return nil
}

// Dial 采用TCP连接，兼容Ipv6和Ipv4
// 可处理：与服务端建立连接,导致的异常/超时
func Dial(ip string, port string) (net.Conn, error) {
	var conn net.Conn
	var err error
	connected := make(chan struct{})

	go func() {
		conn, err = net.DialTimeout("tcp6", fmt.Sprintf("[%s]", ip)+":"+port, clientTimeOut)
		if err != nil {
			logger.Debugln("rpc client: new client error: " + err.Error())
			conn, err = net.DialTimeout("tcp4", fmt.Sprintf("%s:%s", ip, port), clientTimeOut)
			if err != nil {
				logger.Warnln("rpc client: new client error: " + err.Error())
			}
		}
		connected <- struct{}{}
	}()

	select {
	case <-time.After(clientTimeOut):
		return nil, fmt.Errorf("rpc client: new client timeout: expect within %v", clientTimeOut)
	case <-connected:
		return conn, nil
	}
}
