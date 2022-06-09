package tcp

import (
	"bufio"
	"io"
	"log"
	"net"
	"reflect"
	"sync"
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

func TestScanner(t *testing.T) {
	con, err := net.Dial(tcp, addr)
	if err != nil {
		t.Fatal(err)
	}
	defer con.Close()
	scanner := bufio.NewScanner(con)
	scanner.Split(bufio.ScanWords)

	words := make([]string, 0)

	for scanner.Scan() {
		words = append(words, scanner.Text())
	}

	err = scanner.Err()
	if err != nil {
		t.Error(err)
	}
	expected := []string{"The", "bigger", "the", "interface,", "the",
		"weaker", "the", "abstraction."}
	if !reflect.DeepEqual(words, expected) {
		t.Fatal("inaccurate scanned word list")
	}

	t.Logf("Scanned words: %#v", words)
}

func Test_PayloadListener(t *testing.T) {
	conn, err := net.Dial(tcp, addr)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	for i := 0; i < len(payloads); i++ {
		actual, err := decode(conn)
		if err != nil {
			t.Fatal(err)
		}
		if expected := payloads[i]; !reflect.DeepEqual(expected, actual) {
			t.Errorf("value mismatch: %v != %v", expected, actual)
			continue
		}
		t.Logf("[%T] %[1]q", actual)
	}
}

func Test_proxy(t *testing.T) {
	var wg sync.WaitGroup
	proxyServer, err := net.Listen(tcp, "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			conn, err := proxyServer.Accept()
			if err != nil {
				return
			}

			go func(from net.Conn) {
				defer from.Close()
				to, err := net.Dial(tcp, addr)
				if err != nil {
					t.Error(err)
					return
				}
				defer to.Close()
				err = proxy(from, to)
				if err != nil && err != io.EOF {
					t.Error(err)
				}
			}(conn)

		}
	}()

	conn, err := net.Dial(tcp, proxyServer.Addr().String())
	if err != nil {
		t.Fatal(err)
	}

	msgs := []struct{ Message, Reply string }{{"ping", "pong"},
		{"pong", "pong"},
		{"echo", "echo"},
		{"ping", "pong"},
	}
	for i, m := range msgs {
		_, err = conn.Write([]byte(m.Message))
		if err != nil {
			t.Fatal(err)
		}
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			t.Fatal(err)
		}
		actual := string(buf[:n])
		t.Logf("%q -> proxy -> %q", m.Message, actual)
		if actual != m.Reply {
			t.Errorf("%d: expected reply: %q; actual: %q",
				i, m.Reply, actual)
		}
	}
	_ = conn.Close()
	_ = proxyServer.Close()

	wg.Wait()

}
