package utils

import (
	"fmt"
	"strconv"
	"strings"
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

func UrlsListHasCurrentPath(urls []string, url string) bool {
	for _, protectedUrl := range urls {
		if strings.Contains(url, protectedUrl) {
			return true
		}
	}
	return false
}

func GetStartAndEndDatesFromMonthYear(monthYear string) (time.Time, time.Time, error) {
	parts := strings.Split(monthYear, ", ")
	if len(parts) != 2 {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid format")
	}

	selectedMonth := parts[0]
	selectedYear := parts[1]

	year, err := strconv.Atoi(selectedYear)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid year: %v", err)
	}

	monthTime, err := time.Parse("January", selectedMonth)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid month: %v", err)
	}

	start := time.Date(year, monthTime.Month(), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)

	return start, end, nil
}
