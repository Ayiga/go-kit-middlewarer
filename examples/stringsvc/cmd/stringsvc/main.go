package main

import (
	"net/http"
	"strings"

	trans "github.com/ayiga/go-kit-middlewarer/examples/stringsvc/transport/http"
)

// StringService represents an object that will implement the StringService
// interface
type StringService struct{}

// Uppercase implements StringService
func (StringService) Uppercase(str string) (string, error) {
	return strings.ToUpper(str), nil
}

// Count implements StringService
func (StringService) Count(str string) int {
	return len(str)
}

func main() {
	var svc StringService

	trans.HTTPServersForEndpoints(svc)
	http.ListenAndServe(":9000", nil)
}
