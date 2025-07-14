package dto

import (
	"fmt"
	"tender/tools"
)

const PKO_TIME_LAYOUT = "02.01.2006 15:04"

type PkoDTO struct {
	Count    int            `json:"count"`
	Next     string         `json:"next"`
	Previous string         `json:"previous"`
	Results  []PkoResultDTO `json:"results"`
}

type PkoResultDTO struct {
	Path       string        `json:"path"`
	NewsListId int           `json:"newslist_id"`
	Snippet    PkoSnippetDTO `json:"snippet"`
	Id         int           `json:"id"`
	ParentId   int           `json:"parent_id"`
	Filters    PkoFiltersDTO `json:"filters"`
}

type PkoSnippetDTO struct {
	Title               PkoTextDTO `json:"title"`
	Lead                string     `json:"lead"`
	Label               string     `json:"label"`
	LabelColor          string     `json:"label_color"`
	RawPublicationDate  string     `json:"raw_publication_date"`
	PublicationDate     string     `json:"publication_date"`
	ShowPublicationTime bool       `json:"show_publication_time"`
	Featured            bool       `json:"featured"`
	FileName            string     `json:"file_name"`
}

type PkoTextDTO struct {
	Text string `json:"text"`
}

type PkoFiltersDTO struct {
	Categories []int `json:"categories"`
	Customers  []int `json:"customers"`
	Years      []int `json:"years"`
}

func (pkoResultDto PkoResultDTO) GetDataDTO() *DataDTO {
	href := "https://www.pkobp.pl" + pkoResultDto.Path
	//"Termin nadsyłania ofert upływa w dniu 02.01.2006 roku, o godzinie 15:04."
	str := pkoResultDto.Snippet.Lead[40 : len(pkoResultDto.Snippet.Lead)-1]
	str = str[:10] + str[27:]
	timePtr := tools.ParseDate(PKO_TIME_LAYOUT, str)
	return NewDataDTO("pko", pkoResultDto.Snippet.Title.Text, href, timePtr, fmt.Sprintf("%v", pkoResultDto.Id))
}
