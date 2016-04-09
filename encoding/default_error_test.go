package encoding_test

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/ayiga/go-kit-middlewarer/encoding"
)

func TestDefaultErrorEncodingDecode(t *testing.T) {
	testError := errors.New("this is some sort of error")

	buf := bytes.NewBuffer([]byte(testError.Error()))
	ro := new(http.Response)
	ro.StatusCode = 500
	ro.Body = ioutil.NopCloser(buf)
	ro.Header = make(http.Header)
	ro.Header.Set("Content-Type", "text/plain; encoding=utf-8")
	ro.Header.Set("Content-Length", fmt.Sprintf("%d", buf.Len()))

	resp := new(request)
	r, err := encoding.Default().DecodeResponse(resp)(ro)
	if err != nil {
		t.Fatalf("Unable to Decode Response: %s", err)
	}

	err, ok := r.(error)
	if !ok {
		t.Fatal("Unable to cast returned response into an error")
	}

	if got, want := err.Error(), testError.Error(); got != want {
		t.Errorf(".Error():\ngot:\n\t%s\nwant:\n\t%s", got, want)
	}

	if got, want := r == testError, false; got != want {
		t.Errorf(".Error():\ngot:\n\t%t\nwant:\n\t%t", got, want)
	}
}

func TestDefaultErrorEncodingDecodeWithDefinitiveTypeDecoder(t *testing.T) {
	testError := errors.New("this is some sort of error")

	buf := bytes.NewBuffer([]byte(testError.Error()))
	ro := new(http.Response)
	ro.StatusCode = 500
	ro.Body = ioutil.NopCloser(buf)
	ro.Header = make(http.Header)
	ro.Header.Set("Content-Type", "text/plain; encoding=utf-8")
	ro.Header.Set("Content-Length", fmt.Sprintf("%d", buf.Len()))

	resp := new(request)
	r, err := encoding.JSON(0).DecodeResponse(resp)(ro)
	if err != nil {
		t.Fatalf("Unable to Decode Response: %s", err)
	}

	err, ok := r.(error)
	if !ok {
		t.Fatal("Unable to cast returned response into an error")
	}

	if got, want := err.Error(), testError.Error(); got != want {
		t.Errorf(".Error():\ngot:\n\t%s\nwant:\n\t%s", got, want)
	}

	if got, want := r == testError, false; got != want {
		t.Errorf(".Error():\ngot:\n\t%t\nwant:\n\t%t", got, want)
	}
}
