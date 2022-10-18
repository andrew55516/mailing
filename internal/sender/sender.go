package sender

import (
	"bytes"
	"github.com/gocelery/gocelery"
	"github.com/gomodule/redigo/redis"
	"html/template"
	"log"
	"mailing/internal/dbsubs"
	"net/smtp"
	"os"
	"strings"
	"time"
)

var auth = smtp.PlainAuth("", "andrey.aksenov2001@gmail.com", "password", "smtp.gmail.com")
var headers = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";"

// create redis connection pool
var redisPool = &redis.Pool{
	Dial: func() (redis.Conn, error) {
		c, err := redis.DialURL(os.Getenv("REDIS_URL"))
		if err != nil {
			return nil, err
		}
		return c, err
	},
}

// initialize celery client
var cli, _ = gocelery.NewCeleryClient(
	gocelery.NewRedisBroker(redisPool),
	&gocelery.RedisCeleryBackend{Pool: redisPool},
	5, // number of workers
)

var taskName = "worker.send"

func SendEmail(sendtime string, tmpl string, subs []dbsubs.Subscriber) {

	subject := strings.ToUpper(tmpl)
	tmplPath := "templates/msg/" + tmpl + ".html"
	t, err := template.ParseFiles(tmplPath)
	if err != nil {
		log.Println(err)
	}

	msg := "Subject: " + subject + "\n" + headers + "\n\n"
	if sendtime == "now" {
		for _, sub := range subs {
			go sender(t, msg, sub)
		}

	} else {
		log.Println("Later")
		d, _ := time.Parse("2006-01-02 15:04", sendtime)
		//d := t.Unix() - time.Now().Unix()
		cli.Register("worker.send", sendWithDelay)
		_, err = cli.Delay(taskName, d, t, msg, subs)
		if err != nil {
			log.Println(err)
		}
	}
}

func sender(t *template.Template, msg string, sub dbsubs.Subscriber) {
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

func sendWithDelay(d time.Time, t *template.Template, msg string, subs []dbsubs.Subscriber) {
	log.Println(time.Until(d))
	time.Sleep(time.Until(d))
	for _, sub := range subs {
		sender(t, msg, sub)
	}
}
