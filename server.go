package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"

	"github.com/tj/go-debug"
)

var deMain = debug.Debug("go_log:main")

func main() {
	//parse config
	ip := net.ParseIP("0.0.0.0")
	port := 3334
	//chan receive log
	chLog := make(chan *oneLog, 100)
	//listen tcp
	go listenTCP(ip, port, (chan<- *oneLog)(chLog))
	//listen http
	go listenHTTP(ip, port, (chan<- *oneLog)(chLog))
	//listen unix

	//create log handler chain
	//stop here, wait
	chWait := make(chan bool)
	<-chWait
}

func listenHTTP(ip net.IP, port int, ch chan<- *oneLog) {
	lhh := &logHTTPHandler{
		chLog: ch,
	}
	s := http.Server{
		Addr:    fmt.Sprintf("%s:%d", ip.String(), port),
		Handler: lhh,
	}
	s.ListenAndServe()
}

type logHTTPHandler struct {
	chLog chan<- *oneLog
}

func (self *logHTTPHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	content, err := ioutil.ReadAll(req.Body)
	if err != nil {
		deMain("logHTTPHandler get req body:", err)
		return
	}
	l := &oneLog{
		Content: content,
	}
	self.chLog <- l
}

func listenTCP(ip net.IP, port int, ch chan<- *oneLog) {
	addr := &net.TCPAddr{
		IP:   ip,
		Port: port,
	}
	l, err := net.ListenTCP("tcp4", addr)
	if err != nil {
		log.Fatalf("listen tcp:%s\n", err)
	}
	defer l.Close()
	var cs []*net.TCPConn
	for {
		conn, err := l.AcceptTCP()
		if err != nil {
			deMain("go err when accept:%s", err)
			continue
		}
		cs = append(cs, conn)
		go handleTCPConn(conn, ch)
	}
}

const log_DELIM = byte('\n')

func handleTCPConn(c *net.TCPConn, ch chan<- *oneLog) {
	br := bufio.NewReader(c)
	for {
		line, err := br.ReadSlice(log_DELIM)
		if err != nil {
			deMain("read string from conn:%s", err)
			continue
		}
		deMain("[LOG CONTENT]:%s", line)
		l := oneLog{
			Content: line,
		}
		ch <- &l
	}
}
