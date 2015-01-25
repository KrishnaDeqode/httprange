package httprange

import "github.com/justinas/alice"
import "net/http/httptest"
import "io/ioutil"
import "net/http"
import "strconv"
import "testing"
import "errors"
import "fmt"

func resText(res http.ResponseWriter, req *http.Request) {
	body := []byte("Hello, world!")
	size := len(body)

	res.Header().Set("Content-Length", strconv.Itoa(size))
	res.Write(body)
}

var chain = alice.New(New()).Then(http.HandlerFunc(resText))

// tests

func TestGetNoRange(t *testing.T) {
	result, _, _ := getRes(chain, "")
	stringEqual(result, "Hello, world!")
}

func TestGetWithRange01(t *testing.T) {
	result, contentRange, contentLength := getRes(chain, getRangeString(0, 2))
	t.Log(contentRange)
	stringEqual(result, "Hel")
	stringEqual(contentLength, "3")
}

func TestGetWithRange02(t *testing.T) {
	result, contentRange, contentLength := getRes(chain, getRangeString(3, 5))
	t.Log(contentRange)
	stringEqual(result, "lo,")
	stringEqual(contentLength, "3")
}

func TestGetWithRange03(t *testing.T) {
	result, contentRange, contentLength := getRes(chain, getRangeString(11, 13))
	t.Log(result, contentRange, contentLength)
}

func getRes(h http.Handler, rangeString string) (result, contentRange, contentLength string) {
	server := httptest.NewServer(h)
	defer server.Close()

	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		panic(err)
	}
	if rangeString != "" {
		req.Header.Set("Range", rangeString)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	// check `Content-Range` and `Content-Length`
	if rangeString != "" {
		contentRange = res.Header.Get("Content-Range")
		contentLength = res.Header.Get("Content-Length")
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		panic(err)
	}
	result = string(body[:])

	// assert
	stringEqual(res.Header.Get("Accept-Ranges"), "bytes")

	return
}

// utils for test

func stringEqual(a, b string) {
	if a != b {
		panic(errors.New(a + " - not equal - " + b))
	}
}

// Example:
//   "Content-Range": "bytes 100-200/1000"
//   "Content-Range": "bytes 100-200/*"
func parseRangeString(r string) (start, end, total int64) {
	fmt.Sscanf(r, "bytes %d-%d/%d", &start, &end, &total)

	if total != 0 && end > total {
		end = total
	}
	if start >= end {
		start = 0
		end = 0
	}

	return
}

// Example:
//   "Range": "bytes=100-200"
func getRangeString(start, end int64) string {
	return fmt.Sprintf("bytes=%d-%d", start, end)
}
