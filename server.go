package main

import (
	"bufio"
	"github.com/tj/go-debug"
	"log"
	"net"
)

var deMain = debug.Debug("go_log:main")

func main() {
	//parse config
	ip := net.ParseIP("0.0.0.0")
	port := 3334
	//listen
	addr := &net.TCPAddr{
		IP:   ip,
		Port: port,
	}
	l, err := net.ListenTCP("tcp4", addr)
	if err != nil {
		log.Fatalf("listen tcp:%s\n", err)
	}
	defer l.Close()
	listen(l)
}

func listen(l *net.TCPListener) {
	var cs []*net.TCPConn
	for {
		conn, err := l.AcceptTCP()
		if err != nil {
			deMain("go err when accept:%s", err)
			continue
		}
		cs = append(cs, conn)
		go handleConn(conn)
	}
}

const log_DELIM = byte('\n')

func handleConn(c *net.TCPConn) {
	br := bufio.NewReader(c)
	for {
		line, err := br.ReadString(log_DELIM)
		if err != nil {
			deMain("read string from conn:%s", err)
			continue
		}
		deMain("[LOG CONTENT]:%s", line)
	}
}
