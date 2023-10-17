package main

import (
	"net"
	"strings"
)

type User struct {
	Name   string
	Addr   string
	C      chan string
	conn   net.Conn
	server *Server
}

func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()

	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}

	go user.ListenMessage()

	return user
}

func (this *User) Online() {
	this.server.maplock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.maplock.Unlock()

	this.server.BroadCast(this, "online")
}

func (this *User) Offline() {
	this.server.maplock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.maplock.Unlock()

	this.server.BroadCast(this, "offline")
}

func (this *User) SendMsg(msg string) {
	this.conn.Write([]byte(msg))
}

func (this *User) DoMessage(msg string) {
	if msg == "who" {
		this.server.maplock.Lock()
		for _, user := range this.server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + "is online...\n"
			this.SendMsg(onlineMsg)
		}
		this.server.maplock.Unlock()

	} else if len(msg) > 7 && msg[:7] == "rename|" {
		newName := strings.Split(msg, "|")[1]

		_, ok := this.server.OnlineMap[newName]
		if ok {
			this.SendMsg("current user name has been used...\n")
		} else {
			this.server.maplock.Lock()
			delete(this.server.OnlineMap, this.Name)
			this.server.OnlineMap[newName] = this
			this.server.maplock.Unlock()

			this.Name = newName
			this.SendMsg("ur user name has been set to " + newName + "\n")
		}
	} else if len(msg) > 4 && msg[:3] == "to|" {
		remoteName := strings.Split(msg, "|")[1]
		if remoteName == "" {
			this.SendMsg("wrong input...\n")
			return
		}

		remoteUser, ok := this.server.OnlineMap[remoteName]
		if !ok {
			this.SendMsg("no such user...\n")
			return
		}

		content := strings.Split(msg, "|")[2]
		if content == "" {
			this.SendMsg("no message...\n")
			return
		}
		remoteUser.SendMsg(this.Name + ": " + content + "[private]")

	} else {
		this.server.BroadCast(this, msg)
	}
}

func (this *User) ListenMessage() {
	for {
		msg := <-this.C

		this.conn.Write([]byte(msg + "\n"))
	}
}
