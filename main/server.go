package main

import (
	"fmt"
	"net"
	"sync"
	"time"
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
	user := NewUser(conn.RemoteAddr().String(), conn, server)
	user.Online()

	isLive := make(chan bool)

	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}
			if err != nil {
				fmt.Println("conn read err:", err)
				return
			}
			msg := string(buf[:n-1])

			user.DoMsg(msg)

			// 每次用户发消息，会向isLive管道中发送消息
			isLive <- true
		}
	}()

	for {
		select {
		// 当isLive出现消息，会执行isLive的case，
		// 之后本次的select会结束，进入下一次for循环，time.After重置，
		// 这样去保证一个超时机制
		case <- isLive:
		case <- time.After(time.Second * 10):
			// 强制下线
			user.Offline()
			return

		}
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