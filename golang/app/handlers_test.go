package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
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

func connect_mock_db() {
	DSN := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME_MOCK"),
	)

	var err error
	//Creating DB connection
	DB, err = sql.Open("postgres", DSN)
	CheckError(err, nil)
	CheckError(DB.Ping(), nil)
}

func BasicTestTemplate(t *testing.T, tests []RoutingTestModel, r interface{}) {
	connect_mock_db()
	defer DB.Close()

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

		var resp ErrorStatus
		json.Unmarshal(body, &resp)

		// Verify, that the reponse body equals the expected body
		assert.Equalf(t, false, resp.Status.Error, test.description)
	}

}
