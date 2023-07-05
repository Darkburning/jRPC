package codec

import (
	"bufio"
	"io"
	"jRPC/logger"
	"jRPC/protocol"
	"jRPC/serializer"
	"sync"
)

type ServerCodec struct {
	conn       io.ReadWriteCloser
	serializer serializer.JsonSerializer
	w          *bufio.Writer
	r          *bufio.Reader
}

func NewServerCodec(conn io.ReadWriteCloser) *ServerCodec {
	return &ServerCodec{
		conn: conn,
		w:    bufio.NewWriter(conn),
		r:    bufio.NewReader(conn),
	}
}

func (s *ServerCodec) ReadRequest() (*protocol.Request, error) {
	req := new(protocol.Request)
	reqBytes, err := recvFrame(s.r)
	if err != nil {
		logger.Warnln("rpc server: serverCodec ReadRequest: " + err.Error())
		return nil, err
	}
	logger.Debugln("rpc server: ReadRequest JSON:" + string(reqBytes))

	err = s.serializer.Unmarshal(reqBytes, req)
	if err != nil {
		logger.Warnln("rpc server: serverCodec ReadRequest: " + err.Error())
	}
	return req, nil
}

// WriteRes 直接把结果写回，不使用json
func (s *ServerCodec) WriteRes(replies string, mutex *sync.Mutex) {
	mutex.Lock()
	defer func() {
		err := s.w.Flush() // 将所有的缓存数据写入底层的IO接口
		if err != nil {
			_ = s.Close() // 发生错误则关闭
		}
		mutex.Unlock()
	}()

	respBytes := []byte(replies)

	err := sendFrame(s.w, respBytes)
	if err != nil {
		logger.Warnln("rpc server: serverCodec WriteResponse: " + err.Error())
		return
	}
}

func (s *ServerCodec) Close() error {
	return s.conn.Close()
}
