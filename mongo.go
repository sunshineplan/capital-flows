package main

import (
	"fmt"
	"time"

	"github.com/sunshineplan/database/mongodb"
	"github.com/sunshineplan/database/mongodb/driver"
	"github.com/sunshineplan/stock/capitalflows"
	"github.com/sunshineplan/utils/retry"
)

var client driver.Client

func initDB() error {
	if err := retry.Do(func() error {
		return meta.Get("capitalflows_mongo", &client)
	}, 3, 20); err != nil {
		return err
	}
	return client.Connect()
}

func record() {
	flows, err := capitalflows.Fetch()
	if err != nil {
		svc.Debug(err.Error())
		return
	}

	now := time.Now()
	res, err := client.UpdateOne(
		struct {
			Date string `json:"date" bson:"date"`
			Time string `json:"time" bson:"time"`
		}{
			now.Format("2006-01-02"),
			now.Format("15:04"),
		},
		mongodb.M{"$set": mongodb.M{"flows": flows}},
		&mongodb.UpdateOpt{Upsert: true},
	)
	if err != nil {
		svc.Debug(err.Error())
		return
	}

	if n := res.MatchedCount; n != 0 {
		svc.Debug(fmt.Sprintf("Updated %d record", n))
	}
	if n := res.UpsertedCount; n != 0 {
		svc.Debug(fmt.Sprintf("Upserted %d record", n))
	}
}
