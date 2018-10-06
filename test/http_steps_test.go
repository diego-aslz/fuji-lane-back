package fujilane

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"regexp"
	"strings"

	"github.com/nerde/fuji-lane-back/flentities"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
)

var response *httptest.ResponseRecorder

func assertResponseStatusAndNoBody(status string) error {
	if err := assertResponseStatus(status); err != nil {
		return err
	}

	body := response.Body.String()
	if len(body) > 0 {
		return fmt.Errorf("Expected no body in response, got %s", body)
	}

	return nil
}

func assertResponseStatusAndBody(status string, table *gherkin.DataTable) error {
	if err := assertResponseStatus(status); err != nil {
		return err
	}

	expectedBody, err := assist.ParseMap(table)
	if err != nil {
		return err
	}

	actualBody := map[string]interface{}{}
	if err := json.Unmarshal([]byte(response.Body.String()), &actualBody); err != nil {
		return fmt.Errorf("Unable to unmarshal %s: %s", response.Body.String(), err.Error())
	}

	for k, v := range expectedBody {
		expected := fmt.Sprint(actualBody[k])
		if expected == v {
			continue
		}
		return fmt.Errorf("Expected %s to be %s, got %s", k, v, expected)
	}

	return nil
}

func assertResponseStatusAndJSON(status string, rawJSON *gherkin.DocString) error {
	if err := assertResponseStatus(status); err != nil {
		return err
	}

	expectedBody := map[string]interface{}{}
	err := json.Unmarshal([]byte(rawJSON.Content), &expectedBody)
	if err != nil {
		return err
	}

	actualBody := map[string]interface{}{}
	if err := json.Unmarshal([]byte(response.Body.String()), &actualBody); err != nil {
		return fmt.Errorf("Unable to unmarshal %s: %s", response.Body.String(), err.Error())
	}

	if !reflect.DeepEqual(expectedBody, actualBody) {
		return fmt.Errorf("Response body does not match:\nResponse: %s\nExpected: %s", response.Body.String(),
			rawJSON.Content)
	}

	return nil
}

func assertResponseStatusAndErrors(status string, table *gherkin.DataTable) error {
	if err := assertResponseStatus(status); err != nil {
		return err
	}

	expectedErrors := []string{}
	for _, row := range table.Rows {
		expectedErrors = append(expectedErrors, row.Cells[0].Value)
	}

	body := map[string][]string{}
	if err := json.Unmarshal([]byte(response.Body.String()), &body); err != nil {
		return fmt.Errorf("Unable to unmarshal %s: %s", response.Body.String(), err.Error())
	}

	actualErrors := body["errors"]

	if !reflect.DeepEqual(expectedErrors, actualErrors) {
		return fmt.Errorf(
			"Expected errors [%s], got [%s]",
			strings.Join(expectedErrors, ", "),
			strings.Join(actualErrors, ", "))
	}

	return nil
}

func assertResponseStatusAndCountries(status string, table *gherkin.DataTable) error {
	if err := assertResponseStatus(status); err != nil {
		return err
	}

	countries := []*flentities.Country{}
	if err := json.Unmarshal([]byte(response.Body.String()), &countries); err != nil {
		return fmt.Errorf("Unable to unmarshal %s: %s", response.Body.String(), err.Error())
	}

	return assist.CompareToSlice(countries, table)
}

func performPOSTWithTable(path string, table *gherkin.DataTable) error {
	var body map[string]string

	if len(table.Rows) > 0 {
		var err error
		body, err = assist.ParseMap(table)
		if err != nil {
			return err
		}
	}

	return performPOST(path, body)
}

func performPOST(path string, body interface{}) error {
	var bodyIO io.Reader

	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return err
		}

		bodyIO = strings.NewReader(string(jsonBody))
	}

	return perform("POST", path, bodyIO)
}

func performGET(path string) error {
	return perform("GET", path, nil)
}

func perform(method, path string, body io.Reader) error {
	response = httptest.NewRecorder()

	req, err := http.NewRequest(method, path, body)

	if err != nil {
		return err
	}

	if currentSession != nil {
		req.Header["Authorization"] = []string{fmt.Sprintf("Bearer %s", currentSession.Token)}
	}

	router.ServeHTTP(response, req)

	return nil
}

func assertResponseStatusAndPresignedURL(status string, table *gherkin.DataTable) error {
	if err := assertResponseStatus(status); err != nil {
		return err
	}

	expectedDetails, err := assist.ParseMap(table)
	if err != nil {
		return err
	}

	body := map[string]string{}
	if err := json.Unmarshal([]byte(response.Body.String()), &body); err != nil {
		return fmt.Errorf("Unable to unmarshal %s: %s", response.Body.String(), err.Error())
	}

	url := body["url"]
	reg := regexp.MustCompile("\\/\\/(.*)\\.s3\\.amazonaws\\.com\\/(.*)\\?.*X-Amz-Expires=(\\d+)")
	groups := reg.FindStringSubmatch(url)

	actualDetails := map[string]string{
		"bucket":     groups[1],
		"key":        groups[2],
		"expiration": groups[3],
	}

	for k, v := range expectedDetails {
		expected := actualDetails[k]
		if expected == v {
			continue
		}
		return fmt.Errorf("Expected %s to be %s, got %s", k, v, expected)
	}

	return nil
}

func assertResponseStatus(status string) error {
	if response == nil {
		return fmt.Errorf("Response is nil, are you sure you made any HTTP request?")
	}

	switch status {
	case "CREATED":
		return assertResponseStatusCode(http.StatusCreated)
	case "NOT FOUND":
		return assertResponseStatusCode(http.StatusNotFound)
	case "OK":
		return assertResponseStatusCode(http.StatusOK)
	case "PRECONDITION REQUIRED":
		return assertResponseStatusCode(http.StatusPreconditionRequired)
	case "UNAUTHORIZED":
		return assertResponseStatusCode(http.StatusUnauthorized)
	case "UNPROCESSABLE ENTITY":
		return assertResponseStatusCode(http.StatusUnprocessableEntity)
	default:
		return fmt.Errorf("Unhandled status: %s", status)
	}
}

func assertResponseStatusCode(expected int) error {
	if response.Code != expected {
		return fmt.Errorf("Expected response to be status %d, got %d", expected, response.Code)
	}
	return nil
}

func HTTPContext(s *godog.Suite) {
	s.Step(`^the system should respond with "([^"]*)" and no body$`, assertResponseStatusAndNoBody)
	s.Step(`^the system should respond with "([^"]*)" and the following body:$`, assertResponseStatusAndBody)
	s.Step(`^the system should respond with "([^"]*)" and the following JSON:$`, assertResponseStatusAndJSON)
	s.Step(`^the system should respond with "([^"]*)" and the following errors:$`, assertResponseStatusAndErrors)
	s.Step(`^the system should respond with "([^"]*)" and the following pre-signed URL:$`, assertResponseStatusAndPresignedURL)
	s.Step(`^the system should respond with "([^"]*)" and the following countries:$`, assertResponseStatusAndCountries)
	s.Step(`^the system should respond with "([^"]*)"$`, assertResponseStatus)
}
