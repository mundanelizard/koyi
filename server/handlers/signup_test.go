package handlers

import (
	"encoding/json"
	"github.com/mundanelizard/koyi/server/testutils"
	"log"
	"net/http"
	"testing"
)

func TestEmailSignInHandler(t *testing.T) {
	server := testutils.NewEngine("/v1/", CreateSignUpRoutes)

	t.Run("With wrong credentials", func(t *testing.T) {
		t.Run("[Invalid Email]", func(t *testing.T) {
			user := map[string]string{
				"email":    "mundanelizard.com",
				"password": "SuperC0olC4+",
			}

			r, rr := testutils.NewHTTPTest(http.MethodPost, "/v1/auth/signup/email", user)
			server.ServeHTTP(rr, r)

			testutils.ExpectToBe(t, http.StatusBadRequest, rr.Code)
			testutils.ExpectToBe(t, "{}", rr.Body.String())
		})

		t.Run("[No Email]", func(t *testing.T) {
			user := map[string]string{
				"password": "SuperC0olC4+",
			}

			r, rr := testutils.NewHTTPTest(http.MethodPost, "/v1/auth/signup/email", user)
			server.ServeHTTP(rr, r)

			testutils.ExpectToBe(t, http.StatusBadRequest, rr.Code)
			testutils.ExpectToBe(t, "{}", rr.Body.String())
		})

		t.Run("[Invalid Password]", func(t *testing.T) {
			user := map[string]string{
				"email": "mundanelizard@gmail.com",
			}

			r, rr := testutils.NewHTTPTest(http.MethodPost, "/v1/auth/signup/email", user)
			server.ServeHTTP(rr, r)

			testutils.ExpectToBe(t, http.StatusBadRequest, rr.Code)
			testutils.ExpectToBe(t, "{}", rr.Body.String())
		})
	})

	t.Run("With right credentials", func(t *testing.T) {
		user := map[string]interface{}{
			"email":    "mundanelizard@gmail.com.com",
			"password": "SuperC0olC4+",
			"metadata": map[string]string{
				"firstName": "Mundane",
				"lastName":  "Lizard",
			},
		}

		r, rr := testutils.NewHTTPTest(http.MethodPost, "/v1/auth/signup/email", user)
		server.ServeHTTP(rr, r)

		log.Println(rr.Body.String())

		var body map[string]interface{}
		err := json.Unmarshal(rr.Body.Bytes(), &body)

		if err != nil {
			t.Fatal(err)
		}

		testutils.ExpectToBe(t, http.StatusCreated, rr.Code)
		testutils.ExpectNotToBe(t, body, nil)
		testutils.ExpectNotToBe(t, body["token"], nil)
		testutils.ExpectNotToBe(t, rr.Header().Get("authentication"), nil)
	})

	// Clean Up
	testutils.ClearDatabase()
}

func TestPhoneNumberSignInHandler(t *testing.T) {

	server := testutils.NewEngine("/v1/", CreateSignUpRoutes)

	t.Run("With wrong credentials", func(t *testing.T) {
		t.Run("[Invalid Phone Number]", func(t *testing.T) {
			user := map[string]interface{}{
				"phoneNumber": map[string]string{
					"countryCode":      "4",
					"subscriberNumber": "03399",
				},
				"password": "SuperC0olC4+",
			}

			r, rr := testutils.NewHTTPTest(http.MethodPost, "/v1/auth/signup/phone", user)
			server.ServeHTTP(rr, r)

			testutils.ExpectToBe(t, http.StatusBadRequest, rr.Code)
			testutils.ExpectToBe(t, "{}", rr.Body.String())
		})

		t.Run("[No Phone Number]", func(t *testing.T) {
			user := map[string]string{
				"password": "SuperC0olC4+",
			}

			r, rr := testutils.NewHTTPTest(http.MethodPost, "/v1/auth/signup/phone", user)
			server.ServeHTTP(rr, r)

			testutils.ExpectToBe(t, http.StatusBadRequest, rr.Code)
			testutils.ExpectToBe(t, "{}", rr.Body.String())
		})

		t.Run("[Invalid Password]", func(t *testing.T) {
			user := map[string]interface{}{
				"phoneNumber": map[string]string{
					"countryCode":      "44",
					"subscriberNumber": "033935589",
				},
			}

			r, rr := testutils.NewHTTPTest(http.MethodPost, "/v1/auth/signup/phone", user)
			server.ServeHTTP(rr, r)

			testutils.ExpectToBe(t, http.StatusBadRequest, rr.Code)
			testutils.ExpectToBe(t, "{}", rr.Body.String())
		})
	})

	t.Run("With right credentials", func(t *testing.T) {
		user := map[string]interface{}{
			"countryCode":      "44",
			"subscriberNumber": "033935589",
			"password":         "SuperC0olC4+",
			"metadata": map[string]string{
				"firstName": "Mundane",
				"lastName":  "Lizard",
			},
		}

		r, rr := testutils.NewHTTPTest(http.MethodPost, "/v1/auth/signup/phone", user)
		server.ServeHTTP(rr, r)

		var body map[string]interface{}
		err := json.Unmarshal(rr.Body.Bytes(), &body)

		if err != nil {
			t.Fatal(err)
		}

		testutils.ExpectToBe(t, http.StatusCreated, rr.Code)
		testutils.ExpectNotToBe(t, body, nil)
		testutils.ExpectNotToBe(t, body["token"], nil)
		testutils.ExpectNotToBe(t, rr.Header().Get("authentication"), nil)
	})

	// Clean Up
	//testutils.ClearDatabase()
}
