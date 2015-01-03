package platypus

import "time"

func getLocalTimeLocation() *time.Location {
	t := time.Now()
	local := t.Local()
	return local.Location()
}
