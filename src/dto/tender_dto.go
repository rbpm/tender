package dto

import (
	"tender/interfaces/data"
)

type TenderDTO struct {
	name string
	href string
	date string
	id   string
	isIT bool
}

func NewTenderDTO(name string, href string, date string, id string) *TenderDTO {
	p := TenderDTO{name: name, href: href, date: date, id: id, isIT: data.IsIt(name)}
	return &p
}

func (tender *TenderDTO) Name() string {
	return tender.name
}

func (tender *TenderDTO) Href() string {
	return tender.href
}

func (tender *TenderDTO) Date() string {
	return tender.date
}

func (tender *TenderDTO) Id() string {
	return tender.id
}

func (tender *TenderDTO) IsIT() bool {
	return tender.isIT
}
