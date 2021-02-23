package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/sunshineplan/stock/capitalflows"
	"github.com/sunshineplan/utils/database/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var config mongodb.Config
var collection *mongo.Collection

func initDB() error {
	if err := meta.Get("capitalflows_mongo", &config); err != nil {
		return err
	}

	client, err := config.Open()
	if err != nil {
		return err
	}

	collection = client.Database(config.Database).Collection(config.Collection)

	return nil
}

func record() error {
	flows, err := capitalflows.Fetch()
	if err != nil {
		return err
	}

	t := time.Now().In(tz)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := collection.UpdateOne(
		ctx,
		bson.M{
			"date": fmt.Sprintf("%04d-%02d-%02d", t.Year(), t.Month(), t.Day()),
			"time": fmt.Sprintf("%02d:%02d", t.Hour(), t.Minute()),
		},
		bson.M{"$set": bson.M{"flows": flows}},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		return err
	}

	if n := res.MatchedCount; n != 0 {
		log.Printf("Updated %d record", n)
	}
	if n := res.UpsertedCount; n != 0 {
		log.Printf("Upserted %d record", n)
	}

	return nil
}
