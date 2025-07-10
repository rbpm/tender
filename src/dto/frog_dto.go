package dto

import (
	"time"
)

type FrogDTO struct {
	Id        string
	TitleName string
	TitleID   string
	Href      string
	Date      string
	StartDate string
	EndDate   string
	FrogID    string
	Status    string
}

func NewFrogDTO(id string, titleName string, titleID string, href string, date string, startDate string, endDate string, frogID string, status string) *FrogDTO {
	p := FrogDTO{Id: id, TitleName: titleName, TitleID: titleID, Href: "https://zabka.logintrade.net/" + href, Date: date, StartDate: startDate, EndDate: endDate, FrogID: frogID, Status: status}
	return &p
}

func (frogDto FrogDTO) GetDataDTO() *DataDTO {
	const longForm = "2006-01-02 15:04"
	dateTime, _ := time.Parse(longForm, frogDto.EndDate)
	dateValue := dateTime.Format("2006-01-02")
	// frogDto.FrogID?
	return NewDataDTO("zabka", frogDto.TitleName+"\n"+frogDto.TitleID, frogDto.Href, dateValue, frogDto.Id)
}
