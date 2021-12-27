package main

import (
	"bufio"
	"bytes"
	"io/ioutil"

	// "bytes"
	"fmt"
	// "io/ioutil"

	"log"
	"net"
	"net/http"
)

var addr = net.TCPAddr{
	Port: 3000,
}

func main() {

	fmt.Println("server is running...")

	server, err := net.ListenTCP("tcp", &addr)

	if err != nil {
		log.Fatal(err)
	}
	defer server.Close()

	for {
		conn, _ := server.Accept()

		go func(c net.Conn) {

			req, err := http.ReadRequest(bufio.NewReader(c))

			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(req.URL.Path)
			file, err := ioutil.ReadFile("." + req.URL.Path)
			if err != nil {
				eFile, _ := ioutil.ReadFile("./404.html")
				res := http.Response{
					Body:   ioutil.NopCloser(bytes.NewBuffer(eFile)),
					Status: "200 OK",
					Proto:  "HTTP/1.1",
				}
				res.Write(c)
				c.Close()
			}

			res := http.Response{
				Body:   ioutil.NopCloser(bytes.NewBuffer(file)),
				Status: "200 OK",
				Proto:  "HTTP/1.1",
			}
			res.Write(c)
			c.Close()
		}(conn)
	}

}
