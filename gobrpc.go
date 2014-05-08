package birpc

import (
	"bufio"
	"encoding/gob"
	"io"
)

func NewGobCodec(conn io.ReadWriteCloser) Codec {
	wBuf := bufio.NewWriter(conn)
	return &codec{
		conn: conn,
		dec:  gob.NewDecoder(conn),
		enc:  gob.NewEncoder(wBuf),
		wBuf: wBuf,
	}
}
