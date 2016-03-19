package encoding

import (
	"bytes"
	httptransport "github.com/go-kit/kit/transport/http"
	"io"
	"io/ioutil"
	"net/http"
)

func copyRequestToBuf(r *http.Request) ([]byte, error) {
	buf := new(bytes.Buffer)
	var i int64 = 0
	for i < r.ContentLength {
		l, err := io.CopyN(buf, r.Body, r.ContentLength)
		if err != nil {
			return nil, err
		}
		i += l
	}

	// replace the body, though thisi s probably not necessary
	return buf.Bytes(), nil
}

func copyResponseToBuf(r *http.Response) ([]byte, error) {
	buf := new(bytes.Buffer)
	var i int64 = 0
	for i < r.ContentLength {
		l, err := io.CopyN(buf, r.Body, r.ContentLength)
		if err != nil {
			return nil, err
		}
		i += l
	}

	// replace the body, though thisi s probably not necessary
	return buf.Bytes(), nil
}

type hintResolver int

// EncodeRequest does not implement RequestResponseEncoding
func (hintResolver) EncodeRequest() httptransport.EncodeRequestFunc {
	return func(r *http.Request, request interface{}) error {
		return ErrNotImplemented
	}
}

// DecodeRequest implements RequestResponseEncoding
func (hintResolver) DecodeRequest(request interface{}) httptransport.DecodeRequestFunc {
	return func(r *http.Request) (interface{}, error) {
		byts, err := copyRequestToBuf(r)
		if err != nil {
			return request, err
		}

		rune1 := []rune(string(byts))[0]

		var mimesToSkip = map[string]bool{}

		for mime, hints := range mimeToFirstRunes {
			for _, rune2 := range hints {
				mimesToSkip[mime] = true
				if rune1 == rune2 {
					encoding, err := Get(mime)
					if err != nil {
						// not found... this should be impossible
						// but it's good to check it anyway.
						break
					}

					r.Body = ioutil.NopCloser(bytes.NewBuffer(byts))

					if _, err = encoding.DecodeRequest(request)(r); err != nil {
						// encoding failed... let's retry
						break
					}

					if em, ok := request.(EmbededMime); ok {
						// let's embed the mime type
						accept := parseAccept(r.Header.Get("Accept"))
						ct := parseContentType(mime)
						transferMimeDetails(em, ct, accept)
					}

					// we succeeded
					return request, nil
				}
			}
		}

		// well... I guess we'll just try all of them, on at a time...
		for mime, encoding := range mimeToEncodings {
			if mimesToSkip[mime] {
				continue
			}

			r.Body = ioutil.NopCloser(bytes.NewBuffer(byts))

			if _, err = encoding.DecodeRequest(request)(r); err != nil {
				continue
			}

			if em, ok := request.(EmbededMime); ok {
				// let's embed the mime type
				accept := parseAccept(r.Header.Get("Accept"))
				ct := parseContentType(mime)
				transferMimeDetails(em, ct, accept)
			}

			// we succeeded
			return request, nil
		}

		return request, ErrUnableToDetermineMime
	}
}

// EncodeResponse does not implement RequestResponseEncoding
func (hintResolver) EncodeResponse() httptransport.EncodeResponseFunc {
	return func(w http.ResponseWriter, response interface{}) error {
		return ErrNotImplemented
	}
}

// DecodeResponse implements RequestResponseEncoding
func (hintResolver) DecodeResponse(response interface{}) httptransport.DecodeResponseFunc {
	return func(r *http.Response) (interface{}, error) {
		byts, err := copyResponseToBuf(r)
		if err != nil {
			return response, err
		}

		rune1 := []rune(string(byts))[0]

		var mimesToSkip = map[string]bool{}

		for mime, hints := range mimeToFirstRunes {
			for _, rune2 := range hints {
				mimesToSkip[mime] = true
				if rune1 == rune2 {
					encoding, err := Get(mime)
					if err != nil {
						// not found... this should be impossible
						// but it's good to check it anyway.
						continue
					}

					r.Body = ioutil.NopCloser(bytes.NewBuffer(byts))

					if _, err = encoding.DecodeResponse(response)(r); err != nil {
						// encoding failed... let's retry
						continue
					}
					// we succeeded
					return response, nil
				}
			}
		}

		// well... I guess we'll just try all of them, on at a time...
		for mime, encoding := range mimeToEncodings {
			if mimesToSkip[mime] {
				continue
			}

			r.Body = ioutil.NopCloser(bytes.NewBuffer(byts))

			if _, err = encoding.DecodeResponse(response)(r); err != nil {
				// error decoding, it's likely not this mime type.
				continue
			}

			if em, ok := response.(EmbededMime); ok {
				// let's embed the mime type
				accept := parseAccept(r.Header.Get("Accept"))
				ct := parseContentType(mime)
				transferMimeDetails(em, ct, accept)
			}

			// we succeeded
			return response, nil
		}

		return response, ErrUnableToDetermineMime
	}
}

// encoder does not implement RequestResponseEncoding
func (hintResolver) encoder(w io.Writer) Encoder {
	return nil
}

// decoder does not implement RequestResponseEncoding
func (hintResolver) decoder(r io.Reader) Decoder {
	return nil
}
