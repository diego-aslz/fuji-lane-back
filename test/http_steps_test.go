package fujilane

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"regexp"
	"strings"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/gin-gonic/gin"
	"github.com/go-test/deep"
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

func assertResponseStatusAndEmptyList(status string) error {
	if err := assertResponseStatus(status); err != nil {
		return err
	}

	body := response.Body.String()
	if body != "[]" {
		return fmt.Errorf("Expected no body to be an empty list, got %s", body)
	}

	return nil
}

func assertResponseStatusAndJSONTable(status string, table *gherkin.DataTable) error {
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

	var expectedBody interface{}
	if err := json.Unmarshal([]byte(rawJSON.Content), &expectedBody); err != nil {
		return err
	}

	var actualBody interface{}
	if err := json.Unmarshal([]byte(response.Body.String()), &actualBody); err != nil {
		return fmt.Errorf("Unable to unmarshal %s: %s", response.Body.String(), err.Error())
	}

	if diff := deep.Equal(expectedBody, actualBody); diff != nil {
		return errors.New(strings.Join(diff, "\n"))
	}

	return nil
}

func assertResponseStatusAndXML(status string, xmlDoc *gherkin.DocString) error {
	body := string(regexp.MustCompile("\\n\\s*").ReplaceAll([]byte(xmlDoc.Content), []byte("")))
	return assertResponseStatusAndString(status, body)
}

func assertResponseStatusAndText(status string, textDoc *gherkin.DocString) error {
	return assertResponseStatusAndString(status, textDoc.Content)
}

func assertResponseStatusAndString(status, expectedBody string) error {
	if err := assertResponseStatus(status); err != nil {
		return err
	}

	body := strings.TrimSpace(response.Body.String())
	expectedBody = strings.TrimSpace(expectedBody)
	if expectedBody != body {
		return fmt.Errorf("Expected body:\n%s\nBut got:\n%s", expectedBody, body)
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

func performPOST(path string, body interface{}) (err error) {
	var bodyIO io.Reader

	if body != nil {
		if bodyIO, err = bodyFromObject(body); err != nil {
			return
		}
	}

	return perform("POST", path, bodyIO)
}

func bodyFromObject(body interface{}) (io.Reader, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	return strings.NewReader(string(jsonBody)), nil
}

func slicePayload(src interface{}, fields []string) gin.H {
	payload := gin.H{}

	v := reflect.ValueOf(src)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	for _, f := range fields {
		payload[f] = v.FieldByName(f).Interface()
	}

	return payload
}

func performGETWithParamsStep(path string) func(*gherkin.DataTable) error {
	return func(table *gherkin.DataTable) error {
		queries, err := assist.ParseMap(table)
		if err != nil {
			return err
		}
		finalPath := path
		for key, value := range queries {
			if finalPath == path {
				finalPath += "?"
			} else {
				finalPath += "&"
			}
			finalPath += key + "=" + value
		}

		return performGETStep(finalPath)()
	}
}

func performGETStepWithPage(path string) func(string) error {
	return func(page string) error {
		return perform("GET", path+"?page="+page, nil)
	}
}

func performGETStep(path string) func() error {
	return performStep("GET", path)
}

func performStep(method, path string) func() error {
	return func() error {
		return perform(method, path, nil)
	}
}

func performGETWithQueryStep(path string) func(*gherkin.DataTable) error {
	return func(table *gherkin.DataTable) error {
		params, err := assist.ParseMap(table)
		if err != nil {
			return err
		}

		p := path
		sep := "?"
		for key, value := range params {
			p += sep + key + "=" + value
			sep = "&"
		}

		return perform("GET", p, nil)
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

var statusNames = map[string]int{
	"CREATED":               http.StatusCreated,
	"NOT FOUND":             http.StatusNotFound,
	"NOT MODIFIED":          http.StatusNotModified,
	"OK":                    http.StatusOK,
	"PRECONDITION REQUIRED": http.StatusPreconditionRequired,
	"UNAUTHORIZED":          http.StatusUnauthorized,
	"UNPROCESSABLE ENTITY":  http.StatusUnprocessableEntity,
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

func assertResponseHeaders(table *gherkin.DataTable) error {
	if response == nil {
		return fmt.Errorf("Response is nil, are you sure you made any HTTP request?")
	}

	rowsMap, err := assist.ParseMap(table)
	if err != nil {
		return err
	}

	expected := map[string][]string{}
	actual := map[string][]string{}
	for key, value := range rowsMap {
		expected[key] = strings.Split(value, "|")
		actual[key] = response.HeaderMap[key]
	}

	if diff := deep.Equal(expected, actual); diff != nil {
		return errors.New(strings.Join(diff, "\n"))
	}

	return nil
}

func HTTPContext(s *godog.Suite) {
	s.Step(`^I should receive an? "([^"]*)" response with no body$`, assertResponseStatusAndNoBody)
	s.Step(`^I should receive an? "([^"]*)" response with an empty list$`, assertResponseStatusAndEmptyList)
	s.Step(`^I should receive an? "([^"]*)" response with the following body:$`, assertResponseStatusAndJSONTable)
	s.Step(`^I should receive an? "([^"]*)" response with the following JSON:$`, assertResponseStatusAndJSON)
	s.Step(`^I should receive an? "([^"]*)" response with the following XML:$`, assertResponseStatusAndXML)
	s.Step(`^I should receive an? "([^"]*)" response with the following errors:$`, assertResponseStatusAndErrors)
	s.Step(`^I should receive an? "([^"]*)" response$`, assertResponseStatus)
	s.Step(`^I should receive the following headers:$`, assertResponseHeaders)
	s.Step(`^I should receive an "([^"]*)" response with the following text:$`, assertResponseStatusAndText)
}
