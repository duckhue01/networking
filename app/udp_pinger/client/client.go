package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

var addr = net.UDPAddr{
	Port: 2000,
	IP:   net.ParseIP("192.168.1.240"),
}

func main() {

	udp, err := net.DialUDP("udp", nil, &addr)

	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		mes := make([]byte, 1024)

		start := time.Now().UnixNano()
		udp.Write([]byte("ping"))
		_, err := udp.Read(mes)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(mes), time.Now().UnixNano() - start)
		time.Sleep(1 * time.Second)
		
	}

}
