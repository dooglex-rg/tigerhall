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
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

//common testing inputs
type RoutingTestModel struct {
	//description of the route
	description string
	//url path
	route string
	//is error expected?
	expectedError bool
	//http status code
	expectedCode int
}

var fiber_app *fiber.App

//main init function
func TestMain(m *testing.M) {

	//check if env file is is in relative dir
	env_file := ".env"
	if _, err := os.Stat(env_file); err != nil || !os.IsExist(err) {
		env_file = "../../.env"
	}

	//load env variables from .env file
	err := godotenv.Load(env_file)
	CheckError(err, nil)

	DSN := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_NAME_MOCK"),
	)

	//Creating DB connection
	DB, err = sql.Open("postgres", DSN)
	CheckError(err, nil)
	defer DB.Close()
	CheckError(DB.Ping(), nil)
	fiber_app = fiber.New()
	url_router(fiber_app)

	os.Exit(m.Run())
}

//common test pattern
func BasicTestTemplate(t *testing.T, tests []RoutingTestModel, payload interface{}) {
	for _, test := range tests {

		// jsonize the payload
		req_byte, _ := json.Marshal(payload)

		//create new request object
		req, _ := http.NewRequest(
			"POST",
			test.route,
			bytes.NewBuffer(req_byte),
		)
		req.Header.Set("Content-Type", "application/json")

		res, err := fiber_app.Test(req, -1)

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

		//unmarshalling the response body
		var resp ErrorStatus
		json.Unmarshal(body, &resp)

		// Verify, that the reponse body equals the expected body
		assert.Equalf(t, false, resp.Status.Error, test.description)
	}

}

//create tiger endpoint test
func TestCreateTiger(t *testing.T) {

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
			expectedError: true,
			expectedCode:  400,
		},
	}
	payload := PayloadAddNewTiger{
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
	payload.Image = "testuuid"

	BasicTestTemplate(t, tests, payload)
}

//show tigers endpoint test
func TestShowTigers(t *testing.T) {
	tests := []RoutingTestModel{
		{
			description:   "show tiger",
			route:         "/tiger/show",
			expectedError: false,
			expectedCode:  200,
		},
		{
			description:   "non existing route",
			route:         "/i-dont-exist",
			expectedError: true,
			expectedCode:  404,
		},
	}

	BasicTestTemplate(t, tests, "")
}

//create tiger endpoint testing
func TestCreateSighting(t *testing.T) {
	tests := []RoutingTestModel{
		{
			description:   "create new tiger",
			route:         "/sighting/add",
			expectedError: false,
			expectedCode:  200,
		},
		{
			description:   "non existing route",
			route:         "/i-dont-exist",
			expectedError: true,
			expectedCode:  400,
		},
	}
	payload := PayloadAddSighting{
		SightingInfo{
			LastSeen:  "2007-12-12",
			Latitude:  110.15051,
			Longitude: 75.5644100,
		},
		TigerIdModel{TigerId: 1},
	}

	BasicTestTemplate(t, tests, payload)
}

//show sightings endpoint testing
func TestShowSighting(t *testing.T) {
	tests := []RoutingTestModel{
		{
			description:   "create new tiger",
			route:         "/sighting/show",
			expectedError: false,
			expectedCode:  200,
		},
		{
			description:   "non existing route",
			route:         "/i-dont-exist",
			expectedError: true,
			expectedCode:  400,
		},
	}
	payload := TigerIdModel{TigerId: 1}

	BasicTestTemplate(t, tests, payload)
}
