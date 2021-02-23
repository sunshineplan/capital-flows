package main

import (
	"log"
	"time"
)

func run() {
	if err := initDB(); err != nil {
		log.Fatal(err)
	}

	ticker := time.NewTicker(time.Nanosecond)
Reset:
	for {
		select {
		case t := <-ticker.C:
			if t.Second() == 46 {
				ticker.Stop()
				break Reset
			}
		}
	}

	ticker = time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case t := <-ticker.C:
			if weekday := t.In(tz).Weekday(); weekday >= 1 && weekday <= 5 {
				if hour := t.In(tz).Hour(); hour >= 9 && hour <= 18 {
					if err := record(); err != nil {
						log.Print(err)
					}
				}
			}
		}
	}
}
