package data

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type CustomTime struct {
	time.Time
}

func (ct *CustomTime) UnmarshalJSON(data []byte) error {
	str := strings.Trim(string(data), `"`)
	fmt.Printf("Trimmed string: '%s'\n", str)

	formats := []string{
		"2006-01-02T15:04:05.000000Z",
		"2006-01-02T15:04:05.000Z",
		"2006-01-02T15:04:05Z",
	}

	for i, format := range formats {
		fmt.Printf("Trying format %d: %s\n", i, format)
		if t, err := time.Parse(format, str); err == nil {
			ct.Time = t
			return nil
		} else {
			fmt.Printf("❌ Failed with format %s: %v\n", format, err)
		}
	}

	return fmt.Errorf("invalid timestamp format: %s", str)
}

func (ct *CustomTime) MarshalJSON() ([]byte, error) {
	formatted := ct.Time.Format("2006-01-02T15:04:05:000Z")
	return []byte(strconv.Quote(formatted)), nil
}
