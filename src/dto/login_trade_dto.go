package dto

import (
	"time"
)

type LoginTradeDTO struct {
	Id        string
	TitleName string
	TitleID   string
	Href      string
	Date      string
	StartDate string
	EndDate   string
	ClientID  string
	Status    string
}

func NewLoginTradeDTO(url string, id string, titleName string, titleID string, href string, date string, startDate string, endDate string, frogID string, status string) *LoginTradeDTO {
	p := LoginTradeDTO{Id: id, TitleName: titleName, TitleID: titleID, Href: url + href, Date: date, StartDate: startDate, EndDate: endDate, ClientID: frogID, Status: status}
	return &p
}

func (loginTradeDto LoginTradeDTO) GetDataDTO(client string) *DataDTO {
	const longForm = "2006-01-02 15:04"
	dateTime, _ := time.Parse(longForm, loginTradeDto.EndDate)
	dateValue := dateTime.Format("2006-01-02")
	// loginTradeDto.ClientID?
	return NewDataDTO(client, loginTradeDto.TitleName+"\n"+loginTradeDto.TitleID, loginTradeDto.Href, dateValue, loginTradeDto.Id)
}
