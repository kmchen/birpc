package birpc

import (
	"bufio"
	"encoding/json"
	"io"
)

func NewJsonCodec(conn io.ReadWriteCloser) Codec {
	wBuf := bufio.NewWriter(conn)
	return &codec{
		conn: conn,
		dec:  json.NewDecoder(conn),
		enc:  json.NewEncoder(wBuf),
		wBuf: wBuf,
	}
}
