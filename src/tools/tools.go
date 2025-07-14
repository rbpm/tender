package tools

import (
	"fmt"
	"time"
)

func ParseDate(layout, value string) *time.Time {
	if timeValue, err := time.Parse(layout, value); err == nil {
		return &timeValue
	} else {
		fmt.Println("tools.ParseDate() error:", err.Error())
		return nil
	}
}
