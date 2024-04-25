package main

import (
	"fmt"
	"log"
	"net"
)

type Server struct {
	listenAddr string
	ln         net.Listener
	clients    map[*Client]struct{}
	join       chan net.Conn
	leave      chan net.Conn
	broadcast  chan *Message
}

func NewServer(listenAddr string) *Server {
	return &Server{
		listenAddr: listenAddr,
		clients:    make(map[*Client]struct{}),
		join:       make(chan net.Conn),
		leave:      make(chan net.Conn),
		broadcast:  make(chan *Message),
	}
}

func (s *Server) start() error {
	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()

	s.ln = ln

	go s.manageConnections()
	s.acceptLoop()

	return nil
}

func (s *Server) manageConnections() {
	for {
		select {
		case conn := <-s.join:
			log.Printf("client joined: %s\n", conn.RemoteAddr())
			s.listConnections()
		case conn := <-s.leave:
			log.Printf("client left: %s\n", conn.RemoteAddr())
		case message := <-s.broadcast:
			for client := range s.clients {
				if _, ok := client.rooms[message.room]; ok {
					client.writeMessage(message)
				}
			}
		}
	}
}

func (s *Server) acceptLoop() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	client := NewClient(conn)

	s.clients[client] = struct{}{}
	s.join <- client.conn

	defer func() {
		s.leave <- client.conn
		delete(s.clients, client)
		client.conn.Close()
	}()

	client.loop(s)
}

func (s *Server) listConnections() {
	i := 0
	for client := range s.clients {
		fmt.Printf("%d :: %s\n", i, client.conn.RemoteAddr().String())
		i++
	}
}
