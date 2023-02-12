package main

import (
	"L0/db"
	"flag"
	_ "github.com/lib/pq"
	"github.com/nats-io/nats.go"
	"log"
	"time"
)

var conn = flag.String("conn", "postgres://mokyuser:pwd4moky@localhost/mokydb?sslmode=disable", "database connection string")

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	// Subscribe
	sub, err := nc.SubscribeSync("foo")
	if err != nil {
		log.Fatal(err)
	}

	// Wait for a message
	msg, err := sub.NextMsg(10 * time.Second)
	if err != nil {
		log.Fatal(err)
	}
	msg.Respond([]byte("Krasivo"))
	run(msg.Data)

}

func run(d []byte) error {
	dbh, err := db.Connect(conn)
	if err != nil {
		return err
	}
	defer dbh.Close()
	_ = dbh.InitTable()
	_ = dbh.Insert(d)
	return nil
}
