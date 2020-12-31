package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

var data []byte

const MaxLen = 64

func init() {
	alts := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ01234567890")
	const n = MaxLen
	data = make([]byte, n)
	for i := 0; i < n; i++ {
		data[i] = alts[rand.Intn(len(alts))]
	}
}

func main() {
	var noDelay bool
	var quickAck bool
	flag.BoolVar(&noDelay, "noDelay", false, "Set to true to disable Nagle")
	flag.BoolVar(&quickAck, "quickAck", false, "Set to true to disable delay ack")
	flag.Parse()

	log.Println("Nagle disabled:", noDelay, "QuickAck:", quickAck)
	addr := ":8080"
	if s := os.Getenv("PORT"); s != "" {
		addr = ":" + s
	}
	lis, err := net.Listen("tcp4", addr)
	if err != nil {
		log.Fatalln("faield to listen at", addr, err)
	}

	shutdown := make(chan os.Signal)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())

	go acceptConns(ctx, lis, &wg, quickAck, noDelay)

	<-shutdown
	cancel()
	wg.Wait()
	lis.Close()
	log.Println("Server shutdown")
	os.Exit(0)

}

func acceptConns(ctx context.Context, lis net.Listener, wg *sync.WaitGroup, quickAck, noDelay bool) {
	for {
		conn, err := lis.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				return
			default:
				log.Println("Failed to accept:", err)
			}
			return
		}
		tcpConn, ok := conn.(*net.TCPConn)
		if !ok {
			log.Fatalln("failed to cast as TCPConn")
		}
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

		wg.Add(1)
		go func(conn *net.TCPConn) {
			defer wg.Done()
			defer log.Println("Client connection closed")
			defer conn.Close()
			buf := make([]byte, MaxLen)

			go func() {
				const N = 5
				for i := 0; i < N; i++ {
					conn.Write(data)
				}
			}()
			for {
				select {
				case <-ctx.Done():
					conn.Write([]byte("Server shutdown\n"))
					return
				default:
				}
				// Stop from endless wait
				conn.SetReadDeadline(time.Now().Add(time.Second))
				n, err := conn.Read(buf)
				if err != nil {
					if errors.Is(err, os.ErrDeadlineExceeded) {
						continue
					}
					log.Println("failed to read from connection:", err)
					return
				}
				content := string(buf[:n])
				content = strings.TrimSuffix(content, "\n")
				log.Println("Client said:", content)
				switch content {
				case "close", "exit", "bye":
					log.Println("Received close")
					return
				}
			}
		}(tcpConn)
	}
}
