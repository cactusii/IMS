package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip string
	Port int
}

func NewServer (ip string, port int) *Server {
	return &Server{
		ip,
		port,
	}
}

func (server *Server) Handle(conn net.Conn) {
	defer conn.Close()
	fmt.Println("连接建立成功！")
}

func (server *Server) Start() {
	// socket listen
	listen, err := net.Listen("tcp4", fmt.Sprint("%s:%d", server.Ip, server.Port))
	if err != nil {
		fmt.Println(err)
		return
	}
	// close
	defer listen.Close()

	// accept
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		// handle
		go server.Handle(conn)
	}

}