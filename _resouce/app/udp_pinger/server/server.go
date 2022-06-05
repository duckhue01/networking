package main

import (
	"log"
	"net"
)

var addr = net.UDPAddr{
	Port: 2000,
	IP:   net.ParseIP("192.168.1.240"),
}

func main() {

	udp, err := net.ListenUDP("udp", &addr)

	if err != nil {
		log.Fatal(err)
	}

	for {
		_, clientAddr, err := udp.ReadFromUDP(nil)

		if err != nil {
			log.Fatal(err)
		}

		go func(clientAddr *net.UDPAddr) {
			_, err := udp.WriteToUDP(nil, clientAddr)
			if err != nil {
				log.Fatal(err)
			}
		}(clientAddr)
	}

}
