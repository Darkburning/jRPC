package codec

import (
	"encoding/binary"
	"io"
	"net"
)

// 将信息封装成Frame 头部编码为4字节再发送，避免粘包

func sendFrame(w io.Writer, data []byte) (err error) {
	size := make([]byte, 4)
	if len(data) == 0 {
		binary.LittleEndian.PutUint32(size, uint32(0))
		if err = write(w, size); err != nil {
			return
		}
		return
	}

	binary.LittleEndian.PutUint32(size, uint32(len(data)))
	if err = write(w, size); err != nil {
		return err
	}
	if err = write(w, data); err != nil {
		return err
	}
	return nil
}

func recvFrame(r io.Reader) (data []byte, err error) {
	size := make([]byte, 4)
	if _, err := io.ReadFull(r, size); err != nil {
		return nil, err
	}

	data = make([]byte, binary.LittleEndian.Uint32(size))
	if err := read(r, data); err != nil {
		return nil, err
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
