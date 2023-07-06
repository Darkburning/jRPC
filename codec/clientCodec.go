package codec

import (
	"bufio"
	"io"
	"jRPC/logger"
	"jRPC/protocol"
	"jRPC/serializer"
	"net"
)

// ClientCodec 持有连接
type ClientCodec struct {
	conn       io.ReadWriteCloser
	serializer serializer.JsonSerializer
	w          *bufio.Writer
	r          *bufio.Reader
}

func NewClientCodec(conn net.Conn) *ClientCodec {
	return &ClientCodec{
		conn: conn,
		w:    bufio.NewWriter(conn),
		r:    bufio.NewReader(conn),
	}
}

func (c *ClientCodec) ReadResponse() (*protocol.Response, error) {
	resp := new(protocol.Response)
	byteResp, err := recvFrame(c.r)
	if err != nil {
		logger.Warnln("rpc client: clientCodec ReadResponse: " + err.Error())
		return nil, err
	}
	err = c.serializer.Unmarshal(byteResp, resp)
	if err != nil {
		logger.Warnln("rpc client: clientCodec ReadResponse: " + err.Error())
		return nil, err
	}
	return resp, nil

}
func (c *ClientCodec) WriteRequest(req *protocol.Request) {
	defer func(w *bufio.Writer) {
		err := w.Flush()
		if err != nil {
			logger.Warnln("rpc client: clientCodec WriteRequest: " + err.Error())
		}
	}(c.w)

	reqBytes, err := c.serializer.Marshal(req)
	if err != nil {
		logger.Warnln("rpc client: clientCodec WriteRequest: " + err.Error())
		return
	}
	logger.Debugln("rpc client: Request serialized:" + string(reqBytes))

	err = sendFrame(c.w, reqBytes)
	if err != nil {
		logger.Warnln("rpc client: clientCodec WriteRequest: " + err.Error())
		return
	}

}
func (c *ClientCodec) Close() error {
	return c.conn.Close()
}
