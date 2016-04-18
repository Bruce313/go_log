package main

import (
	"os"
)

type storeEngine interface {
	save(*oneLog) error
}
type storeEngineOptions struct {
	fileName string
}

func newStoreEngine(opts *storeEngineOptions) (storeEngine, error) {
	return newFileStoreEngine(opts.fileName)
}

type fileStoreEngine struct {
	file *os.File
}

func newFileStoreEngine(p string) (*fileStoreEngine, error) {
	f, err := os.OpenFile(p, os.O_RDWR|os.O_CREATE|os.O_APPEND,
		0644)
	if err != nil {
		return nil, err
	}
	fse := &fileStoreEngine{
		file: f,
	}
	return fse, nil
}

func (self *fileStoreEngine) save(l *oneLog) error {
	//todo cache in mem, sync to file
	//mem store engine impl, maybe
	_, err := self.file.Write(l.Content)
	return err
}
