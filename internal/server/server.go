package server

import (
	"fmt"
	"log"
	"net"
)

type Server struct {
	addr string
	port string
}

type SetServer func(*Server)

func WithPort(port string) SetServer {
	return func(s *Server) {
		s.port = port
	}
}

func WithAddr(addr string) SetServer {
	return func(s *Server) {
		s.addr = addr
	}
}

func New(options ...SetServer) *Server {
	s := &Server{
		addr: "localhost",
		port: "6380",
	}
	for _, set := range options {
		set(s)
	}
	return s;
}

func (s *Server) Start() error {
    listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", s.addr, s.port))
    if err != nil {
        return err
    }
    defer listener.Close()

    log.Println("Server started on port "+fmt.Sprintf("%s:%s", s.addr, s.port))
    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Printf("Failed to accept connection: %v", err)
            continue
        }
        go s.handleConnection(conn)
    }
}

func (s *Server) handleConnection(conn net.Conn) {
    defer conn.Close()
    handler := NewRESPHandler(conn)
    for {
        if err := handler.Handle(); err != nil {
            log.Printf("Error handling connection: %v", err)
            return
        }
    }
}