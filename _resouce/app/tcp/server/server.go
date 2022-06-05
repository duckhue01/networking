package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

var addr = net.TCPAddr{
	Port: 2222,
	// IP:   net.ParseIP("195.106.1.240"),
}

func main() {
	// var mess []byte
	fmt.Println("Launching server...")

	// listen on all interfaces
	server, err := net.ListenTCP("tcp", &addr)

	if err != nil {
		log.Fatalf("err1: %v", err)
	}
	defer server.Close()

	for {
		conn, _ := server.Accept()
		go func(c net.Conn) {
			// will listen for message to process ending in newline (\n)
			message, _ := bufio.NewReader(c).ReadString('\n')
			// output message received
			fmt.Print("Message Received:", string(message))
			// sample process for string received
			newmessage := strings.ToUpper(message)
			// send new string back to client
			c.Write([]byte(newmessage + "\n"))
			c.Close()
		}(conn)

	}

}

    
