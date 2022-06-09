package timehelpers

import "time"

const MILLISEC_TO_NANOSEC = 1000000

func NewStopWatch() StopWatch {
	return StopWatch{
		time.Now(),
	}
}

func (sw *StopWatch) EllapsedMillis() int {
	return int(time.Since(sw.Start).Nanoseconds()) / MILLISEC_TO_NANOSEC
}
