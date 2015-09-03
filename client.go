package fastcall

import (
	"bytes"
	"encoding/binary"
	"errors"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/kadirahq/go-tools/secure"
)

// Client calls methods
type Client struct {
	Conn

	address  string
	active   *secure.Bool
	closed   *secure.Bool
	connmtx  sync.Mutex
	inflight map[uint32]chan []byte
	callID   uint32
}

// NewClient creates a client
func NewClient(addr string) (c *Client) {
	c = &Client{
		address:  addr,
		active:   secure.NewBool(true),
		closed:   secure.NewBool(true),
		inflight: make(map[uint32]chan []byte),
	}

	return c
}

// Call calls a method
func (c *Client) Call(p []byte) (res []byte, err error) {
	id := atomic.AddUint32(&c.callID, 1)
	ch := make(chan []byte)
	c.inflight[id] = ch

	// ---- OPTIMIZE !!!
	buff := bytes.NewBuffer(nil)
	if err := binary.Write(buff, binary.LittleEndian, id); err != nil {
		return nil, errors.New("encode")
	}
	buff.Write(p)
	data := buff.Bytes()
	// ----

	c.connmtx.Lock()
	if c.closed.Get() {
		c.connmtx.Unlock()
		return nil, errors.New("closed")
	}

	if err := c.Write(data); err != nil {
		c.conn.Close()
		c.connmtx.Unlock()
		return nil, err
	}
	c.connmtx.Unlock()

	out := <-ch
	if out == nil {
		return nil, errors.New("failed")
	}

	return out, nil
}

// Connect connects
func (c *Client) Connect() {
	var err error

	// stop only if the connection is closed by the user
	// if closed for other reasons, the loop will continue
	if !c.active.Get() {
		return
	}

	// STEP 1/3: connect
	// create a new tcp connection
	c.connmtx.Lock()

	for {
		c.conn, err = net.DialTimeout("tcp", c.address, time.Minute)
		if err != nil {
			time.Sleep(time.Minute)
			continue
		}

		break
	}

	c.closed.Set(false)
	c.connmtx.Unlock()

	// STEP 2/3: read/write until error
	// read and handle response messages
	for {
		in, err := c.Read()
		if err != nil {
			c.closed.Set(true)
			break
		}

		sz := len(in)
		if sz < 4 {
			c.closed.Set(true)
			break
		}

		cid := in[:4]
		msg := in[4:]
		var id uint32

		// ---- OPTIMIZE !!!
		buff := bytes.NewBuffer(cid)
		if err := binary.Read(buff, binary.LittleEndian, &id); err != nil {
			c.closed.Set(true)
			break
		}
		// ----

		if ch, ok := c.inflight[id]; ok {
			delete(c.inflight, id)
			ch <- msg
		}
	}

	// STEP 3/3: disconnect
	// end inflight method calls
	// respond with ErrDisconnect
	for id, ch := range c.inflight {
		delete(c.inflight, id)
		ch <- nil
	}
}

// Close closes the connection
func (c *Client) Close() (err error) {
	c.active.Set(false)
	return nil
}
