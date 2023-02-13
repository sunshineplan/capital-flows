package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/sunshineplan/database/mongodb"
	"github.com/sunshineplan/database/mongodb/driver"
	"github.com/sunshineplan/stock/capitalflows/sector"
)

var client = driver.Client{
	Database:   "stock",
	Collection: "capitalflows",
	Username:   "capitalflows",
	Password:   "capitalflows",
	SRV:        true,
}

var (
	path = flag.String("path", "", "data path")
	cmt  = flag.Bool("commit", false, "commit records")
	del  = flag.Bool("delete", false, "delete records")
)

var (
	tz    = time.FixedZone("CST", 8*60*60)
	now   = time.Now().In(tz)
	today = now.Format("2006-01-02")
)

type date struct {
	Date string `json:"_id" bson:"_id"`
}

func main() {
	flag.StringVar(&client.Server, "mongo", "", "MongoDB Server")
	flag.Parse()

	if err := client.Connect(); err != nil {
		log.Fatal(err)
	}

	var res []date
	if err := client.Aggregate([]mongodb.M{{"$group": mongodb.M{"_id": "$date"}}}, &res); err != nil {
		log.Fatal(err)
	}

	if *cmt {
		if err := commit(res); err != nil {
			log.Fatal(err)
		}
	} else if *del {
		if err := delete(res); err != nil {
			log.Fatal(err)
		}
	}
}

func commit(res []date) error {
	for _, i := range res {
		if i.Date != today {
			res, err := sector.GetTimeLine(i.Date, &client)
			if err != nil {
				return err
			}

			if res[0].TimeLine[0]["09:30"] != res[0].TimeLine[len(res[0].TimeLine)-1]["15:00"] {
				b, err := json.Marshal(res)
				if err != nil {
					return err
				}

				fullpath := filepath.Join(append([]string{*path}, strings.Split(i.Date, "-")...)...) + ".json"
				if err := os.MkdirAll(filepath.Dir(fullpath), 0755); err != nil {
					return err
				}
				if err := os.WriteFile(fullpath, b, 0644); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func delete(res []date) error {
	for _, i := range res {
		if i.Date != today {
			d, _ := time.ParseInLocation("2006-01-02", i.Date, tz)
			if now.Sub(d).Hours() > 7*24 {
				if _, err := client.DeleteMany(mongodb.M{"date": i.Date}); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
