package birpc

import (
	"bufio"
	"io"
)

type codec struct {
	conn io.ReadWriteCloser
	dec  Decoder
	enc  Encoder
	wBuf *bufio.Writer
}

func (c *codec) Write(rs *RepReq, v interface{}) (err error) {
	if err = c.enc.Encode(rs); err != nil {
		return
	}
	if err = c.enc.Encode(v); err != nil {
		return
	}
	return c.wBuf.Flush()
}

func (c *codec) ReadHeader(res *RepReq) (err error) {
	return c.dec.Decode(res)
}

func (c *codec) ReadBody(v interface{}) (err error) {
	return c.dec.Decode(v)
}

func (c *codec) Close() (err error) {
	return c.conn.Close()
}

// A Codec implements reading of RPC requests and writing of RPC
// Responsess for the server side of an RPC session. The server calls
// ReadHeader and ReadBody in pairs to read requests from the
// connnection, an dit calls Write to write a response back. The
// server calls Close when finished with the connection. ReadBody
// may be called with a nil argument to force the body of the
// request to read and discarded
type Codec interface {
	ReadHeader(*RepReq) error
	ReadBody(interface{}) error
	Write(*RepReq, interface{}) error

	Close() error
}

type Decoder interface {
	Decode(interface{}) error
}

type Encoder interface {
	Encode(interface{}) error
}
