package main

import (
	"L0/db"
	"flag"
	_ "github.com/lib/pq"
	"github.com/nats-io/stan.go"
	"log"
	"net/http"
)

var conn = flag.String("conn", "postgres://mokyuser:pwd4moky@localhost/mokydb?sslmode=disable",
	"database connection string: postgres://{user}:{password}@localhost/{dbname}?sslmode=disable")

func main() {
	nc, err := stan.Connect("test", "test-producer")
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	dbh, err := db.Connect(conn)
	if err != nil {
		log.Fatal(err)
	}
	defer dbh.Close()

	err = dbh.InitTable()
	if err != nil {
		log.Fatal(err)
	}

	nc.Subscribe("test-sub", func(m *stan.Msg) {
		dbh.Insert(m.Data)
	}, stan.StartWithLastReceived())

	http.ListenAndServe(":8888", nil)
}
