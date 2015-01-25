package httprange

import "net/http"
import "strconv"

// import "fmt"
// import "io"

type rangeResponseWriter struct {
	http.ResponseWriter
	start  int64
	length int64
	flag   int64
}

func (w *rangeResponseWriter) Write(data []byte) (size int, err error) {
	size = len(data)

	if (w.flag+int64(size) <= w.start) || (w.flag >= w.start+w.length) {
		return
	}

	start := w.start - w.flag
	if start < 0 {
		start = 0
	}
	// add flag
	w.flag += int64(size)
	var end int64
	if w.flag <= w.start+w.length {
		end = int64(size) - 1
	} else {
		end = w.start + w.length - (w.flag - int64(size))
	}

	w.ResponseWriter.Write(data[start:end])
	return
}

func New() func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			res.Header().Set("Accept-Ranges", "bytes")

			rangeString := req.Header.Get("Range")
			if rangeString == "" {
				h.ServeHTTP(res, req)
				return
			}
			// BUG: get total content length - 100 for test
			ranges, err := parseRange(rangeString, 100)
			if err != nil {
				http.Error(res, "Requested Range Not Satisfiable", 416)
				return
			}
			start := ranges[0].start
			length := ranges[0].length

			res.WriteHeader(206)
			res.Header().Set("Content-Range", getRange(start, start+length, 0))
			res.Header().Set("Content-Length", strconv.FormatInt(length, 10))
			h.ServeHTTP(&rangeResponseWriter{
				ResponseWriter: res,
				start:          start,
				length:         length,
				flag:           0,
			}, req)
		})
	}
}
