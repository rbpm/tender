package dto

import "strings"

type TenderDTO struct {
	Name string
	Href string
	Date string
	Id   string
	IsIT bool
}

func NewTenderDTO(name string, href string, date string, id string) *TenderDTO {
	p := TenderDTO{Name: name, Href: href, Date: date, Id: id, IsIT: IsIt(name)}
	return &p
}

func IsIt(name string) bool {
	lowerName := strings.ToLower(name)
	return strings.Contains(lowerName, "oprogramowani") ||
		strings.Contains(lowerName, " it ") ||
		strings.Contains(lowerName, "rozwój i utrzymanie systemu") ||
		strings.Contains(lowerName, "aplikacj")
}

func (tender *TenderDTO) IsIn(tenders []*TenderDTO) bool {
	for _, p := range tenders {
		if p.Id == tender.Id && p.Name == tender.Name {
			return true
		}
	}
	return false
}
