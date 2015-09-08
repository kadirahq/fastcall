package fastcall

import (
	"bytes"
	"sync"
	"testing"
)

func TestServer(t *testing.T) {
	serv, err := Serve("localhost:1337")
	if err != nil {
		t.Fatal(err)
	} else {
		defer serv.Close()
	}

	conn, err := Dial("localhost:1337")
	if err != nil {
		t.Fatal(err)
	} else {
		defer conn.Close()
	}

	var w sync.WaitGroup
	w.Add(2)

	go func() {
		conn := <-serv.Accept()

		if msg, err := conn.Read(); err != nil {
			t.Fatal(err)
		} else if !bytes.Equal(msg, []byte{1, 2, 3, 4}) {
			t.Fatal("wrong value")
		}

		if err := conn.Write([]byte{5, 6, 7, 8}); err != nil {
			t.Fatal(err)
		}

		w.Done()
	}()

	go func() {
		if err := conn.Write([]byte{1, 2, 3, 4}); err != nil {
			t.Fatal(err)
		}

		if msg, err := conn.Read(); err != nil {
			t.Fatal(err)
		} else if !bytes.Equal(msg, []byte{5, 6, 7, 8}) {
			t.Fatal("wrong value")
		}

		w.Done()
	}()

	w.Wait()
}
