package sender

import (
	"bytes"
	"html/template"
	"log"
	"mailing/internal/dbsubs"
	"net/smtp"
	"strings"
)

var auth = smtp.PlainAuth("", "andrey.aksenov2001@gmail.com", "password", "smtp.gmail.com")
var headers = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";"

func SendEmail(sendtime string, tmpl string, subs []dbsubs.Subscriber) {
	if sendtime == "now" {
		log.Println("Sending...")
		subject := strings.ToUpper(tmpl)
		tmplPath := "templates/msg/" + tmpl + ".html"
		t, err := template.ParseFiles(tmplPath)
		if err != nil {
			log.Println(err)
		}

		msg := "Subject: " + subject + "\n" + headers + "\n\n"

		for _, sub := range subs {
			go sender(t, msg, sub)
		}

	} else {
		log.Println("Later")
	}
}

func sender(t *template.Template, msg string, sub dbsubs.Subscriber) {
	var body bytes.Buffer

	err := t.Execute(&body, sub)
	if err != nil {
		log.Println(err)
	}

	msg = msg + body.String()
	err = smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		"andrey.aksenov2001@gmail.com",
		[]string{sub.Email},
		[]byte(msg),
	)
	if err != nil {
		log.Println(err)
	}
}
