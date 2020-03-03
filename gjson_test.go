package gjson

import (
	"math"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	_msg = "hello world"
	_num = 3.1415926535
)

func TestNumber(t *testing.T) {
	const (
		intVal     int     = 127
		int8Val    int8    = 127
		int16Val   int16   = -314
		int32Val   int32   = -314159
		int64Val   int64   = -31415926535
		uintVal    uint    = 127
		uint8Val   uint8   = 127
		uint16Val  uint16  = 314
		uint32Val  uint32  = 314159
		uint64Val  uint64  = 31415926535
		float32Val float32 = 3.14159
		float64Val float64 = 3.1415926535
	)

	assert.Equal(t, intVal, Number(intVal).ToInt())
	assert.Equal(t, int8Val, Number(int8Val).ToInt8())
	assert.Equal(t, int16Val, Number(int16Val).ToInt16())
	assert.Equal(t, int32Val, Number(int32Val).ToInt32())
	assert.Equal(t, int64Val, Number(int64Val).ToInt64())
	assert.Equal(t, uintVal, Number(uintVal).ToUint())
	assert.Equal(t, uint8Val, Number(uint8Val).ToUint8())
	assert.Equal(t, uint16Val, Number(uint16Val).ToUint16())
	assert.Equal(t, uint32Val, Number(uint32Val).ToUint32())
	assert.Equal(t, uint64Val, Number(uint64Val).ToUint64())
	assert.Equal(t, float32Val, Number(float32Val).ToFloat32())
	assert.Equal(t, float64Val, Number(float64Val).ToFloat64())
}

func TestNumber_ToString(t *testing.T) {
	var n Number = 10086
	want := "10086"
	assert.Equal(t, want, n.ToString())
}

func TestAny_ToString(t *testing.T) {
	var v *Any
	assert.Empty(t, v.ToString())

	v = &Any{_num}
	assert.Empty(t, v.ToString())

	v = &Any{_msg}
	assert.Equal(t, _msg, v.ToString())
}

func TestAny_ToBoolean(t *testing.T) {
	var v *Any
	assert.False(t, v.ToBoolean())

	v = &Any{_msg}
	assert.False(t, v.ToBoolean())

	v = &Any{true}
	assert.True(t, v.ToBoolean())
}

func TestAny_ToNumber(t *testing.T) {
	var v *Any
	assert.Equal(t, Number(0), v.ToNumber())

	v = &Any{_msg}
	assert.Equal(t, Number(0), v.ToNumber())

	v = &Any{_num}
	assert.Equal(t, _num, v.ToNumber().ToFloat64())
}

func TestAny_ToObject(t *testing.T) {
	var v *Any
	assert.Empty(t, v.ToObject())

	v = &Any{_msg}
	assert.Nil(t, v.ToObject())

	obj := map[string]interface{}{
		"msg": _msg,
	}
	v = &Any{obj}
	assert.Equal(t, Object(obj), v.ToObject())
}

func TestAny_ToArray(t *testing.T) {
	var v *Any
	assert.Empty(t, v.ToArray())

	v = &Any{_msg}
	assert.Empty(t, v.ToArray())

	arr := []interface{}{_msg, _num}
	v = &Any{arr}
	assert.Len(t, v.ToArray(), len(arr))
}

func TestAny_Interface(t *testing.T) {
	var v *Any
	assert.Nil(t, v.Interface())

	v = &Any{_msg}
	assert.Equal(t, _msg, v.Interface())
}

func TestAny_IsNull(t *testing.T) {
	var v *Any
	assert.False(t, v.IsNull())

	v = &Any{nil}
	assert.True(t, v.IsNull())

	v = &Any{_msg}
	assert.False(t, v.IsNull())
}

type S struct {
	Any *Any `json:"any"`
}

func TestAny_MarshalJSON(t *testing.T) {
	s := S{Any: &Any{map[string]interface{}{
		"code": 200,
		"msg":  "ok",
	}}}
	b, err := Marshal(s)
	if assert.NoError(t, err) {
		want := "{\"any\":{\"code\":200,\"msg\":\"ok\"}}"
		assert.Equal(t, want, string(b))
	}
}

func TestAny_UnmarshalJSON(t *testing.T) {
	var s S
	data := []byte("{\"any\":{\"code\":200,\"msg\":\"ok\"}}")
	err := Unmarshal(data, &s)
	if assert.NoError(t, err) {
		obj := s.Any.ToObject()
		assert.Equal(t, Number(200), obj.GetNumber("code"))
		assert.Equal(t, "ok", obj.GetString("msg"))
	}
}

func TestArray_Filter(t *testing.T) {
	arr := Array{
		&Any{100},
		&Any{200},
		&Any{300},
		&Any{400},
	}
	selector := func(index int, value *Any) bool {
		return index%2 == 0
	}
	assert.Equal(t, 2, arr.Filter(selector).Size())
}

func TestArray_Size(t *testing.T) {
	arr := Array{
		&Any{_msg},
		&Any{_num},
	}
	assert.Equal(t, 2, arr.Size())
}

func TestArray_Index(t *testing.T) {
	arr := Array{
		&Any{_msg},
		&Any{true},
	}
	assert.Equal(t, _msg, arr.Index(0).ToString())
	assert.True(t, arr.Index(1).ToBoolean())
	assert.Nil(t, arr.Index(-1))
	assert.Nil(t, arr.Index(999))
}

func TestArray_ToStrings(t *testing.T) {
	arr := Array{
		&Any{"hello"},
		&Any{"hi"},
	}
	assert.Len(t, arr.ToStrings(), 2)
}

func TestArray_ToBooleans(t *testing.T) {
	arr := Array{
		&Any{true},
		&Any{false},
	}
	assert.Len(t, arr.ToBooleans(), 2)
}

func TestArray_ToNumbers(t *testing.T) {
	arr := Array{
		&Any{3.14},
		&Any{3.14159},
	}
	assert.Len(t, arr.ToNumbers(), 2)
}

func TestArray_ToObjects(t *testing.T) {
	arr := Array{
		&Any{map[string]interface{}{
			"msg": _msg,
		}},
		&Any{map[string]interface{}{
			"num": _num,
		}},
	}
	assert.Len(t, arr.ToObjects(), 2)
}

func TestObject_GetAny(t *testing.T) {
	err := map[string]interface{}{
		"code": 403,
		"msg":  "403 Forbidden",
	}
	obj := Object{
		"err": err,
	}
	v := &Any{"403 Forbidden"}
	assert.Equal(t, v, obj.GetAny("err", "msg"))

	assert.Nil(t, obj.GetAny("noKey"))
}

func TestObject_GetString(t *testing.T) {
	obj := Object{
		"msg": "hello",
		"obj": map[string]interface{}{
			"msg": "hi",
		},
	}
	assert.Equal(t, "hello", obj.GetString("msg"))
	assert.Equal(t, "hi", obj.GetString("obj", "msg"))
}

func TestObject_GetBoolean(t *testing.T) {
	obj := Object{
		"on": true,
		"obj": map[string]interface{}{
			"on": true,
		},
	}
	assert.True(t, obj.GetBoolean("on"))
	assert.True(t, obj.GetBoolean("obj", "on"))
}

func TestObject_GetNumber(t *testing.T) {
	obj := Object{
		"num": 3.14,
		"obj": map[string]interface{}{
			"num": 3.14159,
		},
	}
	assert.Equal(t, Number(3.14), obj.GetNumber("num"))
	assert.Equal(t, Number(3.14159), obj.GetNumber("obj", "num"))
}

func TestObject_GetObject(t *testing.T) {
	err := map[string]interface{}{
		"code": 403,
		"msg":  "403 Forbidden",
	}
	obj := Object{
		"err": err,
	}
	assert.Equal(t, Object(err), obj.GetObject("err"))
}

func TestObject_GetArray(t *testing.T) {
	arr := []interface{}{_msg, _num}
	obj := Object{
		"arr": arr,
	}
	assert.Len(t, obj.GetArray("arr"), 2)
}

func TestObject_String(t *testing.T) {
	obj := Object{
		"msg": "hello world",
	}
	want := "{\n\t\"msg\": \"hello world\"\n}"
	assert.Equal(t, want, obj.String())

	obj = Object{
		"num": math.Inf(1),
	}
	want = "{}"
	assert.Equal(t, want, obj.String())
}

func TestMarshalToString(t *testing.T) {
	obj := map[string]interface{}{
		"msg": "hello world",
	}
	s, err := MarshalToString(obj)
	if assert.NoError(t, err) {
		assert.Equal(t, "{\"msg\":\"hello world\"}", s)
	}
}

func TestMarshalIndentToString(t *testing.T) {
	obj := map[string]interface{}{
		"msg": "hello world",
	}
	s, err := MarshalIndentToString(obj, "", "\t")
	if assert.NoError(t, err) {
		assert.Equal(t, "{\n\t\"msg\": \"hello world\"\n}", s)
	}
}

func TestUnmarshalFromString(t *testing.T) {
	var obj Object
	err := UnmarshalFromString(`{"code":200,"msg":"ok"}`, &obj)
	if assert.NoError(t, err) {
		assert.Equal(t, 200, obj.GetNumber("code").ToInt())
		assert.Equal(t, "ok", obj.GetString("msg"))
	}
}

func TestEncodeToString(t *testing.T) {
	obj := map[string]interface{}{
		"msg": "hello&world",
	}
	s, err := EncodeToString(obj, func(enc *Encoder) {
		enc.SetEscapeHTML(false)
	})
	if assert.NoError(t, err) {
		assert.Equal(t, "{\"msg\":\"hello&world\"}", s)
	}
}

func TestDecodeFromString(t *testing.T) {
	var obj Object
	err := DecodeFromString(`{"code":200,"msg":"ok"}`, &obj, func(dec *Decoder) {
		dec.UseNumber()
	})
	if assert.NoError(t, err) {
		assert.Equal(t, 200, obj.GetNumber("code").ToInt())
		assert.Equal(t, "ok", obj.GetString("msg"))
	}
}

func TestParseFromString(t *testing.T) {
	obj, err := ParseFromString(`{"code":200,"msg":"ok"}`)
	if assert.NoError(t, err) {
		assert.Equal(t, 200, obj.GetNumber("code").ToInt())
		assert.Equal(t, "ok", obj.GetString("msg"))
	}
}

func TestParseFromReader(t *testing.T) {
	r := strings.NewReader("{\n\t\"msg\": \"hello world\"\n}")
	obj, err := ParseFromReader(r)
	if assert.NoError(t, err) {
		want := Object{
			"msg": "hello world",
		}
		assert.Equal(t, want, obj)
	}
}
