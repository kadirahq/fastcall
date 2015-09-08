package fastcall

import "net"

// Server listens for new connections
type Server struct {
	recv chan *Conn
	lsnr net.Listener
}

// Serve creates a listener and accepts connections
func Serve(addr string) (s *Server, err error) {
	lsnr, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	ch := make(chan *Conn)
	s = &Server{recv: ch, lsnr: lsnr}
	go s.accept()

	return s, nil
}

// Close stops accepting connections
func (s *Server) Close() (err error) {
	if err := s.lsnr.Close(); err != nil {
		return err
	}

	return nil
}

// Accept returns a channel of connections
func (s *Server) Accept() (ch chan *Conn) {
	return s.recv
}

func (s *Server) accept() {
	for {
		conn, err := s.lsnr.Accept()
		if err != nil {
			break
		}

		s.recv <- &Conn{conn: conn}
	}

	close(s.recv)
}
