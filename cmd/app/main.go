package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"mailing/internal/dbsubs"
	"mailing/internal/sender"
	"mailing/pkg/helpers/pg"
	"os"
	"strconv"
	"strings"
	"sync"
)

var wg sync.WaitGroup

// MessageTracker provides the amount of opening message with the given messageId and list of subscribers that have opened this one
var MessageTracker = make(map[int]Track)

type Track struct {
	OpenedTimes int      `json:"opened"`
	Openers     []string `json:"openers"`
}

func main() {
	defer wg.Wait()
	cfg := &pg.Config{}
	cfg.Host = "localhost"
	cfg.Username = "db_user"
	cfg.Password = "pwd123"
	cfg.Port = "54320"
	cfg.DbName = "db_subs"
	cfg.Timeout = 5

	poolConfig, err := pg.NewPoolConfig(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Pool config error: %v\n", err)
		os.Exit(1)
	}

	poolConfig.MaxConns = 5

	conn, err := pg.NewConnection(poolConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Connection to database failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Connection OK!")

	_, err = conn.Exec(context.Background(), ";")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ping failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Ping OK!")

	ins := &dbsubs.Instance{Db: conn}
	subs := ins.GetSubscribers(context.Background())

	router := gin.Default()

	router.LoadHTMLGlob("templates/**/*")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		if c.Query("messageId") != "" {
			messageId, _ := strconv.Atoi(c.Query("messageId"))

			c.HTML(200, "Home.html", gin.H{
				"opened":  MessageTracker[messageId].OpenedTimes,
				"openers": strings.Join(MessageTracker[messageId].Openers, ", "),
			})
		} else {
			c.HTML(200, "Home.html", gin.H{})
		}
	})

	router.GET("/Messages", func(c *gin.Context) {
		c.HTML(200, "Messages.html", gin.H{})
	})

	router.GET("/Contacts", func(c *gin.Context) {
		c.HTML(200, "Contacts.html", gin.H{})
	})

	router.GET("/Message/:msg", func(c *gin.Context) {
		msg := c.Param("msg")
		sendtime := c.Query("sendtime")
		if sendtime != "" {
			// if we would send message we need to create a new id for tracking that one
			messageId := len(MessageTracker) + 1
			MessageTracker[messageId] = Track{
				OpenedTimes: 0,
				Openers:     []string{},
			}
			c.String(200, "The message will be sent at time: "+sendtime+
				"\nmessageId for tracking: "+strconv.Itoa(messageId))

			sender.SendEmail(messageId, sendtime, msg, subs, &wg)

		} else {
			c.HTML(200, "Message.html", gin.H{
				"Message": "/" + msg,
			})
		}

	})

	router.GET("/welcome", func(c *gin.Context) {
		c.HTML(200, "welcome.html", gin.H{
			"Email":     subs[0].Email,
			"Firstname": subs[0].Firstname,
			"Lastname":  subs[0].Lastname,
			"Birthday":  subs[0].Birthday,
			"messageId": 0,
		})
	})

	router.GET("/gift", func(c *gin.Context) {
		c.HTML(200, "gift.html", gin.H{
			"Email":     subs[0].Email,
			"Firstname": subs[0].Firstname,
			"Lastname":  subs[0].Lastname,
			"Birthday":  subs[0].Birthday,
			"messageId": 0,
		})
	})

	// Announcing that some subscriber have opened the message with certain id
	router.GET("/tracker/:tr", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Query("id"))
		if err == nil && id > 0 {
			log.Println("tracked id = " + c.Query("id") + " sub = " + c.Query("sub"))
			MessageTracker[id] = Track{
				OpenedTimes: MessageTracker[id].OpenedTimes + 1,
				Openers:     append(MessageTracker[id].Openers, c.Query("sub")),
			}
		}

	})

	router.Run()
}
