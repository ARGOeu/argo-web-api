package nodes

import (
	"strconv"
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
