package fastcall

import (
	"encoding/binary"
	"net"
	"sync"
)

// Conn does net stuff
// TODO write comment
type Conn struct {
	conn  net.Conn
	rmutx sync.Mutex
	wmutx sync.Mutex
}

// NewConn creates a connection
func NewConn(conn net.Conn) (c *Conn) {
	return &Conn{conn: conn}
}

// Read reads a message from connection
func (c *Conn) Read() (b []byte, err error) {
	c.rmutx.Lock()
	defer c.rmutx.Unlock()

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
	c.wmutx.Lock()
	defer c.wmutx.Unlock()

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
