package dto

import (
	"fmt"
	"tender/tools"
)

const BK_TIME_LAYOUT = "2006-01-02 15:04:05"

type BkDTO struct {
	Status string         `json:"status"`
	Data   BkDataTableDTO `json:"data"`
}

type BkDataTableDTO struct {
	Advertisements []BkDataDTO `json:"advertisements"`
	Meta           BkMetaDTO   `json:"meta"`
}

type BkMetaDTO struct {
	Total int `json:"total"`
}

type BkDataDTO struct {
	Id                 int    `json:"id"`
	Title              string `json:"title"`
	Content            string `json:"content"`
	AdvertiserName     string `json:"advertiser_name"`
	IsMine             bool   `json:"is_mine"`
	PublicationDate    string `json:"publication_date"`
	SubmissionDeadline string `json:"submission_deadline"`
	FulfillmentPlace   string `json:"fulfillment_place"`
	Favorite           string `json:"favorite"`
}

func (bk BkDataDTO) GetDataDTO() *DataDTO {
	href := "https://bazakonkurencyjnosci.funduszeeuropejskie.gov.pl/ogloszenia/" + fmt.Sprintf("%v", bk.Id)
	timePtr := tools.ParseDate(BK_TIME_LAYOUT, bk.SubmissionDeadline)
	return NewDataDTO("bazakonkurencyjnosci", bk.Title+"\n"+bk.Content, href, timePtr, fmt.Sprintf("%v", bk.Id))
}
