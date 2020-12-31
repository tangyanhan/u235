package main

import (
	"errors"
	"log"
	"net"
	"os"
	"syscall"
	"testing"
	"time"
)

func TestSend(t *testing.T) {
	const N = 5
	conn, err := net.Dial("tcp4", "localhost:8080")
	if err != nil {
		t.Fatal(err)
	}
	tcpConn, ok := conn.(*net.TCPConn)
	if !ok {
		log.Fatalln("failed to cast as TCPConn")
	}
	quickAck := false
	noDelay := false
	if quickAck {
		rawConn, err := tcpConn.SyscallConn()
		if err != nil {
			log.Fatalln("failed to get syscall conn:", err)
		}
		rawConn.Control(func(fd uintptr) {
			err := syscall.SetsockoptInt(int(fd), syscall.IPPROTO_TCP, syscall.TCP_QUICKACK, 1)
			if err != nil {
				log.Fatalln(err)
			}
		})
	}
	tcpConn.SetNoDelay(noDelay)

	go func() {
		buf := make([]byte, MaxLen)
		for {
			conn.SetReadDeadline(time.Now().Add(time.Second))
			n, err := conn.Read(buf)
			if err != nil {
				if errors.Is(err, os.ErrDeadlineExceeded) {
					continue
				}
				t.Log("failed to read from connection:", err)
				return
			}
			t.Log("Received:", string(buf[:n]))
		}
	}()

	for i := 0; i < N; i++ {
		conn.Write(data)
	}
	conn.Write([]byte("close"))
	time.Sleep(time.Millisecond * 20)
	conn.Close()
}
