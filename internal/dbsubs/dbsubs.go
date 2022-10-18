package dbsubs

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

type Instance struct {
	Db *pgxpool.Pool
}

type Subscriber struct {
	Email     string
	Firstname string
	Lastname  string
	Birthday  string
}

func (i *Instance) Start() {
	log.Println("Db...")
	//i.Db.Exec(context.Background(), "DROP TABLE subscribers;")
	//err := i.createTable(context.Background())
	//if err != nil {
	//	log.Fatal(err)
	//}

	//subs := []Subscriber{{
	//	Email:     "test@gmail.com",
	//	Firstname: "Kolya",
	//	Lastname:  "Ivanov",
	//	Birthday:  "2003-02-01",
	//}, {
	//	Email:     "test2@gmail.com",
	//	Firstname: "Dima",
	//	Lastname:  "Bublikov",
	//	Birthday:  "1999-07-30",
	//}, {
	//	Email:     "aksenovandrey4@gmail.com",
	//	Firstname: "Andrey",
	//	Lastname:  "Aksenov",
	//	Birthday:  "2001-08-11",
	//}, {
	//	Email:     "andrey.aksenov2001@gmail.com",
	//	Firstname: "Andrey",
	//	Lastname:  "Aksenov",
	//	Birthday:  "2001-08-11",
	//},
	//}
	//
	//for _, sub := range subs {
	//	i.addSubscriber(context.Background(), sub)
	//}

	//i.GetSubscribers(context.Background())

}

func (i *Instance) createTable(ctx context.Context) error {
	log.Println("Creating table...")
	_, err := i.Db.Exec(ctx, "CREATE TABLE subscribers ("+
		"email VARCHAR ( 255 ) UNIQUE NOT NULL,"+
		"firstname VARCHAR ( 50 ) NOT NULL,"+
		"lastname VARCHAR ( 50 ) NOT NULL,"+
		"birthday VARCHAR(10) NOT NULL);")
	if err != nil {
		return err
	}
	return nil
}

func (i *Instance) addSubscriber(ctx context.Context, sub Subscriber) error {
	//defer wg.Done()
	log.Println("Adding subscriber...")
	_, err := i.Db.Exec(ctx, "INSERT INTO subscribers (email, firstname, lastname, birthday) VALUES ($1, $2, $3, $4);", sub.Email, sub.Firstname, sub.Lastname, sub.Birthday)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (i *Instance) GetSubscribers(ctx context.Context) []Subscriber {
	var subs []Subscriber
	rows, err := i.Db.Query(ctx, "SELECT email, firstname, lastname, birthday FROM subscribers;")
	if err == pgx.ErrNoRows {
		fmt.Println("No rows")
		return nil
	} else if err != nil {
		fmt.Println(err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		sub := Subscriber{}
		rows.Scan(&sub.Email, &sub.Firstname, &sub.Lastname, &sub.Birthday)
		subs = append(subs, sub)
	}

	return subs
}
