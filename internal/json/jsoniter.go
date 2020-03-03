// +build jsoniter

package json

import jsoniter "github.com/json-iterator/go"

// Methods from jsoniter.
var (
	json          = jsoniter.ConfigCompatibleWithStandardLibrary
	Marshal       = json.Marshal
	MarshalIndent = json.MarshalIndent
	Unmarshal     = json.Unmarshal
	NewEncoder    = json.NewEncoder
	NewDecoder    = json.NewDecoder
)

type (
	// Encoder is an alias of jsoniter.Encoder.
	Encoder = jsoniter.Encoder

	// Decoder is an alias of jsoniter.Decoder.
	Decoder = jsoniter.Decoder

	// Number is an alias of jsoniter.Number.
	Number = jsoniter.Number
)
