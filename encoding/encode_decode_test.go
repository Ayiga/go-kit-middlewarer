package encoding_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/ayiga/go-kit-middlewarer/encoding"

	"testing"
)

func TestJSONEncodeDecodeRequest(t *testing.T) {
	req := &request{
		Str:  "foo",
		Num:  1.5,
		Bool: true,
		Null: false,
	}
	req.embedMime = new(embedMime)
	req.SetMime("application/json")

	ri, err := http.NewRequest("GET", "/does/not/matter", nil)
	if err != nil {
		panic(err)
	}

	err = encoding.Default().EncodeRequest()(ri, req)
	if err != nil {
		t.Log("Error Encoding Request: %s\n", err)
		t.Fail()
	}

	buf := new(bytes.Buffer)
	ri.Body = ioutil.NopCloser(io.TeeReader(ri.Body, buf))

	resp := new(request)
	resp.embedMime = new(embedMime)

	_, err = encoding.Default().DecodeRequest(resp)(ri)

	str := "{\"str\":\"foo\",\"num\":1.5,\"bool\":true,\"null\":false}\n" // trailing new-line?
	if s := buf.String(); s != str {
		t.Logf("Encoding Does not match: %s\n", s)
		t.Fail()
	}

	if err != nil {
		t.Logf("Request Decode Failed: %s\n", err)
		t.Fail()
	}

	if req.Str != resp.Str {
		t.Logf("req.Str != resp.Str \"%s\" vs \"%s\"\n", req.Str, resp.Str)
		t.Fail()
	}

	if req.Num != resp.Num {
		t.Logf("req.Num != resp.Num %f vs %f\n", req.Num, resp.Num)
		t.Fail()
	}

	if req.Bool != resp.Bool {
		t.Logf("req.Bool != resp.Bool %t vs %t\n", req.Bool, resp.Bool)
		t.Fail()
	}

	if req.Null != resp.Null {
		t.Logf("req.Null != resp.Null %s vs %s\n", req.Null, resp.Null)
		t.Fail()
	}
}

func TestXMLEncodeDecodeRequest(t *testing.T) {
	req := &request{
		Str:  "foo",
		Num:  1.5,
		Bool: true,
		Null: nil,
	}
	req.embedMime = new(embedMime)
	req.SetMime("application/xml")

	ri, err := http.NewRequest("GET", "/does/not/matter", nil)
	if err != nil {
		panic(err)
	}

	err = encoding.Default().EncodeRequest()(ri, req)
	if err != nil {
		t.Log("Error Encoding Request: %s\n", err)
		t.Fail()
	}

	buf := new(bytes.Buffer)
	ri.Body = ioutil.NopCloser(io.TeeReader(ri.Body, buf))

	resp := new(request)
	resp.embedMime = new(embedMime)

	_, err = encoding.Default().DecodeRequest(resp)(ri)

	str := "<request><str>foo</str><num>1.5</num><bool>true</bool></request>"
	if s := buf.String(); s != str {
		t.Logf("Encoding Does not match: %s\n", s)
		t.Fail()
	}

	if err != nil {
		t.Logf("Request Decode Failed: %s\n", err)
		t.Fail()
	}

	if req.Str != resp.Str {
		t.Logf("req.Str != resp.Str \"%s\" vs \"%s\"\n", req.Str, resp.Str)
		t.Fail()
	}

	if req.Num != resp.Num {
		t.Logf("req.Num != resp.Num %f vs %f\n", req.Num, resp.Num)
		t.Fail()
	}

	if req.Bool != resp.Bool {
		t.Logf("req.Bool != resp.Bool %t vs %t\n", req.Bool, resp.Bool)
		t.Fail()
	}

	if req.Null != resp.Null {
		t.Logf("req.Null != resp.Null %s vs %s\n", req.Null, resp.Null)
		t.Fail()
	}
}

func TestGobEncodeDecodeRequest(t *testing.T) {
	req := &request{
		Str:  "foo",
		Num:  1.5,
		Bool: true,
		Null: nil,
	}
	req.embedMime = new(embedMime)
	req.SetMime("application/xml")

	ri, err := http.NewRequest("GET", "/does/not/matter", nil)
	if err != nil {
		panic(err)
	}

	err = encoding.Default().EncodeRequest()(ri, req)
	if err != nil {
		t.Log("Error Encoding Request: %s\n", err)
		t.Fail()
	}

	buf := new(bytes.Buffer)
	ri.Body = ioutil.NopCloser(io.TeeReader(ri.Body, buf))

	resp := new(request)
	resp.embedMime = new(embedMime)

	_, err = encoding.Default().DecodeRequest(resp)(ri)

	byts := []byte{0x3c, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x3e,
		0x3c, 0x73, 0x74, 0x72, 0x3e, 0x66, 0x6f, 0x6f, 0x3c,
		0x2f, 0x73, 0x74, 0x72, 0x3e, 0x3c, 0x6e, 0x75, 0x6d,
		0x3e, 0x31, 0x2e, 0x35, 0x3c, 0x2f, 0x6e, 0x75, 0x6d,
		0x3e, 0x3c, 0x62, 0x6f, 0x6f, 0x6c, 0x3e, 0x74, 0x72,
		0x75, 0x65, 0x3c, 0x2f, 0x62, 0x6f, 0x6f, 0x6c, 0x3e,
		0x3c, 0x2f, 0x72, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
		0x3e}

	b := buf.Bytes()

	if bl1, bl2 := len(byts), len(b); bl1 != bl2 {
		t.Logf("bl1, bl2 := len(byts), len(b); bl1 != bl2: %d vs %d\n", bl1, bl2)
		t.Fail()
	}

	for i := 0; i < len(b); i++ {
		if b[i] != byts[i] {
			t.Logf("b[%d] != byts[%d]: %d vs %s", i, i, b[i], byts[i])
			t.Fail()
		}
	}

	if err != nil {
		t.Logf("Request Decode Failed: %s\n", err)
		t.Fail()
	}

	if req.Str != resp.Str {
		t.Logf("req.Str != resp.Str \"%s\" vs \"%s\"\n", req.Str, resp.Str)
		t.Fail()
	}

	if req.Num != resp.Num {
		t.Logf("req.Num != resp.Num %f vs %f\n", req.Num, resp.Num)
		t.Fail()
	}

	if req.Bool != resp.Bool {
		t.Logf("req.Bool != resp.Bool %t vs %t\n", req.Bool, resp.Bool)
		t.Fail()
	}

	if req.Null != resp.Null {
		t.Logf("req.Null != resp.Null %s vs %s\n", req.Null, resp.Null)
		t.Fail()
	}
}

func TestXMLEncodeDecodeResponse(t *testing.T) {
	req := &request{
		Str:  "foo",
		Num:  1.5,
		Bool: true,
		Null: nil,
	}
	req.embedMime = new(embedMime)
	req.SetMime("application/xml")

	buf := new(bytes.Buffer)
	ri := createResponseWriter(buf)

	err := encoding.Default().EncodeResponse()(ri, req)
	if err != nil {
		t.Log("Error Encoding Request: %s\n", err)
		t.Fail()
	}

	ro := new(http.Response)
	ro.Body = ioutil.NopCloser(buf)
	ro.Header = make(http.Header)
	ro.Header.Set("Content-Type", "application/xml")

	resp := new(request)
	resp.embedMime = new(embedMime)

	str := "<request><str>foo</str><num>1.5</num><bool>true</bool></request>"
	if s := buf.String(); s != str {
		t.Logf("Encoding Does not match: %s\n", s)
		t.Fail()
	}

	_, err = encoding.Default().DecodeResponse(resp)(ro)

	if err != nil {
		t.Logf("Request Decode Failed: %s\n", err)
		t.Fail()
	}

	if req.Str != resp.Str {
		t.Logf("req.Str != resp.Str \"%s\" vs \"%s\"\n", req.Str, resp.Str)
		t.Fail()
	}

	if req.Num != resp.Num {
		t.Logf("req.Num != resp.Num %f vs %f\n", req.Num, resp.Num)
		t.Fail()
	}

	if req.Bool != resp.Bool {
		t.Logf("req.Bool != resp.Bool %t vs %t\n", req.Bool, resp.Bool)
		t.Fail()
	}

	if req.Null != resp.Null {
		t.Logf("req.Null != resp.Null %s vs %s\n", req.Null, resp.Null)
		t.Fail()
	}
}

func TestGobEncodeDecodeResponse(t *testing.T) {
	req := &request{
		Str:  "foo",
		Num:  1.5,
		Bool: true,
		Null: false,
	}
	req.embedMime = new(embedMime)
	req.SetMime("application/gob")

	buf := new(bytes.Buffer)
	ri := createResponseWriter(buf)

	err := encoding.Default().EncodeResponse()(ri, req)
	if err != nil {
		t.Log("Error Encoding Request: %s\n", err)
		t.Fail()
	}

	ro := new(http.Response)
	ro.Body = ioutil.NopCloser(buf)
	ro.Header = make(http.Header)
	ro.Header.Set("Content-Type", "application/gob")

	resp := new(request)
	resp.embedMime = new(embedMime)

	byts := []byte{0x37, 0xff, 0x81, 0x03, 0x01, 0x01, 0x07, 0x72, 0x65,
		0x71, 0x75, 0x65, 0x73, 0x74, 0x01, 0xff, 0x82, 0x00,
		0x01, 0x04, 0x01, 0x03, 0x53, 0x74, 0x72, 0x01, 0x0c,
		0x00, 0x01, 0x03, 0x4e, 0x75, 0x6d, 0x01, 0x08, 0x00,
		0x01, 0x04, 0x42, 0x6f, 0x6f, 0x6c, 0x01, 0x02, 0x00,
		0x01, 0x04, 0x4e, 0x75, 0x6c, 0x6c, 0x01, 0x10, 0x00,
		0x00, 0x00, 0x18, 0xff, 0x82, 0x01, 0x03, 0x66, 0x6f,
		0x6f, 0x01, 0xfe, 0xf8, 0x3f, 0x01, 0x01, 0x01, 0x04,
		0x62, 0x6f, 0x6f, 0x6c, 0x02, 0x02, 0x00, 0x00, 0x00}

	b := buf.Bytes()

	if bl1, bl2 := len(byts), len(b); bl1 != bl2 {
		t.Logf("bl1, bl2 := len(byts), len(b); bl1 != bl2: %d vs %d\n", bl1, bl2)
		t.Fail()
	}

	for i := 0; i < len(b); i++ {
		if b[i] != byts[i] {
			t.Logf("b[%d] != byts[%d]: %d vs %s", i, i, b[i], byts[i])
			t.Fail()
		}
	}

	_, err = encoding.Default().DecodeResponse(resp)(ro)

	if err != nil {
		t.Logf("Request Decode Failed: %s\n", err)
		t.Fail()
	}

	if req.Str != resp.Str {
		t.Logf("req.Str != resp.Str \"%s\" vs \"%s\"\n", req.Str, resp.Str)
		t.Fail()
	}

	if req.Num != resp.Num {
		t.Logf("req.Num != resp.Num %f vs %f\n", req.Num, resp.Num)
		t.Fail()
	}

	if req.Bool != resp.Bool {
		t.Logf("req.Bool != resp.Bool %t vs %t\n", req.Bool, resp.Bool)
		t.Fail()
	}

	if req.Null != resp.Null {
		t.Logf("req.Null != resp.Null %s vs %s\n", req.Null, resp.Null)
		t.Fail()
	}
}
