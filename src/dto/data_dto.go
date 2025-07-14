package dto

import (
	"tender/interfaces/data"
	"time"
)

type DataDTO struct {
	src  string
	name string
	href string
	time *time.Time
	id   string
	isIT bool
}

func NewDataDTO(src string, name string, href string, time *time.Time, id string) *DataDTO {
	p := DataDTO{src: src, name: name, href: href, time: time, id: id, isIT: data.IsIT(name)}
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
	if tender.time == nil {
		return ""
	}
	return tender.time.Format("2006-01-02")
}

func (tender *DataDTO) Id() string {
	return tender.id
}

func (tender *DataDTO) IsIT() bool {
	return tender.isIT
}
