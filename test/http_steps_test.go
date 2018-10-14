package fujilane

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"

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

func assertResponseStatusAndListStep(slice interface{}) func(string, *gherkin.DataTable) error {
	return func(status string, table *gherkin.DataTable) error {
		if err := assertResponseStatus(status); err != nil {
			return err
		}

		if err := json.Unmarshal([]byte(response.Body.String()), slice); err != nil {
			return fmt.Errorf("Unable to unmarshal %s: %s", response.Body.String(), err.Error())
		}

		return assist.CompareToSlice(reflect.ValueOf(slice).Elem().Interface(), table)
	}
}

func postTableStep(path string) func(*gherkin.DataTable) error {
	return func(table *gherkin.DataTable) error {
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

func performGETStep(path string) func() error {
	return func() error {
		return perform("GET", path, nil)
	}
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

var statusNames map[string]int

func init() {
	statusNames = map[string]int{}
	statusNames["CREATED"] = http.StatusCreated
	statusNames["NOT FOUND"] = http.StatusNotFound
	statusNames["NOT MODIFIED"] = http.StatusNotModified
	statusNames["OK"] = http.StatusOK
	statusNames["PRECONDITION REQUIRED"] = http.StatusPreconditionRequired
	statusNames["UNAUTHORIZED"] = http.StatusUnauthorized
	statusNames["UNPROCESSABLE ENTITY"] = http.StatusUnprocessableEntity
}

func assertResponseStatus(status string) error {
	if response == nil {
		return fmt.Errorf("Response is nil, are you sure you made any HTTP request?")
	}

	if code, ok := statusNames[status]; ok {
		if response.Code != code {
			return fmt.Errorf("Expected response to be status %d, got %d", code, response.Code)
		}

		return nil
	}

	return fmt.Errorf("Unhandled status: %s", status)
}

func HTTPContext(s *godog.Suite) {
	s.Step(`^the system should respond with "([^"]*)" and no body$`, assertResponseStatusAndNoBody)
	s.Step(`^the system should respond with "([^"]*)" and the following body:$`, assertResponseStatusAndBody)
	s.Step(`^the system should respond with "([^"]*)" and the following JSON:$`, assertResponseStatusAndJSON)
	s.Step(`^the system should respond with "([^"]*)" and the following errors:$`, assertResponseStatusAndErrors)
	s.Step(`^the system should respond with "([^"]*)"$`, assertResponseStatus)
}
