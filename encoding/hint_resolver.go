package encoding

import (
	"bytes"
	httptransport "github.com/go-kit/kit/transport/http"
	"io"
	"io/ioutil"
	"net/http"
)

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
		buf := new(bytes.Buffer)
		var i int64 = 0
		for i < r.ContentLength {
			l, err := io.CopyN(buf, r.Body, r.ContentLength)
			if err != nil {
				return request, err
			}
			i += l
		}

		// replace the body
		r.Body = ioutil.NopCloser(buf)

		rune1, _, err := buf.ReadRune()
		if err != nil {
			buf.UnreadRune()
			return request, err
		}

		buf.UnreadRune()

		for mime, hints := range mimeToFirstRunes {
			for _, rune2 := range hints {
				if rune1 == rune2 {
					encoding, err := Get(mime)
					if err != nil {
						// not found... this should be impossible
						// but it's good to check it anyway.
						continue
					}

					if request, err = encoding.DecodeRequest(request)(r); err != nil {
						buf.Reset()
						// encoding failed... let's retry
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
			}
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
		buf := new(bytes.Buffer)
		var i int64 = 0
		for i < r.ContentLength {
			l, err := io.CopyN(buf, r.Body, r.ContentLength)
			if err != nil {
				return response, err
			}
			i += l
		}
		// replace the body
		r.Body = ioutil.NopCloser(buf)

		rune1, _, err := buf.ReadRune()
		defer buf.Reset()
		if err != nil {
			return response, err
		}

		for mime, hints := range mimeToFirstRunes {
			for _, rune2 := range hints {
				if rune1 == rune2 {
					encoding, err := Get(mime)
					if err != nil {
						// not found... this should be impossible
						// but it's good to check it anyway.
						continue
					}

					if response, err = encoding.DecodeResponse(response)(r); err != nil {
						buf.Reset()
						// encoding failed... let's retry
						continue
					}
					// we succeeded
					return response, nil
				}
			}
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
