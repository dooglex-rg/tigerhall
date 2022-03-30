package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

type RoutingTestModel struct {
	description   string
	route         string
	expectedError bool
	expectedCode  int
}

func TestCreate_tiger(t *testing.T) {
	tests := []RoutingTestModel{
		{
			description:   "create new tiger",
			route:         "/tiger/add",
			expectedError: false,
			expectedCode:  200,
		},
		{
			description:   "non existing route",
			route:         "/i-dont-exist",
			expectedError: false,
			expectedCode:  400,
		},
	}
	r := PayloadAddNewTiger{
		PayloadTigerBio{
			Name: "test_tiger1",
			Dob:  "1999-12-12",
		},
		SightingInfo{
			LastSeen:  "2005-12-12",
			Latitude:  150.051,
			Longitude: 100.100,
		},
	}
	r.Image = "testuuid"

	BasicTestTemplate(t, tests, r)
}

func BasicTestTemplate(t *testing.T, tests []RoutingTestModel, r interface{}) {
	app := fiber.New()
	url_router(app)

	for _, test := range tests {
		req_byte, _ := json.Marshal(r)

		req, _ := http.NewRequest(
			"POST",
			test.route,
			bytes.NewBuffer(req_byte),
		)
		req.Header.Set("Content-Type", "application/json")

		res, err := app.Test(req, -1)

		// verify that no error occured, that is not expected
		assert.Equalf(t, test.expectedError, err != nil, test.description)

		// As expected errors lead to broken responses, the next
		// test case needs to be processed
		if test.expectedError {
			continue
		}

		// Verify if the status code is as expected
		assert.Equalf(t, test.expectedCode, res.StatusCode, test.description)

		// Read the response body
		body, err := ioutil.ReadAll(res.Body)

		// Reading the response body should work everytime, such that
		// the err variable should be nil
		assert.Nilf(t, err, test.description)

		// Verify, that the reponse body equals the expected body
		assert.Equalf(t, 1, len(body), test.description)
	}

}
