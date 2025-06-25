package dto

import "time"

type OrderDTO struct {
	ObjectId                string    `json:"objectId"`
	Title                   string    `json:"title"`
	OrganizationId          string    `json:"organizationId"`
	OrganizationName        string    `json:"organizationName"`
	OrganizationPartName    string    `json:"organizationPartName"`
	OrganizationCity        string    `json:"organizationCity"`
	OrganizationProvince    string    `json:"organizationProvince"`
	BzpNumber               string    `json:"bzpNumber"`
	TenderType              string    `json:"tenderType"`
	CompetitionType         string    `json:"competitionType"`
	ConcessionType          string    `json:"concessionType"`
	SubmissionOffersDate    time.Time `json:"submissionOffersDate"`
	TenderState             string    `json:"tenderState"`
	IsTenderAmountBelowEU   bool      `json:"isTenderAmountBelowEU"`
	TedContractNoticeNumber string    `json:"tedContractNoticeNumber"`
	InitiationDate          time.Time `json:"initiationDate"`
}

func (order OrderDTO) GetDataDTO() *DataDTO {
	href := "https://ezamowienia.gov.pl/mp-client/search/list/" + order.ObjectId
	return NewDataDTO(order.Title, href, order.SubmissionOffersDate.Format("2006.01.02"), order.ObjectId)
}
