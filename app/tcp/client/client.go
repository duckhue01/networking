package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

var addr = net.TCPAddr{
	Port: 2222,
	IP:   net.ParseIP("192.168.1.240"),
}

func main() {

	// connect to this socket
	conn, err := net.DialTCP("tcp", nil, &addr)

	if err != nil {
		log.Fatal(err)
	}
	// read in input from stdin
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Text to send: ")
	text, _ := reader.ReadString('\n')
	// send to socket
	// fmt.Fprintf(conn, text+"\n")
	conn.Write([]byte(text))
	// listen for reply

	message, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Print("Message from server: " + message)
}
