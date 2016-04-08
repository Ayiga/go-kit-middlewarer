package encoding

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"encoding/xml"
	"errors"
	"reflect"
	"strings"
)

func init() {
	gob.Register(WrapperError{})
}

type WrapperError struct {
	Type      string      `json:"type" xml:"type"`
	ErrString string      `json:"errorString" xml:"error-string"`
	Err       interface{} `json:"error" xml:"error"`
}

func (we WrapperError) Error() string {
	if e, ok := we.Err.(error); ok {
		return e.Error()
	}

	return we.ErrString
}

var ErrUnexpectedJSONDelim = errors.New("Unexpected JSON Delim")

// implements encoding/json.Unmarshaler
func (we *WrapperError) UnmarshalJSON(p []byte) error {
	buf := bytes.NewBuffer(p)
	dec := json.NewDecoder(buf)

	typ := reflect.TypeOf(*we)
	getTag := func(name string) string {
		n, _ := typ.FieldByName(name)
		return n.Tag.Get("json")
	}

	for {
		t, err := dec.Token()
		if err != nil {
			return err
		}
		if delim, ok := t.(json.Delim); ok {
			// have a deliminator
			switch delim.String() {
			case "{":
				continue
			case "}":
				return nil
			default:
				// unexpected Delim
				return ErrUnexpectedJSONDelim
			}
		}

		if str, ok := t.(string); ok {
			switch str {
			case getTag("Type"):
				err = dec.Decode(&we.Type)
				if err != nil {
					break
				}

				e, err := GetErrorInstance(we.Type)
				if err == nil {
					we.Err = e
				}
			case getTag("ErrString"):
				err = dec.Decode(&we.ErrString)
			case getTag("Err"):
				err = dec.Decode(&we.Err)
				if we.Err != nil {
					we.Err = reflect.Indirect(reflect.ValueOf(we.Err)).Interface()
				}
			default:
				continue
			}
		}

		if err != nil {
			return err
		}
	}
	return nil
}

var ErrUnexpectedElementType = errors.New("Unexpected XML Element Type")

// implements encoding/xml.Unmarshaler
func (we *WrapperError) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	typ := reflect.TypeOf(*we)

	getTag := func(name string) string {
		n, _ := typ.FieldByName(name)
		return n.Tag.Get("xml")
	}
	for {
		t, err := d.Token()
		if err != nil {
			return err
		}
		if t == start.End() {
			// we've consumed everything there is
			return nil
		}

		startToken, ok := t.(xml.StartElement)
		if t == nil || !ok {
			return ErrUnexpectedElementType
		}

		switch startToken.Name.Local {
		case getTag("Type"):
			err = d.DecodeElement(&we.Type, &startToken)
			if err != nil {
				break
			}
			if e, err := GetErrorInstance(we.Type); err == nil {
				we.Err = e
			}
		case getTag("ErrString"):
			err = d.DecodeElement(&we.ErrString, &startToken)
		case getTag("Err"):
			err = d.DecodeElement(&we.Err, &startToken)
			if we.Err != nil {
				we.Err = reflect.Indirect(reflect.ValueOf(we.Err)).Interface()
			}
		default:
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func WrapError(e error) *WrapperError {
	t := reflect.TypeOf(e)
	if _, err := GetErrorInstance(t.String()); err != nil {
		// don't transmit errors.errorString types
		return &WrapperError{
			Type:      t.String(),
			ErrString: e.Error(),
		}
	}

	// perhaps some further checking to see if it adheres to the encoding
	// requirements.

	return &WrapperError{
		Type:      t.String(),
		ErrString: e.Error(),
		Err:       e,
	}
}

var ErrBlacklisted = errors.New("This Error type isn't able to registered, as it is not encodable / decodable")
var ErrDuplicate = errors.New("You tried to register a duplicate type")
var ErrUnknownError = errors.New("The type specified hasn't be registered")

var registeredErrors = make(map[string]reflect.Type)

// RegisterError will attempt to register the given Error with the encoders /
// decoders and will make it available for Decoding errors for the encoders /
// decoders.
//
// This will not automatically register this error type with encoding/gob
func RegisterError(e error) error {
	t := reflect.TypeOf(e)
	if reflect.TypeOf(ErrBlacklisted) == t {
		return ErrBlacklisted
	}

	// ensure that we have a pointer
	if t.Kind() == reflect.Ptr {
		t = reflect.ValueOf(e).Type()
	}

	if registeredErrors[t.String()] != nil {
		return ErrDuplicate
	}

	// store the type information
	registeredErrors[t.String()] = t
	return nil
}

// GetErrorInstance will use reflection to attempt and instanciate a new error
// of the given type string.  The error returned will be a pointer.
func GetErrorInstance(s string) (interface{}, error) {
	s = strings.TrimPrefix(s, "*")
	if registeredErrors[s] == nil {
		return nil, ErrUnknownError
	}

	return reflect.New(registeredErrors[s]).Interface(), nil
}
