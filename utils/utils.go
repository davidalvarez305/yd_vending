package utils

import "time"

func GetCurrentTimeInEST() (int64, error) {
	var est int64

	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		return est, err
	}

	est = time.Now().In(loc).Unix()

	return est, nil
}
