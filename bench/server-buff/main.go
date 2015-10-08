package main

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/kadirahq/fastcall"
)

func main() {
	var counter uint64

	serv, err := fastcall.Serve("localhost:1337")
	if err != nil {
		panic(err)
	} else {
		defer serv.Close()
	}

	go func() {
		for {
			time.Sleep(time.Second)
			fmt.Println(atomic.LoadUint64(&counter))
			atomic.StoreUint64(&counter, 0)
		}
	}()

	for {
		conn, err := serv.AcceptBuf()
		if err != nil {
			panic(err)
		}

		for {
			msg, err := conn.Read()
			if err != nil {
				break
			}

			if err := conn.Write(msg); err != nil {
				break
			}

			if atomic.AddUint64(&counter, 1)%1000 == 0 {
				conn.FlushWriter()
			}
		}
	}
}
