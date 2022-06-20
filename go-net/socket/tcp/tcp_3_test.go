package tcp

import (
	"net"
	"testing"
)

func TestExampleMonitor(t *testing.T) {

	// monitor := &Monitor{
	// 	Logger: log.New(os.Stdout, "monitor: ", log.LstdFlags),
	// }

	
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatal(err)
	}
	_, err = conn.Write([]byte("Test\n"))
	if err != nil {
		t.Fatal(err)
	}
	_ = conn.Close()
}
