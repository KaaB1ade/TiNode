package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int
}

func NewClient(serverIp string, serverPort int) *Client {
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       999,
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("Dial err: ", err)
		return nil
	}

	client.conn = conn

	return client
}

var serverIp string
var serverPort int

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "SET serverIp to 127.0.0.1(default)")
	flag.IntVar(&serverPort, "port", 8888, "SET serverPort to 8888(default)")
}

func (client *Client) menu() bool {
	var flag int

	fmt.Println("1.public chat")
	fmt.Println("2.private chat")
	fmt.Println("3.change name")
	fmt.Println("0.log out")

	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println(">>>>input legal num<<<<")
		return false
	}
}

func (client *Client) DealResponse() {
	io.Copy(os.Stdout, client.conn)
}

func (client *Client) SelectUser() {
	sendMsg := "who\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn err:", err)
		return
	}
}

func (client *Client) PrivateChat() {
	var remoteUser string
	var chatMsg string

	client.SelectUser()
	fmt.Println(">>>select ur user exit<<<")
	fmt.Scanln(&remoteUser)

	for remoteUser != "exit" {
		fmt.Println(">>>ur message exit<<<")
		fmt.Scanln(&chatMsg)

		for chatMsg != "exit" {
			if len(chatMsg) != 0 {
				sendMsg := "to|" + remoteUser + "|" + chatMsg + "\n\n"
				_, err := client.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn write err", err)
					break
				}
			}
			chatMsg = ""
			fmt.Println(">>>input message exit<<<")
			fmt.Scanln(&chatMsg)
		}

		client.SelectUser()
		fmt.Println(">>>select ur user exit<<<")
		fmt.Scanln(&remoteUser)
	}

}

func (client *Client) UpdateName() bool {
	fmt.Println(">>>input user name<<<")
	fmt.Scanln(client.Name)

	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn write err:", err)
		return false
	}
	return true
}

func (client *Client) PublicChat() {
	var chatMessage string

	fmt.Println(">>>input message exit<<<")
	fmt.Scanln(&chatMessage)

	for chatMessage != "exit" {
		if len(chatMessage) != 0 {
			sendMessage := chatMessage + "\n"
			_, err := client.conn.Write([]byte(sendMessage))
			if err != nil {
				fmt.Println("send err:", err)
				break
			}
		}

		chatMessage = ""
		fmt.Println(">>>input message exit<<<")
		fmt.Scanln(&chatMessage)
	}
}

func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {
		}

		switch client.flag {
		case 1:
			client.PublicChat()
			break
		case 2:
			client.PrivateChat()
			break
		case 3:
			client.UpdateName()
			break
		}
	}
}

func main() {
	//命令行解析
	flag.Parse()

	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println("链接失败...")
		return
	}

	go client.DealResponse()

	fmt.Println("链接成功...")

	client.Run()
}
