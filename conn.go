package fastcall

import (
	"encoding/binary"
	"net"
	"sync"
)

// Conn sends and receives binary messages over tcp
type Conn struct {
	conn net.Conn
	rmtx sync.Mutex
	wmtx sync.Mutex
}

// Dial creates a connection to given address
func Dial(addr string) (c *Conn, err error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	c = &Conn{conn: conn}
	return c, err
}

// Read reads a message from connection
func (c *Conn) Read() (b []byte, err error) {
	c.rmtx.Lock()
	defer c.rmtx.Unlock()

	var sz int
	var szu32 uint32

	if err = binary.Read(c.conn, binary.LittleEndian, &szu32); err != nil {
		return nil, err
	}

	sz = int(szu32)
	b = make([]byte, sz)

	buffer := b[:sz]
	toRead := buffer[:]
	for len(toRead) > 0 {
		n, err := c.conn.Read(toRead)
		if err != nil {
			return nil, err
		}

		toRead = toRead[n:]
	}

	return buffer, nil
}

// Write writes a message to the connection
func (c *Conn) Write(b []byte) (err error) {
	c.wmtx.Lock()
	defer c.wmtx.Unlock()

	sz := len(b)
	szu32 := uint32(sz)

	buffer := b[:sz]
	if err := binary.Write(c.conn, binary.LittleEndian, szu32); err != nil {
		return err
	}

	toWrite := buffer[:]
	for len(toWrite) > 0 {
		n, err := c.conn.Write(toWrite)
		if err != nil {
			return err
		}

		toWrite = toWrite[n:]
	}

	return nil
}

// Close closes the connection
func (c *Conn) Close() (err error) {
	c.rmtx.Lock()
	defer c.rmtx.Unlock()

	c.wmtx.Lock()
	defer c.wmtx.Unlock()

	if err := c.conn.Close(); err != nil {
		return err
	}

	return nil
}
