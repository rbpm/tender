package data

import (
	"strings"
)

type Data interface {
	Src() string
	Name() string
	Href() string
	Date() string
	Time() string
	Id() string
	IsIT() bool
}

// id can be different in the same name and date
func IsIn(tenders []Data, tender Data) bool {
	for _, p := range tenders {
		if p.Date() == tender.Date() && p.Name() == tender.Name() {
			// oneplace... server gives different hour:min
			if p.Time() != tender.Time() {
				println("TODO:", p.Time())
				println("TODO:", tender.Time())
			}
			return true
		}
	}
	return false
}

// "sztucznej inteligencji" - is too much and complicated
func IsIT(name string) bool {
	lowerName := strings.ToLower(name)
	return strings.Contains(lowerName, "oprogramowani") ||
		strings.Contains(lowerName, " it ") ||
		strings.Contains(lowerName, "rozwój i utrzymanie systemu") ||
		(strings.Contains(lowerName, "aplikacj") && !strings.Contains(lowerName, "folii") ||
			strings.Contains(lowerName, "software engineer") ||
			strings.Contains(lowerName, "programist"))
}
