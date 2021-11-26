package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/go-github/v35/github"
	"github.com/sunshineplan/database/mongodb"
	"github.com/sunshineplan/database/mongodb/driver"
	"github.com/sunshineplan/stock/capitalflows/sector"
	"golang.org/x/oauth2"
)

var mongo = driver.Client{
	Database:   "stock",
	Collection: "capitalflows",
	Username:   "capitalflows",
	Password:   "capitalflows",
	SRV:        true,
}
var client mongodb.Client

var token, repository, path string

func main() {
	flag.StringVar(&mongo.Server, "mongo", "", "MongoDB Server")
	flag.StringVar(&token, "token", "", "token")
	flag.StringVar(&repository, "repo", "", "repository")
	flag.StringVar(&path, "path", "", "data path")
	flag.Parse()

	client = &mongo
	if err := client.Connect(); err != nil {
		log.Fatal(err)
	}

	if err := commit(); err != nil {
		log.Fatal(err)
	}
}

func commit() error {
	var date []struct {
		Date string `json:"_id" bson:"_id"`
	}
	if err := client.Aggregate([]mongodb.M{{"$group": mongodb.M{"_id": "$date"}}}, &date); err != nil {
		return err
	}

	tz, _ := time.LoadLocation("Asia/Shanghai")
	t := time.Now().In(tz)
	today := fmt.Sprintf("%d-%02d-%02d", t.Year(), t.Month(), t.Day())

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})

	for _, i := range date {
		if i.Date != today {
			res, err := sector.GetTimeLine(i.Date, client)
			if err != nil {
				return err
			}

			if res[0].TimeLine[0]["09:30"] != res[0].TimeLine[len(res[0].TimeLine)-1]["15:00"] {
				b, err := json.Marshal(res)
				if err != nil {
					return err
				}

				ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
				defer cancel()

				tc := oauth2.NewClient(ctx, ts)
				client := github.NewClient(tc)
				repo := strings.Split(repository, "/")
				fullpath := filepath.Join(append([]string{path}, strings.Split(i.Date, "-")...)...) + ".json"
				opt := &github.RepositoryContentFileOptions{
					Message: github.String(i.Date),
					Content: b,
				}

				if _, _, err := client.Repositories.CreateFile(ctx, repo[0], repo[1], fullpath, opt); err != nil {
					if !strings.Contains(err.Error(), `"sha" wasn't supplied.`) {
						return err
					}
				}
			}

			d, _ := time.ParseInLocation("2006-01-02", i.Date, tz)
			if t.Sub(d).Hours() > 7*24 {
				if _, err := client.DeleteMany(mongodb.M{"date": i.Date}); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
