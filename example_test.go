package gjson_test

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/winterssy/gjson"
)

func Example_parseToObject() {
	data, err := ioutil.ReadFile("./testdata/music.json")
	if err != nil {
		log.Fatal(err)
	}

	obj, err := gjson.Parse(data)
	if err != nil {
		panic(err)
	}

	fmt.Println(obj.GetNumber("code"))
	fmt.Println(obj.GetString("data", "total"))
	fmt.Println(obj.GetArray("data", "list").Index(0).ToObject().GetString("name"))
	// Output:
	// 200
	// 1917
	// 告白气球
}

func Example_bindToStruct() {
	const dummyData = `
{
  "code": 200,
  "data": {
    "list": [
      {
        "artist": "周杰伦",
        "album": "周杰伦的床边故事",
        "name": "告白气球"
      },
      {
        "artist": "周杰伦",
        "album": "说好不哭 (with 五月天阿信)",
        "name": "说好不哭 (with 五月天阿信)"
      }
    ]
  }
}
`
	var s struct {
		Code int `json:"code"`
		Data struct {
			List gjson.Array `json:"list"`
		} `json:"data"`
	}
	err := gjson.UnmarshalFromString(dummyData, &s)
	if err != nil {
		panic(err)
	}

	fmt.Println(s.Data.List.Index(0).ToObject().GetString("name"))
	// Output:
	// 告白气球
}
