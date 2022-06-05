package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

var addr = net.UDPAddr{
	Port: 2000,
	IP:   net.ParseIP("192.168.1.240"),
}

func main() {

	server, err := net.ListenUDP("udp", &addr)

	if err != nil {
		log.Fatalf("err: %v", err)
	}

	for {
		mes := make([]byte, 2048)
		_, clientAddr, err := server.ReadFromUDP(mes)

		if err != nil {
			log.Fatalf("err: %v", err)
		}

		fmt.Printf("read message from (%v): %s \n", clientAddr, mes)
		go res(server, clientAddr, string(mes))
	}

}

func res(server *net.UDPConn, clientAddr *net.UDPAddr, mes string) {

	_, err := server.WriteToUDP([]byte(strings.ToUpper(mes)), clientAddr)
	if err != nil {
		log.Fatal("cannot send response")
	}

}
