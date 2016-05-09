package encoding

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"

	"golang.org/x/net/context"

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
	return func(ctx context.Context, r *http.Request) (interface{}, error) {
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
	return func(ctx context.Context, r *http.Response) (interface{}, error) {
		if r.StatusCode < 200 || r.StatusCode > 299 {
			// I'm assuming we have an error at this point, and we should
			// represent it as such.
			ct := parseContentType(r.Header.Get("Content-Type"))
			if ct.contentType == "text/plain" {
				// ok, this is just a simple plain text document, there's not
				// much to do here... so we'll transmit it plainly.
				c, err := ioutil.ReadAll(r.Body)
				if err != nil {
					return nil, err
				}

				// this is unfortunate, but we have no other information to
				// decode with, so at the very least, let's report it as an
				// error
				return errors.New(string(c)), nil
			}

			var we WrapperError
			// we'll let the current Decoder try to decode the error.  In this
			// case we'll give it the custom error type we've created, to wrap
			// the error so we can ensure it decodes properly...
			if err := gen(r.Body).Decode(&we); err != nil {
				return nil, err
			}

			if _, ok := we.Err.(error); we.Err != nil && ok {
				return we.Err, nil
			}

			return we, nil
		}

		if err := gen(r.Body).Decode(response); err != nil {
			return nil, err
		}
		return response, nil
	}
}
