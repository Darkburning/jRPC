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

	req := &protocol.Request{
		Method: method,
		Args:   args,
	}
	c.clientCodec.WriteRequest(req)
	c.sending.Unlock()

	var err error
	for err == nil {
		resp, err := c.clientCodec.ReadResponse()
		if err != nil {
			logger.Warnln("rpc client: client receive: " + err.Error())
		}
		if resp.Err != "" {
			logger.Warnln("rpc client: client receive: " + resp.Err)
			return nil
		} else {
			logger.Infoln("rpc client: client call success!\n")
			for idx, reply := range resp.Replies {
				logger.Infoln(fmt.Sprintf("Value %d is : %v", idx, reply))
			}
			return resp.Replies
		}
	}
	return nil
}

// Dial 处理建立连接超时
func Dial(network string, addr string) (*Client, error) {
	conn, err := net.DialTimeout(network, addr, timeOutLimit)
	if err != nil {
		return nil, err
	} else {
		defer func() {
			if err != nil {
				_ = conn.Close()
			}
		}()
		// 创建子协程，创建一个客户端
		ch := make(chan *Client)
		go func() {
			//time.Sleep(time.Second * 4) // for test
			ch <- NewClient(conn)
		}()

		// select多路复用处理阻塞IO，在两个信道上监听
		select {
		case <-time.After(timeOutLimit):
			return nil, fmt.Errorf("rpc client: new client timeout: expect within %v", timeOutLimit)
		case result := <-ch:
			return result, nil
		}

	}
}
