package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func main() {
	//	与服务端建立连接
	conn, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		fmt.Println("err0 :", err)
		return
	}

	go sendMsd(conn)
	buf := make([]byte, 1024)
	for true {
		n, err := conn.Read(buf)
		if err == io.EOF {
			return
		} else if err != nil {
			fmt.Println("err1 :", err)
			return
		}
		fmt.Println(string(buf[:n]))

	}
}

func sendMsd(conn net.Conn) {
	var input string
	for true {
		reader := bufio.NewReader(os.Stdin)
		data, _, _ := reader.ReadLine()
		input = string(data)
		if strings.ToUpper(input) == "Q" {
			return
		}
		_, err := conn.Write([]byte(input))
		if err != nil {
			fmt.Println("err2 :", err)
			return
		}
	}
}
