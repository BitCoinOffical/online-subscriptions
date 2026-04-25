package date

import (
	"fmt"
	"time"
)

func ParseMonthDate(s string) (time.Time, error) {
	layouts := []string{"01-2006", "02-01-2006"}

	for _, l := range layouts {
		if t, err := time.Parse(l, s); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("invalid date format: %s", s)
}
