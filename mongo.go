package main

import (
	"fmt"
	"log"
	"time"

	"github.com/sunshineplan/database/mongodb"
	"github.com/sunshineplan/database/mongodb/driver"
	"github.com/sunshineplan/stock/capitalflows"
	"github.com/sunshineplan/utils"
)

var client mongodb.Client

func initDB() error {
	var mongo driver.Client
	if err := utils.Retry(func() error {
		return meta.Get("capitalflows_mongo", &mongo)
	}, 3, 20); err != nil {
		return err
	}
	client = &mongo

	return client.Connect()
}

func record() {
	flows, err := capitalflows.Fetch()
	if err != nil {
		if debug {
			log.Print(err)
		}

		return
	}

	t := time.Now().In(tz)

	res, err := client.UpdateOne(
		struct {
			Date string `json:"date" bson:"date"`
			Time string `json:"time" bson:"time"`
		}{
			fmt.Sprintf("%04d-%02d-%02d", t.Year(), t.Month(), t.Day()),
			fmt.Sprintf("%02d:%02d", t.Hour(), t.Minute()),
		},
		mongodb.M{"$set": mongodb.M{"flows": flows}},
		&mongodb.UpdateOpt{Upsert: true},
	)
	if err != nil {
		if debug {
			log.Print(err)
		}

		return
	}

	if n := res.MatchedCount; n != 0 && debug {
		log.Printf("Updated %d record", n)
	}
	if n := res.UpsertedCount; n != 0 && debug {
		log.Printf("Upserted %d record", n)
	}
}
