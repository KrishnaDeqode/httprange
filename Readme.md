
[![Build status][travis-img]][travis-url]
[![License][license-img]][license-url]
[![GoDoc][doc-img]][doc-url]

### httprange

HTTP handler wrapper for `content range` support, test with [alice](https://github.com/justinas/alice)

### Example

```go
package main

import "github.com/pkg4go/httprange"
import "github.com/justinas/alice"
import "net/http"
import "strconv"

func main() {
  var chain = alice.New(httprange.New()).Then(http.HandlerFunc(resText))

  http.ListenAndServe(":3000", chain)
}

func resText(res http.ResponseWriter, req *http.Request) {
  body := []byte("Hello, world!")
  size := len(body)

  res.Header().Set("Content-Length", strconv.Itoa(size))
  res.Write(body)
}
```

```
$ http get localhost:3000

HTTP/1.1 200 OK
Accept-Ranges: bytes
Content-Length: 13
Content-Type: text/plain; charset=utf-8

Hello, world!

$ http get localhost:3000 range:bytes=0-4

HTTP/1.1 206 Partial Content
Accept-Ranges: bytes
Content-Length: 5
Content-Range: bytes 0-4/*
Content-Type: text/plain; charset=utf-8

Hello
```

### License
MIT

[doc-img]: http://img.shields.io/badge/GoDoc-reference-green.svg?style=flat-square
[doc-url]: http://godoc.org/github.com/pkg4go/httprange
[travis-img]: https://img.shields.io/travis/pkg4go/httprange.svg?style=flat-square
[travis-url]: https://travis-ci.org/pkg4go/httprange
[license-img]: http://img.shields.io/badge/license-MIT-green.svg?style=flat-square
[license-url]: http://opensource.org/licenses/MIT
