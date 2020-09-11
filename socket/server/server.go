package main

import (
	"bufio"
	"fmt"
	"net"
)

var (
	connList []net.Conn
)

func main() {
	listen, err := net.Listen("tcp", "127.0.0.1:20000")
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		conn, err := listen.Accept()
		connList = append(connList,conn)
		//fmt.Println(connList)
		if err != nil {
			fmt.Println(err)
			continue
		}
		go proess(conn)

	}
}

func proess(conn net.Conn) {
	defer conn.Close()
	//对当前连接做数据的接|发操作
	for {
		reader := bufio.NewReader(conn) //创建从当前连接读数据的对象
		var buf [128]byte
		n, err := reader.Read(buf[:])
		if err != nil {
			fmt.Println("read from client failed, err:", err)
			break
		}
		recvStr := string(buf[:n])
		fmt.Println("收到client端发来的数据：", recvStr)
		//把接受到的数据返回到客户端
		conn.Write([]byte("ok"))

	}
}