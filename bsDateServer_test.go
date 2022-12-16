package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/cucumber/godog"
)

const host = "http://localhost:10000"

// Struct to wrap response and error
type apiResponse struct {
	response *http.Response
	err      error
}

// newAPIResponse returns a apiResponse struct
func newAPIResponse() *apiResponse {
	return &apiResponse{response: nil, err: nil}
}

func (a *apiResponse) aRequestIsSentToTheEndpoint(method, endpoint string) error {
	var reader = strings.NewReader("")
	var request, err = http.NewRequest(method, host+endpoint, reader)
	if err != nil {
		return fmt.Errorf("could not create request %s", err.Error())
	}

	a.response, a.err = http.DefaultClient.Do(request)
	if a.err != nil {
		return fmt.Errorf("could not send request %s", err.Error())
	}
	return nil
}

func (a *apiResponse) theHTTPresponseCodeShouldBe(expectedCode int) error {
	if expectedCode != a.response.StatusCode {
		return fmt.Errorf("status code not as expected! Expected '%d', got '%d'", expectedCode, a.response.StatusCode)
	}
	return nil
}

func (a *apiResponse) theResponseContentShouldBe(expectedContent string) error {
	body, _ := ioutil.ReadAll(a.response.Body)
	if expectedContent != strings.TrimSpace(string(body)) {
		return fmt.Errorf("response content not as expected! Expected '%s', got '%s'", expectedContent, string(body))
	}
	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	apiResp := newAPIResponse()
	ctx.Step(`^a "([^"]*)" request is sent to the endpoint "([^"]*)"$`, apiResp.aRequestIsSentToTheEndpoint)
	ctx.Step(`^the HTTP-response code should be "([^"]*)"$`, apiResp.theHTTPresponseCodeShouldBe)
	ctx.Step(`^the response content should be "([^"]*)"$`, apiResp.theResponseContentShouldBe)
}

// Integrating godog with go tests
// Run suite : go test -test.v -test.run ^TestFeatures$
// Run particular scenario: go test -test.v -test.run ^TestFeatures$/^my_scenario$
func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: func(ctx *godog.ScenarioContext) {
			// Add step definitions here.
			apiResp := newAPIResponse()
			ctx.Step(`^a "([^"]*)" request is sent to the endpoint "([^"]*)"$`, apiResp.aRequestIsSentToTheEndpoint)
			ctx.Step(`^the HTTP-response code should be "([^"]*)"$`, apiResp.theHTTPresponseCodeShouldBe)
			ctx.Step(`^the response content should be "([^"]*)"$`, apiResp.theResponseContentShouldBe)

		},
		Options: &godog.Options{
			Format:   "pretty",             // Formatter
			Paths:    []string{"features"}, // Path containing feature files
			TestingT: t,                    // Testing instance that will run subtests.
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
