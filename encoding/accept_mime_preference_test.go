package encoding_test

import (
	"bytes"
	"context"
	"net/http"

	"github.com/ayiga/go-kit-middlewarer/encoding"

	"testing"
)

func TestAcceptParsing1(t *testing.T) {
	ctx := context.Background()
	r, err := http.NewRequest("GET", "/not/important", bytes.NewBuffer([]byte(string("{}"))))
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Accept", "application/xml;q=0.7,application/json;q=0.8,application/gob;q=0.9")
	// text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8

	t.Logf("Accept Header: %s\n", r.Header.Get("Accept"))

	req := new(request)
	req.embedMime = new(embedMime)

	def := encoding.Default()

	_, err = def.DecodeRequest(req)(ctx, r)
	if err != nil {
		t.Logf("Decoding Failed: %s\n", err)
		t.Fail()
	}

	if mime := req.GetMime(); mime != "application/gob" {
		t.Logf("mime != \"application/gob\": %s\n", mime)
		t.Fail()
	}
}
