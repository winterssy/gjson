package gjson_test

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/winterssy/gjson"
)

func ExampleParse() {
	data, err := ioutil.ReadFile("./testdata/music.json")
	if err != nil {
		log.Fatal(err)
	}

	obj, err := gjson.Parse(data)
	if err != nil {
		panic(err)
	}

	fmt.Println("code:", obj.GetNumber("code"))
	fmt.Println("msg:", obj.GetString("msg"))
	fmt.Println("total of data:", obj.GetString("data", "total"))
	fmt.Println("song name:", obj.GetArray("data", "list").Index(0).ToObject().GetString("name"))
	// Output:
	// code: 200
	// msg: success
	// total of data: 1917
	// song name: 告白气球
}
