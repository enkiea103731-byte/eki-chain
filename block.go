package main
import "time"
var startTime = time.Now()
func CurrentBlock() int64 { return int64(1000 + int(time.Since(startTime).Seconds()/10)) }
