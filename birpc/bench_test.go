// Copyright 2013 René Kistl. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bidirpc

import (
	"log"
	"net"
	//"runtime"
	"sync"
	//"sync/atomic"
	"testing"
)

var (
	serverBSON, serverGob, serverMsgPack string
	onceBSON, onceGob, onceMsgPack       sync.Once
)

type Args struct {
	A, B int
}

type Reply struct {
	C int
}

type Arith struct {
}

func (t *Arith) Add(args Args, reply *Reply) error {
	reply.C = args.A + args.B
	return nil
}

func listenTCP() (net.Listener, string) {
	l, e := net.Listen("tcp", "127.0.0.1:0") // any available address
	if e != nil {
		log.Fatalf("net.Listen tcp :0: %v", e)
	}
	return l, l.Addr().String()
}

func startGobServer() {
	var l net.Listener
	l, serverGob = listenTCP()

	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				log.Fatal("accept error:", err)
			}

			c := NewCodec(conn)
			s := NewProtocol(c)
			s.Register(new(Arith))
			go s.Serve()
		}
	}()
}

func TestGobConnection(t *testing.T) {
	onceGob.Do(startGobServer)

	conn, err := net.Dial("tcp", serverGob)
	if err != nil {
		log.Fatal("error dialing:", err)
	}
	args := &Args{7, 8}
	reply := new(Reply)
	client := NewClient(conn)
	err = client.Call("Arith.Add", args, reply)
	if err != nil {
		t.Fatalf("rpc error: Add: expected no error but got string %q", err.Error())
	}
	if reply.C != args.A+args.B {
		t.Fatalf("rpc error: Add: expected %d got %d", reply.C, args.A+args.B)
	}
	client.Close()
}
