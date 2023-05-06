package cron

func init() {
	ScheduleTask(StartAGiveaway, "", 0, "", "00 14 * * SAT")
}
