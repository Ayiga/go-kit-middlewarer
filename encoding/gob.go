package encoding

import (
	"encoding/gob"
	"io"

	httptransport "github.com/go-kit/kit/transport/http"
)

func init() {
	Register("application/gob", Gob(0), nil)
	Register("application/octet-stream+gob", Gob(0), nil)
}

// GobGenerateDecoder returns a GOB Decoder
func GobGenerateDecoder(r io.Reader) Decoder {
	return gob.NewDecoder(r)
}

// GobGenerateEncoder returns a GOB Encoder
func GobGenerateEncoder(w io.Writer) Encoder {
	return gob.NewEncoder(w)
}

// Gob is a simple Gob encoder / decoder that conforms to RequestResponseEncoding
type Gob int

// EncodeRequest implements RequestResponseEncoding
func (Gob) EncodeRequest() httptransport.EncodeRequestFunc {
	return MakeRequestEncoder(GobGenerateEncoder)
}

// DecodeRequest implements RequestResponseEncoding
func (Gob) DecodeRequest(request interface{}) httptransport.DecodeRequestFunc {
	return MakeRequestDecoder(request, GobGenerateDecoder)
}

// EncodeResponse implements RequestResponseEncoding
func (Gob) EncodeResponse() httptransport.EncodeResponseFunc {
	return MakeResponseEncoder(GobGenerateEncoder)
}

// DecodeResponse implements RequestResponseEncoding
func (Gob) DecodeResponse(response interface{}) httptransport.DecodeResponseFunc {
	return MakeResponseDecoder(response, GobGenerateDecoder)
}
