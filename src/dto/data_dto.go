package dto

import (
	"tender/interfaces/data"
)

type DataDTO struct {
	src  string
	name string
	href string
	date string
	id   string
	isIT bool
}

func NewDataDTO(src string, name string, href string, date string, id string) *DataDTO {
	p := DataDTO{src: src, name: name, href: href, date: date, id: id, isIT: data.IsIT(name)}
	return &p
}

func (tender *DataDTO) Src() string {
	return tender.src
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
