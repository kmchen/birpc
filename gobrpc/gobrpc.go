package birpc

import (
	"bufio"
	"encoding/gob"
	"io"
)

func NewCodec(conn io.ReadWriteCloser) Codec {
	wBuf := bufio.NewWriter(conn)
	return &codec{
		conn: conn,
		dec:  gob.NewDecoder(conn),
		enc:  gob.NewEncoder(wBuf),
		wBuf: wBuf,
	}
}
