package flemail

import (
	"bytes"
	"reflect"
	txt "text/template"

	"github.com/nerde/fuji-lane-back/fujilane"
)

var tmpl *txt.Template

func renderTextTemplate(data interface{}) (body string, err error) {
	if tmpl == nil {
		tmpl, err = txt.ParseGlob(fujilane.Root() + "/flemail/templates/*.text")
		if err != nil {
			return
		}
	}

	buff := &bytes.Buffer{}
	tmpl.ExecuteTemplate(buff, fujilane.ToSnake(reflect.TypeOf(data).Name())+".text", data)
	body = buff.String()

	return
}
