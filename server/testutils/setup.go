package testutils

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
)

func NewHTTPTest(method, url string, body interface{}) (*http.Request, *httptest.ResponseRecorder) {
	var b io.Reader
	var err error

	if result, ok := body.(io.Reader); !ok && body != nil {
		var nb []byte
		nb, err = json.Marshal(body)
		b = bytes.NewReader(nb)
	} else {
		b = result
	}

	if err != nil {
		log.Fatalln(err)
	}

	r, _ := http.NewRequest(method, url, b)
	rr := httptest.NewRecorder()

	return r, rr
}
