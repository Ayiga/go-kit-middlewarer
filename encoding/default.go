package encoding

import (
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

const defaultEncoding = "application/json"

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
		if err == nil {
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

// EncodeRequest implements RequestResponseEncoding
func (def) EncodeRequest() httptransport.EncodeRequestFunc {
	return func(r *http.Request, request interface{}) error {
		if em, ok := request.(EmbededMime); ok {
			if mime := em.GetMime(); mime != "" {
				encoding, err := Get(mime)
				if err != nil {
					return err
				}
				r.Header.Set("Content-Type", mime)
				r.Header.Set("Accept", mime)
				return encoding.EncodeRequest()(r, request)
			}
		}
		// pick one?
		r.Header.Set("Content-Type", defaultEncoding)
		r.Header.Set("Accept", defaultEncoding)
		return JSON(0).EncodeRequest()(r, request)
	}
}

// DecodeRequest implements RequestResponseEncoding
func (def) DecodeRequest(request interface{}) httptransport.DecodeRequestFunc {
	return func(r *http.Request) (interface{}, error) {
		ct := parseContentType(r.Header.Get("Content-Type"))
		accept := parseAccept(r.Header.Get("Accept"))
		if ct.contentType == "" {
			// fall back...
			if em, ok := request.(EmbededMime); ok {
				if mime := em.GetMime(); mime != "" {
					encoding, err := Get(mime)
					if err != nil {
						return request, err
					}
					return encoding.DecodeRequest(request)(r)
				}
			}
		} else {
			encoding, err := Get(ct.contentType)
			if err == nil {
				if em, ok := request.(EmbededMime); ok {
					// pass it on
					if len(accept.mime) > 0 {
						mime := accept.highest()
						if _, err := Get(mime); err != nil {
							em.SetMime(ct.contentType)
						} else {
							em.SetMime(mime)
						}
					} else {
						em.SetMime(ct.contentType)
					}
				}

				return encoding.DecodeRequest(request)(r)
			}
			return request, err
		}
		// Not sure what to do here... we'll error out since we have no idea...
		return request, ErrMimeNotFound
	}
}

// EncodeResponse implements RequestResponseEncoding
func (def) EncodeResponse() httptransport.EncodeResponseFunc {
	return func(w http.ResponseWriter, response interface{}) error {
		if em, ok := response.(EmbededMime); ok {
			if mime := em.GetMime(); mime != "" {
				encoding, err := Get(mime)
				if err != nil {
					return err
				}
				w.Header().Set("Content-Type", mime)
				return encoding.EncodeResponse()(w, response)
			}
		}

		// pick one?
		w.Header().Set("Content-Type", defaultEncoding)
		return JSON(0).EncodeResponse()(w, response)
	}
}

// DecodeResponse implements RequestResponseEncoding
func (def) DecodeResponse(response interface{}) httptransport.DecodeResponseFunc {
	return func(r *http.Response) (interface{}, error) {
		ct := r.Header.Get("Content-Type")
		if ct == "" {
			// fall back
			if em, ok := response.(EmbededMime); ok {
				if mime := em.GetMime(); mime != "" {
					encoding, err := Get(mime)
					if err != nil {
						return response, err
					}
					return encoding.DecodeResponse(response)(r)
				}
			}
		} else {
			encoding, err := Get(ct)
			if err == nil {
				return encoding.DecodeResponse(response)(r)
			}
			return response, err
		}
		// Not sure what to do here... we'll error out since we have no idea...
		return response, ErrMimeNotFound
	}
}
