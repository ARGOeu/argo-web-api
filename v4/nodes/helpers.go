package nodes

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// parseZuluDate is used to parse a zulu formatted date to integer
func parseZuluDate(dateStr string) (int, error) {
	parsedTime, err := time.Parse(zuluForm, dateStr)
	if err != nil {
		return -1, err
	}
	return strconv.Atoi(parsedTime.Format(ymdForm))
}

// periodToDays parses '7d' or '1w' strings and returns total days
func periodToDays(period string) (int, error) {
	period = strings.TrimSpace(period)
	if len(period) < 2 {
		return 0, fmt.Errorf("invalid period: %q", period)
	}

	value, err := strconv.Atoi(period[:len(period)-1])
	if err != nil {
		return 0, fmt.Errorf("wrong period defined %q", period)
	}
	if value < 0 {
		return 0, fmt.Errorf("don't specify negative period %q", period)
	}
	unit := period[len(period)-1]

	switch unit {
	case 'd':
		return value, nil
	case 'w':
		return value * 7, nil
	default:
		return 0, fmt.Errorf("unsupported unit: %q", unit)
	}
}

func dateInt(t time.Time) int {
	return t.Year()*10000 + int(t.Month())*100 + t.Day()
}
