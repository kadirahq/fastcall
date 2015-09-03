package fastcall

import (
	"bytes"
	"net"
	"sync"
	"testing"
)

func TestConn(t *testing.T) {
	l, err := net.Listen("tcp", "localhost:1337")
	if err != nil {
		t.Fatal(err)
	} else {
		defer l.Close()
	}

	conn, err := net.Dial("tcp", "localhost:1337")
	if err != nil {
		t.Fatal(err)
	} else {
		defer conn.Close()
	}

	c := NewConn(conn)
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

		if !bytes.Equal(dst, []byte{1, 2, 3, 4}) {
			t.Fatal("wrong value")
		}

		w.Done()
	}()

	go func() {
		err := c.Write([]byte{1, 2, 3, 4})
		if err != nil {
			t.Fatal(err)
		}
	}()

	w.Wait()
}

func BenchmarkConn(b *testing.B) {
	l, err := net.Listen("tcp", "localhost:1337")
	if err != nil {
		b.Fatal(err)
	} else {
		defer l.Close()
	}

	conn, err := net.Dial("tcp", "localhost:1337")
	if err != nil {
		b.Fatal(err)
	} else {
		defer conn.Close()
	}

	c := NewConn(conn)
	var w sync.WaitGroup
	w.Add(b.N)

	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				b.Fatal(err)
			}

			d := NewConn(conn)

			for {
				if _, err := d.Read(); err == nil {
					w.Done()
				}
			}
		}
	}()

	go func() {
		src := make([]byte, 1024)
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				c.Write(src)
			}
		})
	}()

	b.ResetTimer()
	w.Wait()
}
