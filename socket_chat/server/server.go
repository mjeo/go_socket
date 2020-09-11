package main

import (
	"fmt"
	"net"
	"strings"
	"time"
)

type (
	// 用户类型
	client struct {
		c    chan string
		name string
		addr string
	}
)

var (
	onlineUser  = make(map[string]client) // 在线用户
	messageList = make(chan string)       // 消息列表
)

func whiteMsgToUser(user client, conn net.Conn) {
	for msg := range user.c {
		conn.Write([]byte(msg + "\n"))

	}
}

func makeMsg(user client, msg string) (buf string) {
	buf = fmt.Sprintf("[%s]{%s}:%s\n", user.addr, user.name, msg)
	return
}

func contentHandler(conn net.Conn) {
	defer conn.Close()
	//创建chan用于判断用户是否活跃
	activeStatus := make(chan bool)
	//获取用户网络地址ip+port
	addr := conn.RemoteAddr().String()
	//	创建新连接用户的结构体信息
	user := client{
		make(chan string),
		addr,
		addr,
	}
	//创建用来专门给用户发送消息的go程
	go whiteMsgToUser(user, conn)
	//	将新连接用户添加到在线用户列表
	onlineUser[addr] = user
	//	发送用户上线消息到全局消息
	messageList <- makeMsg(user, "log in!!")
	//创建chan 用来判断用户退出状态
	quitList := make(chan bool) // 用户退出状态
	//	创建匿名go程,处理用户发送的消息
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				quitList <- true
				fmt.Printf("member {%s} log out", user.name)
				return
			}
			if err != nil {
				fmt.Println("conn.Read err:", err)
				return
			}
			//	将读到的用户消息写入到消息列表
			msg := string(buf[:n])
			//msg = strings.Split(msg,"#")
			fmt.Println(msg)
			if msg == "online" {
				conn.Write([]byte("online user list\n"))
				for _, v := range onlineUser {
					userInfo := fmt.Sprintf("[%s]{%s}\n", v.addr, v.name)
					conn.Write([]byte(userInfo))
				}
			} else if len(msg) > 7 && msg[:6] == "rename" {
				name := msg[8:]
				oldname := user.name
				user.name = name
				onlineUser[user.addr] = user
				messageList <- makeMsg(user, oldname+" rename "+name)
			} else if len(msg) > 15 && msg[:3] == "to#" {
				//私信模式
				//格式 to#ip:port#用户消息
				content := strings.Split(msg, "#")
				msg = "form to " + user.name + "  :  " + content[2]
				toUserMsgList := onlineUser[content[1]].c
				toUserMsgList <- msg
			} else {
				messageList <- makeMsg(user, msg)
			}
			activeStatus <- true

		}
	}()

	for {
		//监听chan的数据流动
		select {
		case <-quitList:
			delete(onlineUser, user.addr) //将用户从在线列表移除
			messageList <- makeMsg(user, user.name+" log out")
			return
		case <-activeStatus:
		//	用户重置计时器
		case <-time.After(time.Second * 15):
			delete(onlineUser, user.addr) //将用户从在线列表移除
			messageList <- makeMsg(user, user.name+" time out")
			return
		}
	}
}

func manager() {
	//监听信息列表
	for {
		msg := <-messageList
		//	循环发送给所有在线用户
		for _, v := range onlineUser {
			v.c <- msg
		}
	}

}

func main() {
	//	创建监听socket
	listen, err := net.Listen("tcp", "127.0.0.1:8000")
	if err != nil {
		fmt.Println("Listen Err", err)
		return
	}
	defer listen.Close()

	//创建管理者go程,用于管理用户在线列表和监听消息列表
	go manager()
	//	循环监听客户端请求
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("Listen Err", err)
			return
		}
		//	启动go程处理客户端数据请求
		go contentHandler(conn)
	}
}
