package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalln(os.Args)
	}
	p := strings.Split(os.Args[1], ":")
	if len(p) != 2 {
		log.Fatalln("invalid addr:", p)
	}
	ip := net.ParseIP(p[0])
	port, err := strconv.Atoi(p[1])
	if err != nil {
		log.Fatalln("invalid port:", p[1])
	}
	conn, err := net.DialTCP("tcp4", nil, &net.TCPAddr{
		IP:   ip,
		Port: port,
	})
	if err != nil {
		log.Fatalln("Dial failed:", err)
	}
	defer conn.Close()
	reader := bufio.NewReader(os.Stdin)

	buf := make([]byte, 1024)
	for {
		fmt.Print("Enter text: ")
		input, _ := reader.ReadString('\n')
		if input == "exit" {
			return
		}
		if _, err := conn.Write([]byte(input)); err != nil {
			log.Fatalln("failed to write:", err)
		}
		n, err := conn.Read(buf)
		if err != nil {
			log.Fatalln("failed to read:", err)
		}
		fmt.Println("Response:", string(buf[:n]))
	}
}
