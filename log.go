package main

import (
	"fmt"
)

type oneLog struct {
	Namespace string
	Category  string
	Content   []byte
}

func (self *oneLog) String() string {
	return fmt.Sprintf("namespace: %s, category: %s, content:\n %s", self.Namespace,
		self.Category, self.Content)
}
