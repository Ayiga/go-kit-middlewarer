package encoding

import (
	"errors"
	"io/ioutil"
	"mime"
	"net/http"
	"strconv"
	"strings"

	httptransport "github.com/go-kit/kit/transport/http"
)

// EmbededMime applies to the intermediary transmission representation of method
// arguments and results.  Using the structure we will attempt to communicate
// with the encoders / decoders what mime type to use.
type EmbededMime interface {
	// Retrieves the currently set mime type of this structure. If nothing has
	// been set, it would be best to specify a def one.
	GetMime() string

	// Sets the current mime type to use. If set this mime type should be
	// retrieved by future calls to GetMime on the same variable.
	SetMime(mime string)
}

// Default returns the default RequestResponseEncoding
func Default() RequestResponseEncoding {
	return def(0)
}

// def is the default Encoding handler.  It will attempt to resolve all
// http transmitted encodings based on the information contained within the
// HTTP Headers.
type def int

const DefaultEncoding = "application/json"

// (\w+\/\w+[;q=score],?)+
type acceptContentHeader struct {
	mime  []string
	value []float32
}

func parseAccept(accept string) acceptContentHeader {
	var mimes []string
	var values []float32

	parts := strings.Split(accept, ",")
	for _, p := range parts {
		mediaType, params, err := mime.ParseMediaType(p)
		if err != nil {
			continue
		}

		if params["q"] == "" {
			mimes = append(mimes, mediaType)
			values = append(values, 1.0)
			continue
		}

		v, err := strconv.ParseFloat(params["q"], 32)
		if err != nil {
			continue
		}
		mimes = append(mimes, mediaType)
		values = append(values, float32(v))
	}

	return acceptContentHeader{
		mime:  mimes,
		value: values,
	}
}

func (a acceptContentHeader) highest() string {
	var score float32
	var max = ""
	for i, m := range a.mime {
		if a.value[i] > score {
			score = a.value[i]
			max = m
		}
	}

	return max
}

// specific;options,...
type contentTypeValue struct {
	contentType string
}

func parseContentType(contentType string) contentTypeValue {
	parts := strings.Split(contentType, ",")
	mediaType, _, err := mime.ParseMediaType(parts[0])
	if err != nil {
		return contentTypeValue{}
	}

	return contentTypeValue{
		contentType: mediaType,
	}
}

func getFromEmbededMime(em EmbededMime) (mime string, encoding RequestResponseEncoding, err error) {
	if mime = em.GetMime(); mime != "" {
		encoding, err = Get(mime)
		return
	}

	return "", nil, ErrMimeNotSpecified
}

func encodeRequest(mime string, encoding RequestResponseEncoding, r *http.Request, request interface{}) error {
	r.Header.Set("Content-Type", mime)
	r.Header.Set("Accept", mime)
	return encoding.EncodeRequest()(r, request)
}

func encodeResponse(mime string, encoding RequestResponseEncoding, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", mime)
	return encoding.EncodeResponse()(w, response)
}

func transferMimeDetails(em EmbededMime, ct contentTypeValue, accept acceptContentHeader) {
	if len(accept.mime) > 0 {
		acpt := accept.highest()
		if _, err := Get(acpt); err == nil {
			em.SetMime(acpt)
			return
		}
	}

	em.SetMime(ct.contentType)
}

// EncodeRequest implements RequestResponseEncoding
func (def) EncodeRequest() httptransport.EncodeRequestFunc {
	return func(r *http.Request, request interface{}) error {
		if em, ok := request.(EmbededMime); ok {
			if mime, encoding, err := getFromEmbededMime(em); err == nil {
				return encodeRequest(mime, encoding, r, request)
			}
		}

		// we failed, unfortunately.  However, we are making a request
		// so we can just specify the default encoding
		encoding, err := Get(DefaultEncoding)
		if err != nil {
			// we have really big problems at this point
			return err
		}

		return encodeRequest(DefaultEncoding, encoding, r, request)
	}
}

// DecodeRequest implements RequestResponseEncoding
func (def) DecodeRequest(request interface{}) httptransport.DecodeRequestFunc {
	return func(r *http.Request) (interface{}, error) {
		ct := parseContentType(r.Header.Get("Content-Type"))
		accept := parseAccept(r.Header.Get("Accept"))

		if ct.contentType == "" {
			// let's try to guess the type based on the request
			return hintResolver(0).DecodeRequest(request)(r)
		} else if encoding, err := Get(ct.contentType); err == nil {
			if em, ok := request.(EmbededMime); ok {
				transferMimeDetails(em, ct, accept)
			}

			return encoding.DecodeRequest(request)(r)
		}

		// let's try to guess the type based on the request
		return hintResolver(0).DecodeRequest(request)(r)
	}
}

// EncodeResponse implements RequestResponseEncoding
func (def) EncodeResponse() httptransport.EncodeResponseFunc {
	return func(w http.ResponseWriter, response interface{}) error {
		if em, ok := response.(EmbededMime); ok {
			if mime, encoding, err := getFromEmbededMime(em); err == nil {
				return encodeResponse(mime, encoding, w, response)
			}
		}

		// we failed, but we'll try to use our default, so that we will
		// at least make some forward progress

		encoding, err := Get(DefaultEncoding)
		if err != nil {
			return err

		}

		return encodeResponse(DefaultEncoding, encoding, w, response)
	}
}

// DecodeResponse implements RequestResponseEncoding
func (def) DecodeResponse(response interface{}) httptransport.DecodeResponseFunc {
	return func(r *http.Response) (interface{}, error) {
		ct := parseContentType(r.Header.Get("Content-Type"))
		if ct.contentType == "" {
			// fall back
			// let's try to guess the type based on the response
			return hintResolver(0).DecodeResponse(response)(r)
		} else if encoding, err := Get(ct.contentType); err == nil {
			return encoding.DecodeResponse(response)(r)
		} else if r.StatusCode < 200 || r.StatusCode > 299 {
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
			return hintResolver(0).DecodeResponse(&we)(r)
		}

		// let's try to guess the type based on the response
		return hintResolver(0).DecodeResponse(response)(r)
	}
}
