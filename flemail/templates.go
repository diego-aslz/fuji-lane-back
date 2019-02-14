package flemail

import (
	"bytes"
	html "html/template"
	"reflect"
	text "text/template"

	"github.com/nerde/fuji-lane-back/fujilane"
)

var textTmpl = text.Must(text.ParseGlob(fujilane.Root() + "/flemail/templates/*.text"))
var htmlTmpl = html.Must(html.ParseGlob(fujilane.Root() + "/flemail/templates/*.html"))

func renderTextTemplate(data interface{}) (string, error) {
	return renderTemplate(data, "text")
}

func renderHTMLTemplate(data interface{}) (string, error) {
	return renderTemplate(data, "html")
}

func renderTemplate(data interface{}, ext string) (string, error) {
	buff := &bytes.Buffer{}

	fn := textTmpl.ExecuteTemplate
	if ext == "html" {
		fn = htmlTmpl.ExecuteTemplate
	}

	err := fn(buff, fujilane.ToSnake(reflect.TypeOf(data).Name())+"."+ext, data)

	return buff.String(), err
}
