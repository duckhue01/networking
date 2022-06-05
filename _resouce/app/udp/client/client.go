package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

var addr = net.UDPAddr{
	Port: 2000,
	IP:   net.ParseIP("192.168.1.240"),
}

func main() {
	mes := make([]byte, 2048)
	conn, err := net.DialUDP("udp", nil, &addr)
	if err != nil {
		log.Fatalf("err: %v", err)
		return
	}

	fmt.Fprintf(conn, "hi w are u do in udp ?")

	_, err = bufio.NewReader(conn).Read(mes)

	if err != nil {
		log.Fatalf("err: %v", err)
	}

	fmt.Printf("%s", mes)

}
