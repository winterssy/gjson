package gjson

import (
	"bytes"
	"io"
	"strconv"
	"unsafe"

	"github.com/winterssy/bufferpool"
	"github.com/winterssy/gjson/internal/json"
)

// Methods from encoding/json or jsoniter.
var (
	Marshal       = json.Marshal
	MarshalIndent = json.MarshalIndent
	Unmarshal     = json.Unmarshal
)

type (
	// Number is a shortcut for float64, represents JSON numbers.
	Number float64

	// Any is a wrapper around an interface{}, represents any JSON value.
	Any struct {
		data interface{}
	}

	// Array is a shortcut for []*Any, represents JSON arrays.
	Array []*Any

	// Object is a shortcut for map[string]interface{}, represents JSON objects.
	Object map[string]interface{}

	// Encoder is an alias of json.Encoder or jsoniter.Encoder.
	Encoder = json.Encoder

	// Decoder is an alias of json.Decoder or jsoniter.Decoder.
	Decoder = json.Decoder
)

// ToInt casts n to an int.
func (n Number) ToInt() int {
	return int(n)
}

// ToInt8 casts n to an int8.
func (n Number) ToInt8() int8 {
	return int8(n)
}

// ToInt16 casts n to an int16.
func (n Number) ToInt16() int16 {
	return int16(n)
}

// ToInt32 casts n to an int32.
func (n Number) ToInt32() int32 {
	return int32(n)
}

// ToInt64 casts n to an int64.
func (n Number) ToInt64() int64 {
	return int64(n)
}

// ToUint casts n to a uint.
func (n Number) ToUint() uint {
	return uint(n)
}

// ToUint8 casts n to a uint8.
func (n Number) ToUint8() uint8 {
	return uint8(n)
}

// ToUint16 casts n to a uint16.
func (n Number) ToUint16() uint16 {
	return uint16(n)
}

// ToUint32 casts n to a uint32.
func (n Number) ToUint32() uint32 {
	return uint32(n)
}

// ToUint64 casts n to a uint64.
func (n Number) ToUint64() uint64 {
	return uint64(n)
}

// ToFloat32 casts n to a float32.
func (n Number) ToFloat32() float32 {
	return float32(n)
}

// ToFloat64 casts n to a float64.
func (n Number) ToFloat64() float64 {
	return float64(n)
}

// ToString casts n to a string.
func (n Number) ToString() string {
	return strconv.FormatFloat(n.ToFloat64(), 'f', -1, 64)
}

// ToString casts v to a string.
// If v is nil or the data type doesn't match, it will return the zero value.
func (v *Any) ToString() string {
	if v == nil {
		return ""
	}

	vv, ok := v.data.(string)
	if !ok {
		return ""
	}

	return vv
}

// ToBoolean casts v to a bool.
// If v is nil or the data type doesn't match, it will return the zero value.
func (v *Any) ToBoolean() bool {
	if v == nil {
		return false
	}

	vv, ok := v.data.(bool)
	if !ok {
		return false
	}

	return vv
}

// ToNumber casts v to a Number.
// If v is nil or the data type doesn't match, it will return the zero value.
func (v *Any) ToNumber() Number {
	if v == nil {
		return 0
	}

	switch vv := v.data.(type) {
	case float64:
		return Number(vv)
	case json.Number:
		n, _ := vv.Float64()
		return Number(n)
	}

	return 0
}

// ToObject casts v to an Object.
// If v is nil or the data type doesn't match, it will return the zero value.
func (v *Any) ToObject() Object {
	if v == nil {
		return nil
	}

	vv, _ := v.data.(map[string]interface{})
	return vv
}

// ToArray casts v to an Array.
// If v is nil or the data type doesn't match, it will return the zero value.
func (v *Any) ToArray() Array {
	if v == nil {
		return nil
	}

	vv, ok := v.data.([]interface{})
	if !ok {
		return nil
	}

	vs := make(Array, len(vv))
	for i, data := range vv {
		vs[i] = &Any{
			data: data,
		}
	}
	return vs
}

// Interface returns the original data of v or nil if v is nil.
func (v *Any) Interface() interface{} {
	if v == nil {
		return nil
	}
	return v.data
}

// IsNull reports whether the original data of v is null or not.
// Note: if v is nil, it returns false.
func (v *Any) IsNull() bool {
	return v != nil && v.data == nil
}

// MarshalJSON implements json.Marshaler interface.
func (v *Any) MarshalJSON() ([]byte, error) {
	return Marshal(v.data)
}

// UnmarshalJSON implements json.Unmarshaler interface.
func (v *Any) UnmarshalJSON(data []byte) error {
	return Unmarshal(data, &v.data)
}

// Filter returns the elements of arr that match selector.
func (arr Array) Filter(selector func(index int, value *Any) bool) Array {
	arr2 := make(Array, 0, arr.Size())
	for i, v := range arr {
		if selector(i, v) {
			arr2 = append(arr2, v)
		}
	}
	return arr2
}

// Size returns the length of arr.
func (arr Array) Size() int {
	return len(arr)
}

// Index provides an elegant way to get arr's i'th element without panic even index out of range.
func (arr Array) Index(i int) *Any {
	if i < 0 || i >= arr.Size() {
		return nil
	}

	return arr[i]
}

// ToStrings casts arr to a string array.
func (arr Array) ToStrings() []string {
	vs := make([]string, arr.Size())
	for i, v := range arr {
		vs[i] = v.ToString()
	}
	return vs
}

// ToBooleans casts arr to a bool array.
func (arr Array) ToBooleans() []bool {
	vs := make([]bool, arr.Size())
	for i, v := range arr {
		vs[i] = v.ToBoolean()
	}
	return vs
}

// ToNumbers casts arr to a Number array.
func (arr Array) ToNumbers() []Number {
	vs := make([]Number, arr.Size())
	for i, v := range arr {
		vs[i] = v.ToNumber()
	}
	return vs
}

// ToObjects casts arr to an Object array.
func (arr Array) ToObjects() []Object {
	vs := make([]Object, arr.Size())
	for i, v := range arr {
		vs[i] = v.ToObject()
	}
	return vs
}

// GetAny gets any value associated with key recursively.
// If the key not exists, it returns nil.
func (obj Object) GetAny(key string, childNodes ...string) *Any {
	data, ok := obj[key]
	if !ok {
		return nil
	}

	v := &Any{data}
	for _, child := range childNodes {
		v = v.ToObject().GetAny(child)
	}
	return v
}

// GetString gets string value associated with key recursively.
// If the key not exists, it returns the zero value.
func (obj Object) GetString(key string, childNodes ...string) string {
	return obj.GetAny(key, childNodes...).ToString()
}

// GetBoolean gets boolean value associated with key recursively.
// If the key not exists, it returns the zero value.
func (obj Object) GetBoolean(key string, childNodes ...string) bool {
	return obj.GetAny(key, childNodes...).ToBoolean()
}

// GetNumber gets number value associated with key recursively.
// If the key not exists, it returns the zero value.
func (obj Object) GetNumber(key string, childNodes ...string) Number {
	return obj.GetAny(key, childNodes...).ToNumber()
}

// GetObject gets object value associated with key recursively.
// If the key not exists, it returns the zero value.
func (obj Object) GetObject(key string, childNodes ...string) Object {
	return obj.GetAny(key, childNodes...).ToObject()
}

// GetArray gets array value associated with key recursively.
// If the key not exists, it returns the zero value.
func (obj Object) GetArray(key string, childNodes ...string) Array {
	return obj.GetAny(key, childNodes...).ToArray()
}

// String implements fmt.Stringer interface, it returns the pretty JSON encoding of obj.
func (obj Object) String() string {
	s, err := EncodeToString(obj, func(enc *Encoder) {
		enc.SetIndent("", "\t")
		enc.SetEscapeHTML(false)
	})
	if err != nil {
		return "{}"
	}

	return s
}

// MarshalToString is like Marshal, but it returns the string instead of []byte.
func MarshalToString(v interface{}) (string, error) {
	b, err := Marshal(v)
	return b2s(b), err
}

// MarshalIndentToString is like MarshalIndent, but it returns the string instead of []byte.
func MarshalIndentToString(v interface{}, prefix, indent string) (string, error) {
	b, err := MarshalIndent(v, prefix, indent)
	return b2s(b), err
}

// UnmarshalFromString is like Unmarshal, but it reads from string instead of []byte.
func UnmarshalFromString(s string, v interface{}) error {
	return Unmarshal([]byte(s), v)
}

var jsonSuffix = []byte{'\n'}

// Encode is like Marshal, but it uses a encoder instead of a marshaler.
func Encode(v interface{}, opts ...func(enc *Encoder)) ([]byte, error) {
	buf := bufferpool.Get()
	defer buf.Free()

	err := NewEncoder(buf, opts...).Encode(v)
	if err != nil {
		return nil, err
	}

	return bytes.TrimSuffix(buf.Bytes(), jsonSuffix), nil
}

// EncodeToString is like Encode, but it returns the string instead of []byte.
func EncodeToString(v interface{}, opts ...func(enc *Encoder)) (string, error) {
	b, err := Encode(v, opts...)
	return b2s(b), err
}

// Decode is like Unmarshal, but it uses a decoder instead of an unmarshaler.
func Decode(data []byte, v interface{}, opts ...func(dec *Decoder)) error {
	return NewDecoder(bytes.NewReader(data), opts...).Decode(v)
}

// DecodeFromString is like Decode, but it reads from string instead of []byte.
func DecodeFromString(s string, v interface{}, opts ...func(dec *Decoder)) error {
	return Decode([]byte(s), v, opts...)
}

// NewEncoder returns a new encoder that writes to w.
func NewEncoder(w io.Writer, opts ...func(enc *Encoder)) *Encoder {
	encoder := json.NewEncoder(w)
	for _, opt := range opts {
		opt(encoder)
	}
	return encoder
}

// NewDecoder returns a new decoder that reads from r.
func NewDecoder(r io.Reader, opts ...func(dec *Decoder)) *Decoder {
	decoder := json.NewDecoder(r)
	for _, opt := range opts {
		opt(decoder)
	}
	return decoder
}

// Parse parses the JSON-encoded data and stores the result in an Object.
func Parse(data []byte) (obj Object, err error) {
	return obj, Unmarshal(data, &obj)
}

// ParseFromString is like Parse, but it reads from string instead of []byte.
func ParseFromString(s string) (obj Object, err error) {
	return Parse([]byte(s))
}

// ParseFromReader reads the next JSON-encoded value from its input and stores it in an Object.
func ParseFromReader(r io.Reader, opts ...func(dec *Decoder)) (obj Object, err error) {
	return obj, NewDecoder(r, opts...).Decode(&obj)
}

func b2s(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
