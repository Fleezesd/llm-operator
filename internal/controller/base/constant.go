package base

import "time"

// wait time duration
const (
	waitLonger  = time.Hour
	waitSmaller = time.Second * 3
	waitMedium  = time.Minute
)

const (
	statusNilResponse = "No err replied but response is not string"
)
