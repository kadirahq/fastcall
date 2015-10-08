package main

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/kadirahq/fastcall"
)

func main() {
	var counter uint64

	conn, err := fastcall.DialBuf("localhost:1337")
	if err != nil {
		panic(err)
	} else {
		defer conn.Close()
	}

	go func() {
		for {
			time.Sleep(time.Second)
			fmt.Println(atomic.LoadUint64(&counter))
			atomic.StoreUint64(&counter, 0)
		}
	}()

	echo := func() {
		pld := make([]byte, 1024)

		for {
			if err := conn.Write(pld); err != nil {
				break
			}

			if _, err := conn.Read(); err != nil {
				break
			}

			if atomic.AddUint64(&counter, 1)%1000 == 0 {
				conn.FlushWriter()
			}
		}
	}

	for i := 0; i < 5000; i++ {
		go echo()
	}

	select {}
}
