package main

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	Ip string
	Port int
	OnlineMap map[string]*User
	mapLock sync.RWMutex

	Message chan string
}

func NewServer (ip string, port int) *Server {
	return &Server{
		Ip: ip,
		Port: port,
		OnlineMap: make(map[string]*User),
		Message: make(chan string),
	}
}

func (server *Server) ListenMsg() {
	for {
		msg := <- server.Message
		server.mapLock.Lock()
		for _, user := range server.OnlineMap {
			user.c <- msg
		}
		server.mapLock.Unlock()
	}
}

func (server *Server) Broadcast(user *User, msg string) {
	sendMsg := "[" + user.Name + "]" + user.Addr + ":" + msg
	server.Message <- sendMsg
}

func (server *Server) Handle(conn net.Conn) {
	defer conn.Close()
	user := NewUser(conn.RemoteAddr().String(), conn)

	server.mapLock.Lock()
	server.OnlineMap[conn.RemoteAddr().String()] = user
	server.mapLock.Unlock()

	server.Broadcast(user, "上线！")

	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				server.Broadcast(user, "下线")
				return
			}
			if err != nil {
				fmt.Println("conn read err:", err)
				return
			}
			msg := string(buf[:n-1])
			server.Broadcast(user, msg)
		}
	}()

	select {

	}
}

func (server *Server) Start() {
	// socket 监听
	listen, err := net.Listen("tcp4", fmt.Sprintf("%s:%d", server.Ip, server.Port))
	if err != nil {
		fmt.Println(err)
		return
	}
	// close
	defer listen.Close()

	// 监听用户的Massage
	go server.ListenMsg()

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