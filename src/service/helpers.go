package main

import (
	"L0/db"
	"encoding/json"
	"errors"
	"fmt"
	"log"
)

func msgHandler(msg []byte) (string, error) {
	buf := make(map[string]any)

	err := json.Unmarshal(msg, &buf)
	if err != nil {
		return "", errors.New(fmt.Sprintf("Error: Message: %s is invalid ", string(msg)))
	}

	uid, ok := buf["order_uid"]
	if !ok {
		return "", errors.New(fmt.Sprintf("Error: Message: %s is invalid ", string(msg)))
	}

	return uid.(string), nil
}

func uploadCache() map[string][]byte {
	dbh, err := db.Connect(conn)
	if err != nil {
		log.Fatal(err)
	}
	defer dbh.Close()

	err = dbh.InitTable()
	if err != nil {
		log.Fatal(err)
	}

	return dbh.Get()
}

func saveMsg(uid string, msg []byte) error {
	_, ok := cache[uid]
	if ok {
		return errors.New(fmt.Sprintf("Error: order withd Uid %s already exists ", uid))
	}

	dbh, err := db.Connect(conn)
	if err != nil {
		return err
	}
	defer dbh.Close()

	err = dbh.Insert(uid, msg)
	if err != nil {
		return err
	}

	cache[uid] = msg

	return nil
}
