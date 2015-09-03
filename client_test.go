package fastcall

import (
	"bytes"
	"net"
	"sync"
	"testing"
)

func TestClient(t *testing.T) {
	l, err := net.Listen("tcp", "localhost:1337")
	if err != nil {
		t.Fatal(err)
	} else {
		defer l.Close()
	}

	c := NewClient("localhost:1337")
	go c.Connect()

	var w sync.WaitGroup
	w.Add(1)

	go func() {
		conn, err := l.Accept()
		if err != nil {
			t.Fatal(err)
		}

		d := NewConn(conn)
		dst, err := d.Read()
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(dst, []byte{1, 0, 0, 0, 1, 2, 3, 4}) {
			t.Fatal("wrong value")
		}

		if err := d.Write([]byte{1, 0, 0, 0, 4, 5, 6, 7}); err != nil {
			t.Fatal(err)
		}
	}()

	go func() {
		res, err := c.Call([]byte{1, 2, 3, 4})
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(res, []byte{4, 5, 6, 7}) {
			t.Fatal("wrong value")
		}

		w.Done()
	}()

	w.Wait()
}

func ExampleClient() {
	p := []byte{1, 2, 3, 4}
	c := NewClient("localhost:1337")
	go c.Connect()

	res, err := c.Call(p)
	if err != nil {
		panic(err)
	}

	if !bytes.Equal(res, p) {
		panic("wrong echo")
	}
}
