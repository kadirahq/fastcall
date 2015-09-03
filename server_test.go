package fastcall

import (
	"bytes"
	"net"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	s := NewServer("localhost:1337", func(req []byte) (res []byte) {
		exp := []byte{5, 5, 5, 5}
		if !bytes.Equal(req, exp) {
			t.Fatal("wrong value")
		}

		return []byte{6, 6, 6, 6}
	})

	go s.Listen()
	time.Sleep(time.Millisecond)

	conn, err := net.Dial("tcp", "localhost:1337")
	if err != nil {
		t.Fatal(err)
	} else {
		defer conn.Close()
	}

	c := NewConn(conn)

	// first 4 bytes make the call id
	src := []byte{1, 1, 1, 1, 5, 5, 5, 5}
	if err := c.Write(src); err != nil {
		t.Fatal(err)
	}

	dst, err := c.Read()
	if err != nil {
		t.Fatal(err)
	}

	exp := []byte{1, 1, 1, 1, 6, 6, 6, 6}
	if !bytes.Equal(exp, dst) {
		t.Fatal("wrong value")
	}
}

func ExampleServer() {
	s := NewServer("localhost:1337", func(req []byte) (res []byte) {
		return req
	})

	if err := s.Listen(); err != nil {
		panic(err)
	}
}
