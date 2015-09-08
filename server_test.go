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
	w.Add(1)

	go func() {
		conn := <-serv.Accept()
		dst, err := conn.Read()
		if err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(dst, []byte{1, 2, 3, 4}) {
			t.Fatal("wrong value")
		}

		w.Done()
	}()

	go func() {
		err := conn.Write([]byte{1, 2, 3, 4})
		if err != nil {
			t.Fatal(err)
		}
	}()

	w.Wait()
}
