package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/tj/go-debug"
)

var deMain = debug.Debug("go_log:main")

func main() {
	//parse config
	ip := net.ParseIP("127.0.0.1")
	portTCP := 3334
	portHTTP := 3344
	portWebUI := 3354
	//chan receive log
	chLog := make(chan *oneLog, 100)
	//listen tcp
	go listenTCP(ip, portTCP, (chan<- *oneLog)(chLog))
	//listen http
	go listenHTTP(ip, portHTTP, (chan<- *oneLog)(chLog))
	//listen unix
	//create log handler chain
	lh := newLogHandler((<-chan *oneLog)(chLog))
	//add debug suber
	desb := newDebugerLogSuber("go_log:debugsuber")
	lh.addSuber(desb)
	//add socket suber
	sioSb, err := newSocketIOLogSuber()
	if err != nil {
		log.Fatal(err)
	}
	lh.addSuber(sioSb)

	//create web ui server
	http.Handle("/", http.FileServer(http.Dir("./webui")))
	http.Handle("/socket.io/", sioSb)
	go http.ListenAndServe(fmt.Sprintf("%s:%d", ip.String(), portWebUI), nil)
	//add file suber to put log on ground
	//log handle, run
	go lh.handle()
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

var (
	rep_OK         = []byte("OK")
	rep_NEED_PARAM = []byte("need namespace and category")
)

func (self *logHTTPHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	content, err := ioutil.ReadAll(req.Body)
	if err != nil {
		deMain("logHTTPHandler get req body:", err)
		return
	}
	reqPath := req.URL.Path
	ps := strings.Split(reqPath, "/")
	deMain("ps:%d, reqpath:%s", len(ps), reqPath)
	if len(ps) < 3 || ps[2] == "" || ps[1] == "" {
		_, _ = w.Write(rep_NEED_PARAM)
		return
	}
	l := &oneLog{
		Namespace: ps[1],
		Category:  ps[2],
		Content:   content,
	}
	self.chLog <- l
	_, _ = w.Write(rep_OK)
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
