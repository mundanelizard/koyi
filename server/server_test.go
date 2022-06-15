package main

import (
	"fmt"
	"github.com/mundanelizard/koyi/server/testutils"
	"net/http"
	"testing"
)

var (
	get = []string{
		"/v1/auth/verify/:intentId/:code",
	}
	post = []string{
		"/v1/auth/signup/email",
		"/v1/auth/signup/phone",
		"/v1/auth/verify",
	}
	update []string
	patch  []string
	del    []string
)

func TestEndpointsHealthCheck(t *testing.T) {
	server := setUpServer()
	methods := map[string][]string{
		"GET":    get,
		"POST":   post,
		"UPDATE": update,
		"PATCH":  patch,
		"DELETE": del,
	}

	for method, urls := range methods {
		for _, url := range urls {
			testName := fmt.Sprintf("[%s]%s", method, url)

			t.Run(testName, func(t *testing.T) {
				r, rr := testutils.NewHTTPTest(method, url, nil)
				server.ServeHTTP(rr, r)

				if rr.Code == http.StatusNotFound {
					t.Fatal("Failed: ", testName)
				}
			})
		}
	}
}
