package main

import (
	"net/textproto"

	"github.com/nerde/fuji-lane-back/flconfig"
	"github.com/nerde/fuji-lane-back/flservices"

	"github.com/jordan-wright/email"
)

func main() {
	e := &email.Email{
		To:      []string{"diegoselzlein@gmail.com"},
		From:    "Jordan Wright <test@gmail.com>",
		Subject: "Awesome Subject",
		Text:    []byte("Text Body is, of course, supported!"),
		HTML:    []byte("<h1>Fancy HTML is supported, too!</h1>"),
		Headers: textproto.MIMEHeader{},
	}

	flconfig.LoadConfiguration()

	if err := flservices.NewSMTPMailer().Send(e); err != nil {
		panic(err)
	}
}
