package main

import (
	"log"
	"time"
)

func run() {
	if err := initDB(); err != nil {
		log.Fatal(err)
	}

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case t := <-ticker.C:
			if weekday := t.In(tz).Weekday(); weekday >= 1 && weekday <= 5 {
				if hour := t.In(tz).Hour(); hour >= 9 && hour <= 18 {
					record()
				}
			}
		}
	}
}
