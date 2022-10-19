package sender

import (
	"bytes"
	"html/template"
	"log"
	"mailing/internal/dbsubs"
	"net/smtp"
	"strings"
	"sync"
	"time"
)

var auth = smtp.PlainAuth("", "andrey.aksenov2001@gmail.com", "hdcirobywejtbibo", "smtp.gmail.com")
var headers = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";"

func SendEmail(sendtime string, tmpl string, subs []dbsubs.Subscriber, wg *sync.WaitGroup) {

	subject := strings.ToUpper(tmpl)
	tmplPath := "templates/msg/" + tmpl + ".html"
	t, err := template.ParseFiles(tmplPath)
	if err != nil {
		log.Println(err)
	}

	msg := "Subject: " + subject + "\n" + headers + "\n\n"
	if sendtime == "now" {
		//for _, sub := range subs {
		//	wg.Add(1)
		//	go sendToOne(t, msg, sub, wg)
		//}
		sendToAll(t, msg, subs, wg)

	} else {
		sendtime = strings.Replace(sendtime, "T", " ", -1)
		sendtime = strings.Replace(sendtime, "%3A", ":", -1)
		log.Println("message will be sent at time: " + sendtime)
		d, err := time.Parse("2006-01-02 15:04", sendtime)
		d = d.Add(-3 * time.Hour)
		if err != nil {
			log.Println(err)
		}

		go func(d time.Time, t *template.Template, msg string, subs []dbsubs.Subscriber, wg *sync.WaitGroup) {
			time.Sleep(time.Until(d))
			sendToAll(t, msg, subs, wg)
		}(d, t, msg, subs, wg)

	}
}

func sendToOne(t *template.Template, msg string, sub dbsubs.Subscriber, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Println("Sending...")
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

func sendToAll(t *template.Template, msg string, subs []dbsubs.Subscriber, wg *sync.WaitGroup) {
	for _, sub := range subs {
		wg.Add(1)
		go sendToOne(t, msg, sub, wg)
	}
}
