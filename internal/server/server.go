package server

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

type RedisObject interface {
	Type() string
	LastAccess() time.Time
	Touch()
	String() string
	Len() int
}

type Server struct {
	addr string
	port string
	store map[string]RedisObject
	mu sync.Mutex
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
		store: make(map[string]RedisObject),
	}
	for _, set := range options {
		set(s)
	}
	return s;
}

func (s *Server) Set(key string, value RedisObject) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.store[key] = value
}

func (s *Server) Get(key string) RedisObject {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.store[key]
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
    handler := NewRESPHandler(conn, s)
    for {
        if err := handler.Handle(); err != nil {
            log.Printf("Error handling connection: %v", err)
            return
        }
    }
}