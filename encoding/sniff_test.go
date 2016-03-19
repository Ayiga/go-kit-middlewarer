package encoding_test

import (
	"bytes"
	"net/http"

	"github.com/ayiga/go-kit-middlewarer/encoding"

	"testing"
)

type request struct {
	Str  string      `json:"str" xml:"str"`
	Num  float64     `json:"num" xml:"num"`
	Bool bool        `json:"bool" xml:"bool"`
	Null interface{} `json:"null" xml:"null"`
}

func TestJSONRequestSniff1(t *testing.T) {
	var e request

	str := "{\"str\":\"bar\",\"num\": 10,\"bool\":true,\"null\":null}"
	t.Logf("Data: %s\n", str)

	buf := bytes.NewBuffer([]byte(str))
	request, err := http.NewRequest("GET", "/test", buf)
	if err != nil {
		panic(err)
	}

	def := encoding.Default()

	e1, err := def.DecodeRequest(&e)(request)
	if err != nil {
		t.Logf("Decode Request Failed: %s\n", err)
		t.Fail()
	}

	if e1 != &e {
		t.Logf("Returned result is the same value: %#v\n", e1)
		t.Fail()
	}

	if e.Str != "bar" {
		t.Logf("e.Str != \"bar\": \"%s\"\n", e.Str)
		t.Fail()
	}

	if e.Num != 10.0 {
		t.Logf("e.Num != 10.0: %f\n", e.Num)
		t.Fail()
	}

	if !e.Bool {
		t.Logf("!e.Bool: %f\n", e.Bool)
		t.Fail()
	}

	if e.Null != nil {
		t.Logf("e.Null != nil: %f\n", e.Null)
		t.Fail()
	}
}

func TestXMLRequestSniff1(t *testing.T) {
	var e request

	str := "<request><str>bar</str><num>10.0</num><bool>true</bool><null>null</null></request>"
	t.Logf("Data: %s\n", str)

	buf := bytes.NewBuffer([]byte(str))
	request, err := http.NewRequest("GET", "/test", buf)
	if err != nil {
		panic(err)
	}

	def := encoding.Default()

	e1, err := def.DecodeRequest(&e)(request)
	if err != nil {
		t.Logf("Decode Request Failed: %s\n", err)
		t.Fail()
	}

	if e1 != &e {
		t.Logf("Returned result is the same value: %#v\n", e1)
		t.Fail()
	}

	if e.Str != "bar" {
		t.Logf("e.Str != \"bar\": \"%s\"\n", e.Str)
		t.Fail()
	}

	if e.Num != 10.0 {
		t.Logf("e.Num != 10.0: %f\n", e.Num)
		t.Fail()
	}

	if !e.Bool {
		t.Logf("!e.Bool: %t\n", e.Bool)
		t.Fail()
	}

	if e.Null != nil {
		t.Logf("e.Null != nil: %s\n", e.Null)
		t.Fail()
	}
}

func TestGobRequestSniff1(t *testing.T) {
	var e request

	b := []byte{0x37, 0xff, 0x81, 0x03, 0x01, 0x01, 0x07, 0x72, 0x65,
		0x71, 0x75, 0x65, 0x73, 0x74, 0x01, 0xff, 0x82, 0x00,
		0x01, 0x04, 0x01, 0x03, 0x53, 0x74, 0x72, 0x01, 0x0c,
		0x00, 0x01, 0x03, 0x4e, 0x75, 0x6d, 0x01, 0x08, 0x00,
		0x01, 0x04, 0x42, 0x6f, 0x6f, 0x6c, 0x01, 0x02, 0x00,
		0x01, 0x04, 0x4e, 0x75, 0x6c, 0x6c, 0x01, 0x10, 0x00,
		0x00, 0x00, 0x0e, 0xff, 0x82, 0x01, 0x03, 0x62, 0x61,
		0x72, 0x01, 0xfe, 0x24, 0x40, 0x01, 0x01, 0x00}

	buf := bytes.NewBuffer(b)
	request, err := http.NewRequest("GET", "/test", buf)
	if err != nil {
		panic(err)
	}

	def := encoding.Default()

	e1, err := def.DecodeRequest(&e)(request)
	if err != nil {
		t.Logf("Decode Request Failed: %s\n", err)
		t.Fail()
	}

	if e1 != &e {
		t.Logf("Returned result is the same value: %#v\n", e1)
		t.Fail()
	}

	if e.Str != "bar" {
		t.Logf("e.Str != \"bar\": \"%s\"\n", e.Str)
		t.Fail()
	}

	if e.Num != 10.0 {
		t.Logf("e.Num != 10.0: %f\n", e.Num)
		t.Fail()
	}

	if !e.Bool {
		t.Logf("!e.Bool: %t\n", e.Bool)
		t.Fail()
	}

	if e.Null != nil {
		t.Logf("e.Null != nil: %s\n", e.Null)
		t.Fail()
	}

}
