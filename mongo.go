package main

import (
	"context"
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

	if _, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, tz)},
		bson.M{"$set": bson.M{"flows": flows}},
		options.Update().SetUpsert(true),
	); err != nil {
		return err
	}

	return nil
}
