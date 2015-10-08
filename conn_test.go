package fastcall

import (
	"bytes"
	"net"
	"sync"
	"sync/atomic"
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

	c := New(conn)
	var w sync.WaitGroup
	w.Add(1)

	go func() {
		conn, err := l.Accept()
		if err != nil {
			t.Fatal(err)
		}

		d := New(conn)
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

	c := New(conn)
	var w sync.WaitGroup
	w.Add(b.N)

	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				b.Fatal(err)
			}

			d := New(conn)

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

func BenchmarkConnBuff(b *testing.B) {
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

	c := NewBuf(conn)
	var w sync.WaitGroup
	w.Add(b.N)

	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				b.Fatal(err)
			}

			d := NewBuf(conn)

			for {
				if _, err := d.Read(); err == nil {
					w.Done()
				}
			}
		}
	}()

	var counter uint64

	go func() {
		src := make([]byte, 1024)
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				c.Write(src)
				if atomic.AddUint64(&counter, 1)%1000 == 0 {
					c.FlushWriter()
				}
			}

			c.FlushWriter()
		})
	}()

	b.ResetTimer()
	w.Wait()
}
