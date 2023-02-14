package main

import (
	"flag"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/nats-io/stan.go"
	"log"
	"net/http"
	"time"
)

var conn = flag.String("conn", "postgres://mokyuser:pwd4moky@localhost/mokydb?sslmode=disable",
	"database connection string: postgres://{user}:{password}@localhost/{dbname}?sslmode=disable")

var cache = uploadCache()

func main() {
	addr := flag.String("addr", ":8888", "Сетевой адрес HTTP")
	flag.Parse()
	nc, err := stan.Connect("test", "test-producer")
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	nc.Subscribe("test-sub", func(m *stan.Msg) {
		uid, err := msgHandler(m.Data)
		if err != nil {
			log.Println(err)
			return
		}

		saveMsg(uid, m.Data)
	})

	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/order", uidHandler)

	server := http.Server{
		Addr:         *addr,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Printf("Start listening at http://localhost%s\n", *addr)
	server.ListenAndServe()
}
