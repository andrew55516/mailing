package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"mailing/internal/dbsubs"
	"mailing/internal/sender"
	"mailing/pkg/helpers/pg"
	"os"
	"sync"
)

var wg sync.WaitGroup

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
		c.HTML(200, "Home.html", gin.H{})
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
			c.String(200, "The message will be sent at time: "+sendtime)

			sender.SendEmail(sendtime, msg, subs, &wg)

		} else {
			c.HTML(200, "Message.html", gin.H{
				"Message": "/" + msg,
			})
		}

	})

	router.GET("/welcome", func(c *gin.Context) {
		c.HTML(200, "welcome.html", subs[0])
	})

	router.GET("/gift", func(c *gin.Context) {
		c.HTML(200, "gift.html", subs[0])
	})

	router.Run()
}
