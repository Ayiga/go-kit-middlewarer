package main

import (
	l "log"
	"os"
)

var log *l.Logger

func init() {
	log = l.New(os.Stdout, "go-kit-middlewarer", l.Lshortfile|l.Ltime|l.Ldate|l.Lmicroseconds)
}
