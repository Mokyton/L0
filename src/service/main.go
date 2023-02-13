package main

import (
	"L0/db"
	"flag"
	"fmt"
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
	nc.Subscribe("test-sub", func(m *stan.Msg) {
		fmt.Printf("Got: %s\n", string(m.Data))
	})

	http.ListenAndServe(":8888", nil)
	//msg.Respond([]byte("Krasivo"))
	//dbh, err := db.Connect(conn)
	//if err != nil {
	//	return
	//}
	//defer dbh.Close()
	//run(msg.Data)
	//http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//	w.Write(dbh.Get())
	//})
	//http.ListenAndServe(":8888", nil)
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
