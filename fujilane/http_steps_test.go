package fujilane

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"

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

	actualBody := map[string]string{}
	if err := json.Unmarshal([]byte(response.Body.String()), &actualBody); err != nil {
		return fmt.Errorf("Unable to unmarshal %s: %s", response.Body.String(), err.Error())
	}

	for k, v := range expectedBody {
		if actualBody[k] == v {
			continue
		}
		return fmt.Errorf("Expected %s to be %s, got %s", k, v, actualBody[k])
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

func makePOSTRequest(path string, table *gherkin.DataTable) error {
	response = httptest.NewRecorder()

	body, err := assist.ParseMap(table)
	if err != nil {
		return err
	}

	json, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", path, strings.NewReader(string(json)))

	if err != nil {
		return err
	}

	router.ServeHTTP(response, req)

	return nil
}

func assertResponseStatusText(status string) error {
	if response == nil {
		return fmt.Errorf("Response is nil, are you sure you made any HTTP request?")
	}

	switch status {
	case "OK":
		return assertResponseStatus(http.StatusOK)
	case "UNAUTHORIZED":
		return assertResponseStatus(http.StatusUnauthorized)
	case "CREATED":
		return assertResponseStatus(http.StatusCreated)
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
	s.Step(`^the system should respond with "([^"]*)"$`, assertResponseStatusText)
}
