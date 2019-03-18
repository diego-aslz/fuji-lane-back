package fujilane

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
	"github.com/nerde/assistdog"
	"github.com/nerde/assistdog/defaults"
	"github.com/nerde/fuji-lane-back/flconfig"
	"github.com/nerde/fuji-lane-back/flentities"
	"github.com/nerde/fuji-lane-back/optional"
)

func setupAssistdog() {
	os.Setenv("STAGE", "test")

	flconfig.LoadConfiguration()

	assist = assistdog.NewDefault()
	assist.RegisterComparer(flentities.Date{}, dateComparer)
	assist.RegisterComparer(time.Time{}, timeComparer)
	assist.RegisterComparer(&time.Time{}, timePtrComparer)
	assist.RegisterComparer(true, boolComparer)
	assist.RegisterComparer(uint(0), uintComparer)
	assist.RegisterComparer(refStr("a"), strPtrComparer)
	assist.RegisterComparer(refInt(1), intPtrComparer)
	assist.RegisterComparer(float32(1.0), floatComparer)
	assist.RegisterComparer(float64(1.0), floatComparer)
	assist.RegisterComparer(refUint(1), uintPtrComparer)
	assist.RegisterComparer(nil, nilComparer)
	assist.RegisterParser(uint(0), uintParser)
	assist.RegisterParser(true, boolParser)
	assist.RegisterParser(refStr("a"), strPtrParser)
	assist.RegisterParser(refInt(1), intPtrParser)
	assist.RegisterParser(refFloat(1), float32PtrParser)
	assist.RegisterParser(float32(1.0), float32Parser)
	assist.RegisterParser(refUint(1), uintPtrParser)
	assist.RegisterParser(&time.Time{}, timePtrParser)
	assist.RegisterParser(flentities.Date{}, dateParser)

	assist.RegisterParser(optional.String{}, func(raw string) (interface{}, error) {
		return optional.String{Value: &raw, Set: true}, nil
	})

	assist.RegisterParser(optional.Int{}, func(raw string) (interface{}, error) {
		i, err := strconv.Atoi(raw)
		return optional.Int{Value: &i, Set: true}, err
	})

	assist.RegisterParser(optional.Uint{}, func(raw string) (interface{}, error) {
		i, err := strconv.Atoi(raw)
		ui := uint(i)
		return optional.Uint{Value: &ui, Set: true}, err
	})

	assist.RegisterParser(optional.Float32{}, func(raw string) (interface{}, error) {
		i, err := strconv.ParseFloat(raw, 32)
		f := float32(i)
		return optional.Float32{Value: &f, Set: true}, err
	})
}

func AssistdogContext(s *godog.Suite) {
	s.BeforeSuite(setupAssistdog)
}

func timeComparer(raw string, rawActual interface{}) error {
	at, ok := rawActual.(time.Time)
	if !ok {
		return fmt.Errorf("%v is not a time.Time", rawActual)
	}

	et, err := defaults.ParseTime(raw)
	if err != nil {
		return err
	}

	expected := et.(time.Time).UTC()
	actual := at.UTC()
	if expected != actual {
		return fmt.Errorf("Expected %v, but got %v", expected, actual)
	}

	return nil
}

func dateComparer(raw string, rawActual interface{}) error {
	d, ok := rawActual.(flentities.Date)
	if !ok {
		return fmt.Errorf("%v is not a Date", rawActual)
	}

	if raw == "" && d.IsZero() {
		return nil
	}

	if raw != d.String() {
		return fmt.Errorf("Expected %s, got %s", raw, d.String())
	}

	return nil
}

func timePtrComparer(raw string, rawActual interface{}) error {
	at, ok := rawActual.(*time.Time)
	if !ok {
		return fmt.Errorf("%v is not *time.Time", rawActual)
	}

	if at == nil && raw == "" {
		return nil
	}

	if raw == "" || at == nil {
		if raw == "" {
			raw = "<nil>"
		}

		return fmt.Errorf("Expected %v, but got %v", raw, at)
	}

	return timeComparer(raw, *at)
}

func uintComparer(raw string, rawActual interface{}) error {
	rawInt, err := strconv.Atoi(raw)
	if err != nil {
		return err
	}

	if uint(rawInt) == rawActual.(uint) {
		return nil
	}

	return fmt.Errorf("Expected %d, but got %d", rawInt, rawActual)
}

func boolComparer(raw string, rawActual interface{}) error {
	actual := fmt.Sprint(rawActual)
	if raw == actual {
		return nil
	}

	return fmt.Errorf("Expected %s, but got %s", raw, actual)
}

func strPtrComparer(raw string, rawActual interface{}) error {
	actual := derefStr(rawActual.(*string))
	if raw == actual {
		return nil
	}

	return fmt.Errorf("Expected %s, but got %s", raw, actual)
}

func nilComparer(raw string, rawActual interface{}) error {
	if raw == "" {
		return nil
	}

	return fmt.Errorf("Expected <nil>, but got %s", rawActual)
}

func intPtrComparer(raw string, rawActual interface{}) error {
	actual := strconv.Itoa(derefInt(rawActual.(*int)))
	if raw == actual {
		return nil
	}

	return fmt.Errorf("Expected %s, but got %s", raw, actual)
}

func floatComparer(raw string, rawActual interface{}) error {
	expected, err := float64Parser(raw)
	if err != nil {
		return err
	}

	ex := expected.(float64)
	failed := false

	switch ac := rawActual.(type) {
	case float32:
		failed = float32(ex) != ac
	case float64:
		failed = ex != ac
	}

	if failed {
		return fmt.Errorf("Expected %s, but got %f", raw, rawActual)
	}

	return nil
}

func uintPtrComparer(raw string, rawActual interface{}) error {
	actual := ""
	actualUint := rawActual.(*uint)
	if actualUint != nil {
		actual = fmt.Sprint(*actualUint)
	}

	if raw == actual {
		return nil
	}

	return fmt.Errorf("Expected %s, but got %s", raw, actual)
}

func uintParser(raw string) (interface{}, error) {
	i, err := strconv.Atoi(raw)
	if err != nil {
		return nil, err
	}

	return uint(i), nil
}

func boolParser(raw string) (interface{}, error) {
	if raw != "true" && raw != "false" {
		return nil, fmt.Errorf("Don't know how to parse \"%s\" to bool", raw)
	}

	return raw == "true", nil
}

func strPtrParser(raw string) (interface{}, error) {
	if raw == "" {
		return nil, nil
	}

	return &raw, nil
}

func intPtrParser(raw string) (interface{}, error) {
	var i *int
	if raw == "" {
		return i, nil
	}

	in, err := strconv.Atoi(raw)
	if err != nil {
		return i, err
	}

	return &in, nil
}

func float32PtrParser(raw string) (interface{}, error) {
	f, err := float32Parser(raw)

	if err != nil {
		return nil, err
	}

	fl := f.(float32)

	return &fl, err
}

func float32Parser(raw string) (interface{}, error) {
	if raw == "" {
		return nil, nil
	}

	i, err := strconv.ParseFloat(raw, 32)
	if err != nil {
		return nil, err
	}

	f := float32(i)

	return f, nil
}

func float64PtrParser(raw string) (interface{}, error) {
	f, err := float64Parser(raw)

	if err != nil {
		return nil, err
	}

	fl := f.(float64)

	return &fl, err
}

func float64Parser(raw string) (interface{}, error) {
	var f float64
	if raw == "" {
		return f, nil
	}

	f, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return f, err
	}

	return f, nil
}

func uintPtrParser(raw string) (interface{}, error) {
	var u *uint
	if raw == "" {
		return u, nil
	}

	i, err := strconv.Atoi(raw)
	if err != nil {
		return u, err
	}

	ui := uint(i)
	return &ui, nil
}

func timePtrParser(raw string) (interface{}, error) {
	var t *time.Time
	if raw == "" {
		return t, nil
	}

	i, err := time.Parse(fullTimeFormat, raw)
	if err != nil {
		return t, err
	}

	return &i, nil
}

func dateParser(raw string) (interface{}, error) {
	if raw == "" {
		return flentities.Date{}, nil
	}

	return flentities.ParseDate(raw)
}

func tableColumn(table *gherkin.DataTable, idx int) []string {
	column := []string{}

	for _, row := range table.Rows {
		column = append(column, row.Cells[idx].Value)
	}

	return column
}
