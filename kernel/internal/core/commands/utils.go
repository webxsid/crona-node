package commands

import "time"

func elapsedSeconds(startTime string, now string) int {
	start, err := time.Parse(time.RFC3339, startTime)
	if err != nil {
		return 0
	}
	end, err := time.Parse(time.RFC3339, now)
	if err != nil {
		return 0
	}
	seconds := int(end.Sub(start).Seconds())
	if seconds < 0 {
		return 0
	}
	return seconds
}
