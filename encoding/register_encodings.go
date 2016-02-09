package encoding

import (
	httptransport "github.com/go-kit/kit/transport/http"
)

// Err are the errors that can be returned from Register or Get
type Err int

const (
	// ErrUnknown represents a non-error
	ErrUnknown Err = iota
	// ErrAlreadyRegistered represents an mime type that has already been
	// registered
	ErrAlreadyRegistered
	// ErrMimeNotFound represents a mime type with no associated encoding
	ErrMimeNotFound
	// ErrNoRegistrationsExist represents that nothing has been registered with
	// this Encoder / Decoder
	ErrNoRegistrationsExist
)

// Error implements the error interface
func (e Err) Error() string {
	switch e {
	case ErrAlreadyRegistered:
		return "That mime type already has already been registered"
	case ErrMimeNotFound:
		return "That mime type does not have an associated Encoder/Decoder"
	case ErrNoRegistrationsExist:
		return "Nothing has been registered, nothing to use for encoding/decoding"
	case ErrUnknown:
		fallthrough
	default:
		return ""
	}
}

// RequestResponseEncoding represents a type that can be used to automatically
// Encode and Decode on HTTP requests used by files generated with
// go-kit-middlewarer
type RequestResponseEncoding interface {
	// EncodeRequest should be able to return an EncodeRequestFunc that can
	// encode the given requests with the encoding type represented by this
	// type.
	EncodeRequest() httptransport.EncodeRequestFunc
	// DecodeRequest should be able to return a DecodeRequestFunc that, when
	// provided with an underlying type, can be used to decode a request with
	// the encoding type represented by this type.
	DecodeRequest(request interface{}) httptransport.DecodeRequestFunc
	// EncodeResponse should be able to return an EncodeResponseFunc that can
	// encode a given response with the encoding type represented by this
	// type.
	EncodeResponse() httptransport.EncodeResponseFunc
	// DecodeResponse should be able to return a DecodeResponseFunc that, when
	// provided with an underlying type, can be used to decode a response with
	// the encoding type represented by this type.
	DecodeResponse(response interface{}) httptransport.DecodeResponseFunc
}

var mimeToEncodings = map[string]RequestResponseEncoding{}

// Register will register the associated encoding with the given mime type
func Register(mime string, encoding RequestResponseEncoding) error {

	if mimeToEncodings[mime] != nil {
		return ErrAlreadyRegistered
	}

	mimeToEncodings[mime] = encoding
	return nil
}

// Get will retrieve the encoding registered with the mime-type
func Get(mime string) (RequestResponseEncoding, error) {
	if len(mimeToEncodings) == 0 {
		return nil, ErrNoRegistrationsExist
	}

	if mimeToEncodings[mime] == nil {
		return nil, ErrNoRegistrationsExist
	}

	return mimeToEncodings[mime], nil
}
