package main

import (
	"strings"
)

// Clean the values in the Linkk
func (msg *Linkk) Clean() {
	if strings.HasSuffix(msg.Path, "/") {
		msg.Path = msg.Path[:len(msg.Path)-1]
	}

	if !strings.HasPrefix(msg.Path, "/") {
		msg.Path = "/" + msg.Path
	}

	msg.Path = strings.ToLower(msg.Path)
}
