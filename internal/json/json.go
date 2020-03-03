// +build !jsoniter

package json

import "encoding/json"

// Methods from encoding/json.
var (
	Marshal       = json.Marshal
	MarshalIndent = json.MarshalIndent
	Unmarshal     = json.Unmarshal
	NewEncoder    = json.NewEncoder
	NewDecoder    = json.NewDecoder
)

type (
	// Encoder is an alias of json.Encoder.
	Encoder = json.Encoder

	// Decoder is an alias of json.Decoder.
	Decoder = json.Decoder

	// Number is an alias of json.Number.
	Number = json.Number
)
