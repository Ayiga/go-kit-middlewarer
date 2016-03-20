package encoding

import (
	"io"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
)

// Decoder is any type that can populate a given structure via a Decode method.
type Decoder interface {
	// Decode should populate the given interface with the information stored
	// within the Decoder.
	Decode(interface{}) error
}

// GenerateDecoder is a function that takes an io.Reader and returns a Decoder
type GenerateDecoder func(io.Reader) Decoder

// MakeRequestDecoder exists to help bridge the gaps for encoding.  It takes a
// request interface type, and a GenerateDecoder, and ultimately returns a
// function that can decode the given request.
func MakeRequestDecoder(request interface{}, gen GenerateDecoder) httptransport.DecodeRequestFunc {
	return func(r *http.Request) (interface{}, error) {
		if err := gen(r.Body).Decode(request); err != nil {
			return nil, err
		}
		return request, nil
	}
}

// MakeResponseDecoder exists to help bridge the gaps for encoding.  It takes a
// response interface type, and a GenerateDecoder, and ultimately returns a
// function that can decode a given Response.
func MakeResponseDecoder(response interface{}, gen GenerateDecoder) httptransport.DecodeResponseFunc {
	return func(r *http.Response) (interface{}, error) {
		if err := gen(r.Body).Decode(response); err != nil {
			return nil, err
		}
		return response, nil
	}
}
