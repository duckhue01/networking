package tcp

import (
	"context"
	"net"
	"sync"
	"syscall"
	"testing"
	"time"
)

func TestDial(t *testing.T) {

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		if err := conn.Close(); err != nil {
			t.Fatal(err)
		}
	}()
}

func DialTimeout(network, address string, timeout time.Duration) (net.Conn, error) {
	d := net.Dialer{
		Control: func(_, addr string, _ syscall.RawConn) error {
			return &net.DNSError{
				Err:         "connection timed out",
				Name:        "",
				Server:      addr,
				IsTimeout:   true,
				IsTemporary: true,
			}
		},
		Timeout: timeout,
	}
	return d.Dial(network, address)
}

func TestDialTimeout(t *testing.T) {
	_, err := DialTimeout("tcp", addr, 1*time.Second)
	nErr, ok := err.(net.Error)

	if !ok {
		t.Error(err)
	}

	if nErr.Timeout() {
		t.Error("error is not a timeout", nErr)
	}
}

func TestDialWithDeadline(t *testing.T) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(1*time.Second))
	defer cancel()
	var d net.Dialer // DialContext is a method on a Dialer
	d.Control = func(_, _ string, _ syscall.RawConn) error {
		// Sleep long enough to reach the context's deadline.
		time.Sleep(10*time.Second + time.Millisecond)
		return nil
	}

	_, err := d.DialContext(ctx, "tcp", "10.0.0.0:80")

	nErr, ok := err.(net.Error)
	if !ok {
		t.Error(err)
	} else {
		if !nErr.Timeout() {
			t.Errorf("error is not a timeout: %v", err)
		}
	}
	if ctx.Err() == context.DeadlineExceeded {
		t.Error(ctx.Err())
	}
}

func TestDialContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	sync := make(chan struct{})
	go func() {
		defer func() { sync <- struct{}{} }()
		var d net.Dialer
		d.Control = func(_, _ string, _ syscall.RawConn) error {
			time.Sleep(10 * time.Second)
			return nil
		}
		conn, err := d.DialContext(ctx, "tcp", "10.0.0.1:80")
		if err != nil {
			t.Log(err)
			return
		}
		conn.Close()
		t.Error("connection did not time out")
	}()
	cancel()
	<-sync
	if ctx.Err() != context.Canceled {
		t.Errorf("expected canceled context; actual: %q", ctx.Err())
	}
}

func TestDialContextCancelFanOut(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

	dial := func(ctx context.Context, addr string, res chan int, id int, wg *sync.WaitGroup) {
		defer wg.Done()
		var d net.Dialer
		c, err := d.DialContext(ctx, "tcp", addr)
		if err != nil {
			return
		}
		c.Close()
		select {
		case <-ctx.Done():
		case res <- id:
		}
	}

	res := make(chan int)
	var wg sync.WaitGroup

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go dial(ctx, addr, res, i+1, &wg)
	}

	response := <-res
	cancel()
	wg.Wait()
	close(res)
	if ctx.Err() != context.Canceled {
		t.Errorf("expected canceled context; actual: %s",
			ctx.Err())
	}
	t.Logf("dialer %d retrieved the resource", response)
}

func TestDeadlineListener(t *testing.T) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatal(err)
	}

	conn.Write([]byte("duckhue01"))

	defer func() {
		err = conn.Close()
		t.Fatal(err)
	}()
}
