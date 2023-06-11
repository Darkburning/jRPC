package codec

import (
	"encoding/binary"
	"io"
	"net"
)

// 将信息封装成Frame再发送，避免粘包

func sendFrame(w io.Writer, data []byte) (err error) {
	var size [binary.MaxVarintLen64]byte

	if len(data) == 0 {
		n := binary.PutUvarint(size[:], uint64(0))
		if err = write(w, size[:n]); err != nil {
			return
		}
		return
	}

	// 对消息体长度进行编码得到消息头，变长头部
	n := binary.PutUvarint(size[:], uint64(len(data)))
	// 先往io写入消息头
	if err = write(w, size[:n]); err != nil {
		return
	}
	// 再写入消息体
	if err = write(w, data); err != nil {
		return
	}
	return
}

func recvFrame(r io.Reader) (data []byte, err error) {
	// 逐字节解析，直至得到变长参数
	size, err := binary.ReadUvarint(r.(io.ByteReader))
	if err != nil {
		return nil, err
	}
	if size != 0 {
		data = make([]byte, size)
		if err = read(r, data); err != nil {
			return nil, err
		}
	}
	return data, nil
}

// 处理粘包问题
func write(w io.Writer, data []byte) error {
	for index := 0; index < len(data); {
		n, err := w.Write(data[index:])
		if _, ok := err.(net.Error); !ok {
			return err
		}
		index += n
	}
	return nil
}

func read(r io.Reader, data []byte) error {
	for index := 0; index < len(data); {
		n, err := r.Read(data[index:])
		if err != nil {
			if _, ok := err.(net.Error); !ok {
				return err
			}
		}
		index += n
	}
	return nil
}
