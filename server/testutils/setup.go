package testutils

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
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

func NewEngine(prefix string, g func(router *gin.RouterGroup)) *gin.Engine {
	engine := gin.Default()
	g(engine.Group(prefix))
	return engine
}

func ExpectToBe(t *testing.T, expected, value interface{}) {
	if expected != value {
		t.Fatalf("Expected '%v' received '%v'", expected, value)
	}
}

func ExpectNotToBe(t *testing.T, expected, value interface{}) {
	if expected == value {
		t.Fatalf("Expected '%v' received '%v'", expected, value)
	}
}
