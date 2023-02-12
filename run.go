package main

import (
	"fmt"
	"log"
	"time"

	"github.com/sunshineplan/stock/capitalflows"
	"github.com/sunshineplan/utils/scheduler"
)

func test() (err error) {
	_, e1 := capitalflows.Fetch()
	if e1 != nil {
		fmt.Println("Failed to fetch capital flows data:", e1)
	}

	e2 := initDB()
	if e2 != nil {
		fmt.Println("Failed to initialize mongodb:", e2)
	}

	if e1 != nil || e2 != nil {
		err = fmt.Errorf("test is failed")
	}

	return
}

func run() {
	if err := initDB(); err != nil {
		log.Fatal(err)
	}

	if isTrading(time.Now()) {
		go record()
	}

	if err := scheduler.NewScheduler().
		At(scheduler.Second(5), scheduler.Second(20), scheduler.Second(35), scheduler.Second(50)).
		Do(func(t time.Time) {
			if isTrading(t) {
				go record()
			}
		}); err != nil {
		log.Fatal(err)
	}
	select {}
}

func isTrading(t time.Time) bool {
	if weekday := t.In(tz).Weekday(); weekday >= 1 && weekday <= 5 {
		hour := t.In(tz).Hour()
		minute := t.In(tz).Minute()
		if (hour == 9 && minute >= 30) ||
			(hour > 9 && hour < 11) ||
			(hour == 11 && minute <= 30) ||
			(hour == 13 && minute >= 1) ||
			(hour > 13 && hour < 15) ||
			(hour == 15 && minute == 0) {
			return true
		}
	}

	return false
}
