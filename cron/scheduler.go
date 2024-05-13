package cron

import (
	"time"

	"github.com/go-co-op/gocron"
)

var (
	S = gocron.NewScheduler(time.UTC)
)

/*
Schedule a job to be run on the given schedule.
...args are in the order of
    `redoflow(second, day, etc.), momentum(integer),
	 realtime("10:30", "15:00", etc.), crontime(* * * * *),` */
func ScheduleTask(task interface{}, args ...interface{}) {
	// extract arg from args
	redoflow := args[0].(string) // seconds, minutes, hours, days, weeks, months
	momentum := args[1].(int)    // integer
	realtime := args[2].(string) // 10:30, 15:00, etc.
	crontime := args[3].(string) // * * * * *

	switch redoflow {
	case "second":
		S.Every(momentum).Seconds().Do(task)
	case "minute":
		S.Every(momentum).Minutes().Do(task)
	case "hour":
		S.Every(momentum).Hours().Do(task)
	case "day":
		if realtime != "" {
			S.Every(momentum).Days().At(realtime).Do(task)
		} else {
			S.Every(momentum).Days().Do(task)
		}
	case "week":
		if realtime != "" {
			S.Every(momentum).Weeks().At(realtime).Do(task)
		} else {
			S.Every(momentum).Weeks().Do(task)
		}
	case "month":
		if realtime != "" {
			S.Every(momentum).Months().At(realtime).Do(task)
		} else {
			S.Every(momentum).Months().Do(task)
		}
	default:
		S.Cron(crontime).Do(task)
	}
}
