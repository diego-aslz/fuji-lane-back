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

func assertResponseStatusTextAndNoBody(status string) error {
	if err := assertResponseStatusText(status); err != nil {
		return err
	}

	body := response.Body.String()
	if len(body) > 0 {
		return fmt.Errorf("Expected no body in response, got %s", body)
	}

	return nil
}

func assertResponseStatusTextAndBody(status string, table *gherkin.DataTable) error {
	if err := assertResponseStatusText(status); err != nil {
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

func assertResponseStatusTextAndErrors(status string, table *gherkin.DataTable) error {
	if err := assertResponseStatusText(status); err != nil {
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
	if err := assertResponseStatusText(status); err != nil {
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

func assertResponseStatusTextAndPresignedURL(status string, table *gherkin.DataTable) error {
	if err := assertResponseStatusText(status); err != nil {
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

func assertResponseStatusText(status string) error {
	if response == nil {
		return fmt.Errorf("Response is nil, are you sure you made any HTTP request?")
	}

	switch status {
	case "CREATED":
		return assertResponseStatus(http.StatusCreated)
	case "NOT FOUND":
		return assertResponseStatus(http.StatusNotFound)
	case "OK":
		return assertResponseStatus(http.StatusOK)
	case "PRECONDITION REQUIRED":
		return assertResponseStatus(http.StatusPreconditionRequired)
	case "UNAUTHORIZED":
		return assertResponseStatus(http.StatusUnauthorized)
	case "UNPROCESSABLE ENTITY":
		return assertResponseStatus(http.StatusUnprocessableEntity)
	default:
		return fmt.Errorf("Unhandled status: %s", status)
	}
}

func assertResponseStatus(expected int) error {
	if response.Code != expected {
		return fmt.Errorf("Expected response to be status %d, got %d", expected, response.Code)
	}
	return nil
}

func HTTPContext(s *godog.Suite) {
	s.Step(`^the system should respond with "([^"]*)" and the following body:$`, assertResponseStatusTextAndBody)
	s.Step(`^the system should respond with "([^"]*)" and no body$`, assertResponseStatusTextAndNoBody)
	s.Step(`^the system should respond with "([^"]*)" and the following errors:$`, assertResponseStatusTextAndErrors)
	s.Step(`^the system should respond with "([^"]*)" and the following pre-signed URL:$`, assertResponseStatusTextAndPresignedURL)
	s.Step(`^the system should respond with "([^"]*)" and the following countries:$`, assertResponseStatusAndCountries)
	s.Step(`^the system should respond with "([^"]*)"$`, assertResponseStatusText)
}
