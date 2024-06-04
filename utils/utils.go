package utils

import (
	"fmt"
	"strconv"
	"time"

	"github.com/twinj/uuid"
)

const zuluDateOnly = "2006-01-02"
const ymdForm = "20060102"

// ParseZuluDate is used to parse a zulu formatted date to integer. If no date is entered
// ParseZuluDate returns a Zulu representation of current time
func ParseZuluDate(dateStr string) (int, string, error) {
	parsedTime := time.Now().UTC()
	var err error
	if dateStr != "" {
		parsedTime, err = time.Parse(zuluDateOnly, dateStr)
		if err != nil {
			return -1, dateStr, fmt.Errorf("date parameter value: %s is not in the valid form of YYYY-MM-DD", dateStr)
		}
	} else {
		dateStr = parsedTime.Format(zuluDateOnly)
	}
	dateInt, err := strconv.Atoi(parsedTime.Format(ymdForm))
	return dateInt, dateStr, err

}

func NewUUID() string {
	return uuid.NewV4().String()
}

func DistinctCast(distinctRes []interface{}) ([]string, error) {
	results := make([]string, 0, len(distinctRes))
	for _, result := range distinctRes {
		if value, ok := result.(string); ok {
			results = append(results, value)
		} else {
			return nil, fmt.Errorf("expected value of string type, instead got %T", result)
		}
	}
	return results, nil
}
