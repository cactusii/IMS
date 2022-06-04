package main

import (
	"net"
	"strings"
)

type User struct {
	Name string
	Addr string
	c chan string
	conn net.Conn

	server *Server
}

func NewUser(name string, conn net.Conn, server *Server) *User {
	user := &User{
		name,
		conn.RemoteAddr().String(),
		make(chan string),
		conn,
		server,
	}
	go user.ListenMeg()
	return user
}

func (user *User) Online() {
	user.server.mapLock.Lock()
	user.server.OnlineMap[user.Name] = user
	user.server.mapLock.Unlock()

	user.server.Broadcast(user, "上线！")
}

func (user *User) sendMsg(msg string) {
	user.conn.Write([]byte(msg))
}

func (user *User) DoMsg(msg string) {
	if msg == "who" {
		user.server.mapLock.Lock()
		for _, user := range user.server.OnlineMap {
			onlineMsg := user.Addr + "\n"
			user.sendMsg(onlineMsg)
		}
		user.server.mapLock.Unlock()
	} else if msg[:6] == "rename" {
		newName := strings.Split(msg, " ")[1]
		_, ok := user.server.OnlineMap[newName]
		if ok {
			user.sendMsg(newName + "已存在！")
		} else {
			user.server.mapLock.Lock()
			delete(user.server.OnlineMap, user.Name)
			user.server.OnlineMap[newName] = user
			user.server.mapLock.Unlock()

			user.Name = newName
			user.sendMsg("更新成功！")
		}
	} else {
		user.server.Broadcast(user, msg)
	}
}

func (user *User) Offline() {
	user.server.Broadcast(user, "下线！")

	user.server.mapLock.Lock()
	delete(user.server.OnlineMap, user.Name)
	user.conn.Close()
	user.server.mapLock.Unlock()

}


func (user *User) ListenMeg() {
	for {
		msg := <-user.c
		user.conn.Write([]byte(msg + "\n"))
	}
}

