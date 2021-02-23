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
				hour := t.In(tz).Hour()
				minute := t.In(tz).Minute()
				if (hour == 9 && minute >= 30) ||
					(hour > 9 && hour < 11) ||
					(hour == 11 && minute <= 30) ||
					(hour == 13 && minute >= 1) ||
					(hour > 13 && hour < 15) ||
					(hour == 15 && minute == 0) {
					record()
				}
			}
		}
	}
}
