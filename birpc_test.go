package birpc

import (
	"bufio"
	"encoding/json"
	"log"
	"net"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	client     *Protocol
	clientConn net.Conn
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (ts *TestSuite) SetupSuite() {
	var err error
	var listener net.Listener

	// Acquire an available ipport
	listener, serverAddr := listenTCP()

	// Client connecting to server
	ts.clientConn, err = net.Dial("tcp", serverAddr)
	assert.NoError(ts.T(), err, "Connecting a TCP server should not return an error : %v", err)

	// Start server
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Fatal("Connection accept error : %v", err)
			}
			c := NewJsonCodec(conn)
			s := NewProtocol(c)
			// Prepare service and methods
			s.Register(new(Arith))
			// Parse request/response
			go s.Serve()
		}
	}()
}

type Arith struct{}
type Args struct {
	A, B int
}

type Reply struct {
	C int
}

func (a *Arith) Add(args Args, reply *Reply) error {
	reply.C = args.A + args.B
	return nil
}

func listenTCP() (net.Listener, string) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal("net.Listen tcp :0: %v", err)
	}
	return l, l.Addr().String()
}

func (ts *TestSuite) TestIsExported() {
	result := isExported("A")
	assert.Equal(ts.T(), result, true, "Unexpected error for isExported(\"A\"), [expect : %v, got : %v]", result, true)
	result = isExported("a")
	assert.Equal(ts.T(), result, false, "Unexpected error for isExported(\"a\"), [expect : %v, got : %v]", result, true)
}

func (ts *TestSuite) TestIsExportedOrBuiltin() {
	arithType := reflect.TypeOf(&Arith{})
	result := isExportedOrBuiltinType(arithType)
	assert.Equal(ts.T(), result, true, "Unexpected error for isExported(\"A\"), [expect : %v, got : %v]", result, true)
}

func (ts *TestSuite) TestSendRequest() {

	call := &Call{
		ServiceMethod: "Arith.Add",
		Args:          &Args{7, 8},
		Reply:         new(Reply),
		Done:          make(chan *Call, 10),
	}

	// Prepare a client codec
	bWriter := bufio.NewWriter(ts.clientConn)
	cc := &codec{
		conn: ts.clientConn,
		dec:  json.NewDecoder(ts.clientConn),
		enc:  json.NewEncoder(bWriter),
		wBuf: bWriter,
	}

	// Create a protocol and prepare
	ts.client = NewProtocol(cc)
	go ts.client.Serve()

	go ts.client.sendRequest(call)
	<-call.Done
	ts.client.Close()
}
