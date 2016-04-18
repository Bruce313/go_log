package main

type logHandler struct {
	chLog chan<- *oneLog
	nsMap map[string]*namespace
}

func newLogHandler(ch chan<- *oneLog) *logHandler {
	return &logHandler{
		chLog: ch,
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
