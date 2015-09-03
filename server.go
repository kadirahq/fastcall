package fastcall

import (
	"net"

	"github.com/kadirahq/go-tools/logger"
)

// Handler handles all method calls
type Handler func(req []byte) (res []byte)

// Server responds to method calls
type Server struct {
	address  string
	listener net.Listener
	handler  Handler
}

// NewServer creates a server
func NewServer(addr string, fn Handler) (s *Server) {
	return &Server{address: addr, handler: fn}
}

// Listen starts listening
func (s *Server) Listen() (err error) {
	s.listener, err = net.Listen("tcp", s.address)
	if err != nil {
		return err
	}

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			logger.Print("fastcall", "error:", err)
			break
		}

		go s.handleConn(conn)
	}

	return nil
}

// Close stops the listener
func (s *Server) Close() (err error) {
	return s.listener.Close()
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()
	c := NewConn(conn)
	p := []byte{}

	for {
		in, err := c.Read()
		if err != nil {
			logger.Print("fastcall", "error:", err)
			break
		}

		sz := len(in)
		if sz < 4 {
			logger.Print("fastcall", "error:", "invalid id")
			break
		}

		cid := in[:4]
		req := in[4:]
		res := s.handler(req)

		if sz+4 > len(p) {
			p = make([]byte, sz+4)
		}

		out := p[:0]
		out = append(out, cid...)
		out = append(out, res...)

		if err := c.Write(out); err != nil {
			logger.Print("fastcall", "error:", err)
			break
		}
	}
}
