package main

import "net"

type User struct {
	Name string
	Addr string
	c chan string
	conn net.Conn
}

func NewUser(name string, conn net.Conn) *User {
	user := &User{
		name,
		conn.RemoteAddr().String(),
		make(chan string),
		conn,
	}
	go user.ListenMeg()
	return user
}

func (user *User) ListenMeg() {
	for {
		msg := <-user.c
		user.conn.Write([]byte(msg + "\n"))
	}
}

