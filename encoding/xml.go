package encoding

import (
	"encoding/xml"
	"io"

	httptransport "github.com/go-kit/kit/transport/http"
)

func init() {
	Register("text/xml", XML(0))
	Register("application/xml", XML(0))
}

// XMLGenerateDecoder returns an XML Decoder
func XMLGenerateDecoder(r io.Reader) Decoder {
	return xml.NewDecoder(r)
}

// XMLGenerateEncoder returns an XML Encoder
func XMLGenerateEncoder(w io.Writer) Encoder {
	return xml.NewEncoder(w)
}

// XML is a simple XML encoder / decoder that conforms to RequestResponseEncoding
type XML int

// EncodeRequest implements RequestResponseEncoding
func (XML) EncodeRequest() httptransport.EncodeRequestFunc {
	return MakeRequestEncoder(XMLGenerateEncoder)
}

// DecodeRequest implements RequestResponseEncoding
func (XML) DecodeRequest(request interface{}) httptransport.DecodeRequestFunc {
	return MakeRequestDecoder(request, XMLGenerateDecoder)
}

// EncodeResponse implements RequestResponseEncoding
func (XML) EncodeResponse() httptransport.EncodeResponseFunc {
	return MakeResponseEncoder(XMLGenerateEncoder)
}

// DecodeResponse implements RequestResponseEncoding
func (XML) DecodeResponse(response interface{}) httptransport.DecodeResponseFunc {
	return MakeResponseDecoder(response, XMLGenerateDecoder)
}
