package dto

import (
	"tender/interfaces/data"
)

type DataDTO struct {
	name string
	href string
	date string
	id   string
	isIT bool
}

func NewTenderDTO(name string, href string, date string, id string) *DataDTO {
	p := DataDTO{name: name, href: href, date: date, id: id, isIT: data.IsIT(name)}
	return &p
}

func (tender *DataDTO) Name() string {
	return tender.name
}

func (tender *DataDTO) Href() string {
	return tender.href
}

func (tender *DataDTO) Date() string {
	return tender.date
}

func (tender *DataDTO) Id() string {
	return tender.id
}

func (tender *DataDTO) IsIT() bool {
	return tender.isIT
}
