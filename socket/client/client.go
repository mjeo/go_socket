package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	//	与服务端建立连接
	conn, err := net.Dial("tcp", "127.0.0.1:20000")
	if err != nil {
		fmt.Println("err :", err)
		return
	}
	//	利用该连接进行数据的发送和接受

	input := bufio.NewReader(os.Stdin) ////创建从当前连接读数据的对象 从终端读取用户输入
	for {
		s, _ := input.ReadString('\n') //读到\n发送
		s = strings.TrimSpace(s)
		if strings.ToUpper(s) == "Q" {
			return
		}
		//	给服务端发消息
		_, err := conn.Write([]byte(s))
		if err != nil {
			return
		}
		//	从服务端接收回复的消息
		var buf [1024]byte
		n, err := conn.Read(buf[:]) //返回读的数量和err
		if err != nil {
			fmt.Println("recv failed, err:", err)
			return
		}
		fmt.Println("收到服务端回复", string(buf[:n]))
	}
}
