package main

import (
	"strings"
)

// Clean the values in the Linkk
func (msg *Linkk) Clean() {
	if !strings.HasPrefix(msg.Path, "/") {
		msg.Path = "/" + msg.Path
	}
	if strings.HasSuffix(msg.Path, "/") {
		msg.Path = msg.Path[:len(msg.Path)-1]
	}
}
