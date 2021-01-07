package web

import (
	"net"
	"net/http"
	"strconv"
)

type ServerOptions struct {
	Socket  string
	Address string
	Port    int
}

func (o *ServerOptions) Server() (*Server, error) {
	var network string
	var address string
	if o.Socket != "" {
		network = "unix"
		address = o.Socket
	} else {
		network = "tcp"
		address = net.JoinHostPort(o.Address, strconv.Itoa(o.Port))
	}

	listen, err := net.Listen(network, address)
	if err != nil {
		return nil, err
	}

	server := &Server{
		listener: listen,
	}
	return server, nil
}

type Server struct {
	listener net.Listener
}

func (s *Server) Network() string {
	return s.listener.Addr().Network()
}

func (s *Server) String() string {
	return s.listener.Addr().String()
}

func (s *Server) Run(handler http.Handler) error {
	srv := &http.Server{
		Addr:    s.listener.Addr().String(),
		Handler: handler,
	}

	return srv.Serve(s.listener)
}
