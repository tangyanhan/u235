package main

import (
	"flag"
	"log"
	"net"
	"os"
)

func main() {
	var mode string
	flag.StringVar(&mode, "mode", "", "Mode to serve rst, emtpy for nothing")
	flag.Parse()

	log.Println("Running at mode:", mode)
	addr := ":8080"
	if s := os.Getenv("PORT"); s != "" {
		addr = ":" + s
	}
	lis, err := net.Listen("tcp4", addr)
	if err != nil {
		log.Fatalln("faield to listen at", addr, err)
	}

	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Fatalln("failed to accept conn:", err)
		}
		tcpConn, ok := conn.(*net.TCPConn)
		if !ok {
			log.Fatalln("failed to cast as TCPConn")
		}
		switch mode {
		case "reset":
			tcpConn.SetLinger(0)
			conn.Close()
		case "close":
			os.Exit(0)
		}
	}
}
