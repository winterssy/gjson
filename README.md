# gjson

**[gjson](https://pkg.go.dev/github.com/winterssy/gjson)** provides a convenient way to read arbitrary JSON in Go.

![Build](https://img.shields.io/github/workflow/status/winterssy/gjson/Test/master?logo=appveyor) [![codecov](https://codecov.io/gh/winterssy/gjson/branch/master/graph/badge.svg)](https://codecov.io/gh/winterssy/gjson) [![Go Report Card](https://goreportcard.com/badge/github.com/winterssy/gjson)](https://goreportcard.com/report/github.com/winterssy/gjson) [![GoDoc](https://img.shields.io/badge/godoc-reference-5875b0)](https://pkg.go.dev/github.com/winterssy/gjson) [![License](https://img.shields.io/github/license/winterssy/gjson.svg)](LICENSE)

## Install

```sh
go get -u github.com/winterssy/gjson
```

## Usage

```go
import "github.com/winterssy/gjson"
```

## Quick Start

- Parse to `gjson.Object`

```go
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
```

- Bind to struct

```go
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
```

## Build with [jsoniter](https://github.com/json-iterator/go)

`gjson` uses `encoding/json` as default json package but you can change to [jsoniter](https://github.com/json-iterator/go) by build from other tags.

```sh
go build -tags=jsoniter .
```

## License

**[MIT](LICENSE)**