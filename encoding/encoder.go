package encoding

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
)

// Encoder is anything that, given an interface, can store an encoding of the
// structure passed into Encode.
type Encoder interface {
	// Encode takes an interface, and should be able to translate it within the
	// given encoding, or it will fail with an error.
	Encode(interface{}) error
}

// GenerateEncoder is a function which takes an io.Writer, and returns an
// Encoder
type GenerateEncoder func(w io.Writer) Encoder

// MakeRequestEncoder takes a GenerateEncoder and returns an
// httptransport.EncodeRequestFunc
func MakeRequestEncoder(gen GenerateEncoder) httptransport.EncodeRequestFunc {
	return func(ctx context.Context, r *http.Request, request interface{}) error {
		var buf bytes.Buffer
		err := gen(&buf).Encode(request)
		r.Body = ioutil.NopCloser(&buf)
		return err
	}
}

// MakeResponseEncoder takes a GenerateEncoder and returns an
// httpstransport.EncodeResponseFunc
func MakeResponseEncoder(gen GenerateEncoder) httptransport.EncodeResponseFunc {
	return func(ctx context.Context, w http.ResponseWriter, response interface{}) error {
		if e, ok := response.(error); ok {
			// we have an error, we'll wrap it, to ensure it's transmission
			// encode-ability

			we := WrapError(e)
			return gen(w).Encode(we)
		}

		return gen(w).Encode(response)
	}
}

// MakeErrorEncoder will take a generic GenerateEncoder function and will
// return an ErrorEncoder
func MakeErrorEncoder(gen RequestResponseEncoding) httptransport.ErrorEncoder {
	return func(ctx context.Context, err error, w http.ResponseWriter) {
		gen.EncodeResponse()(ctx, w, err)
	}
}
