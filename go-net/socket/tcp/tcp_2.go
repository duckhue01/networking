package tcp

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
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
	payload := "The bigger the interface, the weaker the abstraction."

	listener, err := net.Listen(tcp, addr)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	con, err := listener.Accept()
	if err != nil {
		log.Fatal(err)

	}
	defer con.Close()
	_, err = con.Write([]byte(payload))
	if err != nil {
		log.Fatal(err)
	}

}

const (
	BinaryType uint8 = iota + 1
	StringType
	MaxPayloadSize uint32 = 10 << 20 // 10 MB

)

var (
	ErrMaxPayloadSize = errors.New("maximum payload size exceeded")
	b1                = Binary("Clear is better than clever.")
	b2                = Binary("Don't panic.")
	s1                = String("Errors are values.")
	payloads          = []Payload{&b1, &s1, &b2}
)

type (
	Payload interface {
		fmt.Stringer
		io.ReaderFrom
		io.WriterTo
		Bytes() []byte
	}

	Binary []byte
	String string
)

func (m Binary) Bytes() []byte {
	return m
}

func (m Binary) String() string {
	return string(m)
}

func (m Binary) WriteTo(w io.Writer) (int64, error) {
	err := binary.Write(w, binary.BigEndian, BinaryType) // 1-byte type
	if err != nil {
		return 0, err
	}
	n := int64(1)

	err = binary.Write(w, binary.BigEndian, uint32(len(m))) // 4-byte size
	if err != nil {
		return n, err
	}
	n += 4
	o, err := w.Write(m) //payload
	return n + int64(o), err
}

func (m *Binary) ReadFrom(r io.Reader) (int64, error) {
	var typ uint8
	err := binary.Read(r, binary.BigEndian, &typ)
	if err != nil {
		return 0, err
	}

	var n int64 = 1
	if typ != BinaryType {
		return n, errors.New("invalid Binary")
	}
	var size uint32
	err = binary.Read(r, binary.BigEndian, &size) // 4-byte size
	if err != nil {
		return n, err
	}
	n += 4
	if size > MaxPayloadSize {
		return n, ErrMaxPayloadSize
	}
	*m = make([]byte, size)
	o, err := r.Read(*m) // payload
	return n + int64(o), err

}

func (m String) Bytes() []byte {
	return []byte(m)
}
func (m String) String() string {
	return string(m)

}

func (m String) WriteTo(w io.Writer) (n int64, err error) {
	err = binary.Write(w, binary.BigEndian, StringType) // 1-byte type
	if err != nil {
		return 0, err
	}
	n = 1
	err = binary.Write(w, binary.BigEndian, uint32(len(m))) // 4-byte size
	if err != nil {
		return n, err
	}
	n += 4
	o, err := w.Write([]byte(m)) // payload
	return n + int64(o), err

}

func (m *String) ReadFrom(r io.Reader) (n int64, err error) {
	var typ uint8
	err = binary.Read(r, binary.BigEndian, &typ) // 1-byte type
	if err != nil {
		return 0, err
	}
	n = 1
	if typ != StringType {
		return n, errors.New("invalid String")
	}
	var size uint32
	err = binary.Read(r, binary.BigEndian, &size) // 4-byte size
	if err != nil {
		return n, err
	}
	n += 4
	buf := make([]byte, size)
	o, err := r.Read(buf) // payload
	if err != nil {
		return n, err
	}
	*m = String(buf)
	return n + int64(o), nil

}

func decode(r io.Reader) (Payload, error) {
	var typ uint8
	err := binary.Read(r, binary.BigEndian, &typ)
	if err != nil {
		return nil, err
	}

	var payload Payload
	switch typ {
	case BinaryType:
		payload = new(Binary)
	case StringType:
		payload = new(String)
	default:
		return nil, errors.New("unknown type")
	}

	_, err = payload.ReadFrom(
		io.MultiReader(bytes.NewReader([]byte{typ}), r))
	if err != nil {
		return nil, err
	}
	return payload, nil

}

func PayloadListener() {
	listener, err := net.Listen(tcp, addr)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := listener.Accept()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer conn.Close()
	for _, p := range payloads {
		_, err = p.WriteTo(conn)
		if err != nil {
			log.Fatal(err)
			break
		}
	}
}

func ProxyConnection(source, destination string) error {
	connSource, err := net.Dial(tcp, source)
	if err != nil {
		return err
	}
	defer connSource.Close()
	connDestination, err := net.Dial("tcp", destination)
	if err != nil {
		return err
	}
	defer connDestination.Close()

	// connDestination replies to connSource
	go func() { _, _ = io.Copy(connSource, connDestination) }()

	// connSource messages to connDestination
	_, err = io.Copy(connDestination, connSource)
	return err
}

func proxy(from, to io.ReadWriter) error {

	go func() {
		_, _ = io.Copy(to, from)
	}()

	_, err := io.Copy(from, to)
	return err
}

func ProxyListener() {

	listener, err := net.Listen(tcp, addr)
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		go func(c net.Conn) {
			defer c.Close()
			for {
				buf := make([]byte, 1024)
				n, err := c.Read(buf)
				if err != nil {
					if err != io.EOF {
						log.Fatal(err)
					}
					return
				}
				switch msg := string(buf[:n]); msg {
				case "ping":
					_, err = c.Write([]byte("pong"))
				default:
					_, err = c.Write(buf[:n])
				}
				if err != nil {
					if err != io.EOF {
						log.Fatal(err)
					}
					return
				}
			}
		}(conn)
	}
}
