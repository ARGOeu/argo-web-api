package utils

import (
	"strconv"
	"time"
)

const zuluDateOnly = "2006-01-02"
const ymdForm = "20060102"

// ParseZuluDate is used to parse a zulu formatted date to integer. If no date is entered
// ParseZuluDate returns a Zulu representation of current time
func ParseZuluDate(dateStr string) (int, string, error) {
	parsedTime := time.Now().UTC()
	if dateStr != "" {
		parsedTime, _ = time.Parse(zuluDateOnly, dateStr)
	} else {
		dateStr = parsedTime.Format(zuluDateOnly)
	}
	dateInt, err := strconv.Atoi(parsedTime.Format(ymdForm))
	return dateInt, dateStr, err

}
