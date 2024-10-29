package utils

import (
	"fmt"
	"time"
)

func GetCurrentTimeInEST() (int64, error) {
	var est int64

	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		return est, err
	}

	est = time.Now().In(loc).Unix()

	return est, nil
}

func ParseDateInLocation(timestamp int64) (int64, error) {
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		return 0, fmt.Errorf("error loading location: %w", err)
	}

	localTime := time.Unix(timestamp, 0).In(loc)
	return localTime.Unix(), nil
}
