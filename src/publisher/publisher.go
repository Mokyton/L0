package main

import (
	"flag"
	"github.com/nats-io/stan.go"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	path := flag.String("path", "./publisher/data-sets", "path to data-sets dir")
	flag.Parse()
	nc, err := stan.Connect("test", "test-publisher")
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()
	files, err := ioutil.ReadDir(*path)
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range files {
		data, err := readData(*path + "/" + v.Name())
		if err != nil {
			return
		}
		nc.Publish("test-sub", data)
	}

}

func readData(path string) ([]byte, error) {
	file, err := os.Open(path)
	buf, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return buf, nil
}
