package runtime

import "time"

func dateStamp() string {
	return time.Now().Format("2006-01-02")
}
