package data

import "strings"

type Data interface {
	Name() string
	Href() string
	Date() string
	Id() string
	IsIT() bool
}

func IsIn(tender Data, tenders []Data) bool {
	for _, p := range tenders {
		if p.Id() == tender.Id() && p.Name() == tender.Name() {
			return true
		}
	}
	return false
}

func IsIt(name string) bool {
	lowerName := strings.ToLower(name)
	return strings.Contains(lowerName, "oprogramowani") ||
		strings.Contains(lowerName, " it ") ||
		strings.Contains(lowerName, "rozwój i utrzymanie systemu") ||
		strings.Contains(lowerName, "aplikacj")
}
