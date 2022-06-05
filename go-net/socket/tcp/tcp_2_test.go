package tcp

import (
	"io"
	"log"
	"net"
	"testing"
)

func TestReadToFixedBuffer(t *testing.T) {

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatal(err)
	}
	buf := make([]byte, 1<<19) // 512 KB

	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Println(err)
			}
			break
		}

		log.Printf("read %d bytes", n) // buf[:n] is the data read from conn
	}

}
