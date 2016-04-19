package main

import (
	"net/http"

	"github.com/googollee/go-socket.io"
	"github.com/tj/go-debug"
)

//read/write oneLog
type logSuber interface {
	Touch(*oneLog) error
}

type debugerLogSuber struct {
	debuger debug.DebugFunction
}

func newDebugerLogSuber(name string) *debugerLogSuber {
	d := debug.Debug(name)
	return &debugerLogSuber{
		debuger: d,
	}
}

func (self *debugerLogSuber) Touch(l *oneLog) error {
	self.debuger("got the log:%s", l)
	return nil
}

type socketIOLogSuber struct {
	conns []socketio.Socket
	sosvr *socketio.Server
}

const (
	room_NAME          = ""
	event_NAME_NEW_LOG = "new_log"
)

func newSocketIOLogSuber() (*socketIOLogSuber, error) {
	server, err := socketio.NewServer(nil)
	if err != nil {
		return nil, err
	}
	slsb := &socketIOLogSuber{
		sosvr: server,
	}
	server.On("connection", func(so socketio.Socket) {
		deLH("on connection")
		//so.Join(room_NAME)
		slsb.conns = append(slsb.conns, so)
		so.On("disconnection", func() {
			deLH("on disconnect")
			for i, c := range slsb.conns {
				if c.Id() == so.Id() {
					slsb.conns = append(slsb.conns[:i], slsb.conns[i+1:]...)
					break
				}
			}
		})
	})
	server.On("error", func(so socketio.Socket, err error) {
		deLH("error:%s", err)
	})
	return slsb, nil
}

func (self *socketIOLogSuber) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	self.sosvr.ServeHTTP(w, r)
	//TODO: embed
}

func (self *socketIOLogSuber) Touch(l *oneLog) error {
	for _, c := range self.conns {
		deLH("write to socket")
		//c.BroadcastTo(room_NAME, event_NAME_NEW_LOG, l.String())
		c.Emit(event_NAME_NEW_LOG, l.String())
	}
	return nil
}
