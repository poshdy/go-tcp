package main

import (
	"fmt"
	"log"
	"net"
)

type Message struct {
	sender  string
	content []byte
}

type Server struct {
	Addr     string
	ln       net.Listener
	quitchan chan struct{}
	msgchan  chan Message
}

func NewServer(addr string) *Server {
	return &Server{
		Addr:     addr,
		quitchan: make(chan struct{}),
		msgchan:  make(chan Message, 10),
	}
}

func (s *Server) Start() error {

	ln, e := net.Listen("tcp", s.Addr)
	if e != nil {
		fmt.Println("Listennig ERR", e)
	}
	defer ln.Close()
	s.ln = ln
	go s.acceptLoop()

	<-s.quitchan
	close(s.msgchan)
	return nil

}

func (s *Server) acceptLoop() {

	for {
		conn, err := s.ln.Accept()

		if err != nil {
			fmt.Println("Accept Error", err)
			continue
		}
		fmt.Println("New Connection to server", conn.RemoteAddr())
		go s.readLoop(conn)
	}
}
func (s *Server) readLoop(conn net.Conn) {

	defer conn.Close()
	buff := make([]byte, 3000)
	for {
		n, err := conn.Read(buff)
		if err != nil {
			fmt.Println("READ ERROR", err)
			continue
		}

		s.msgchan <- Message{
			sender:  conn.RemoteAddr().String(),
			content: buff[:n],
		}
	}
}

func main() {
	s := NewServer(":3000")
	go func() {
		for msg := range s.msgchan {
			fmt.Println("Message from:", msg.sender, "content: ", string(msg.content))
		}
	}()
	log.Fatal(s.Start())

}
