package tools

import (
	"errors"
	"fmt"
	"os"
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

func MkDirIfNotExist(path string) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			fmt.Println(err)
		}
	}
}
