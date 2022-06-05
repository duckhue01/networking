package tcp

import (
	"crypto/rand"
	"log"
	"net"
)

func ReadToFixedBuffer() {
	listener, err := net.Listen(tcp, addr)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	payload := make([]byte, 1<<24) // 16MB

	_, err = rand.Read(payload)
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Panic(err)
		}

		go func(c net.Conn) {
			defer conn.Close()
			log.Println("payload length: ", len(payload))

			n, err := conn.Write(payload)
			if err != nil {
				log.Fatal(err)
			}
			log.Println("sent payload: ", n)
		}(conn)
	}
}

func Scanner() {

}
