package main

import (
	"fmt"
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

func run() error {
	if err := initDB(); err != nil {
		return err
	}

	if err := scheduler.NewScheduler().AtCondition(
		scheduler.Weekdays,
		scheduler.MultiSchedule(
			scheduler.ClockSchedule(scheduler.AtClock(9, 30, 0), scheduler.AtClock(11, 30, 0), 15*time.Second),
			scheduler.ClockSchedule(scheduler.AtClock(13, 0, 1), scheduler.AtHour(15), 15*time.Second),
		),
	).Do(func(_ time.Time) {
		record()
	}); err != nil {
		return err
	}

	scheduler.Forever()
	return nil
}
