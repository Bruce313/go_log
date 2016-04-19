package main

import (
	"github.com/tj/go-debug"
)

var deLH = debug.Debug("go_log:loghandler")

type logHandler struct {
	subs  []logSuber
	chLog <-chan *oneLog
	nsMap map[string]*namespace
}

func newLogHandler(ch <-chan *oneLog) *logHandler {
	return &logHandler{
		chLog: ch,
	}
}

func (self *logHandler) addSuber(sb logSuber) {
	self.subs = append(self.subs, sb)
}

func (self *logHandler) handle() {
	for {
		lp, ok := <-self.chLog
		if !ok {
			deLH("receive close from chLog, exit")
			break
		}
		//pass log to the chain
		for _, sb := range self.subs {
			err := sb.Touch(lp)
			if err != nil {
				deLH("sb touch err:", err)
			}
		}
	}
}

func (self *logHandler) addNamespace(name string) {
	_, ok := self.nsMap[name]
	if !ok {
		self.nsMap[name] = newNamespace(name)
	}
}

type namespace struct {
	name    string
	cateMap map[string]storeEngine
}

func newNamespace(name string) *namespace {
	return &namespace{
		name: name,
	}
}

func (self *namespace) addCate(name string) error {
	_, ok := self.cateMap[name]
	if !ok {
		se, err := newStoreEngine(&storeEngineOptions{fileName: name})
		if err != nil {
			return err
		}
		self.cateMap[name] = se
	}
	return nil
}
