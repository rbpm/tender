package order_page

import (
	"encoding/json"
	"fmt"
	"github.com/Noooste/azuretls-client"
	"tender/dto"
)

func ProcessGetOrderPage(page int, session *azuretls.Session, tendersIT []*dto.TenderDTO, tenders []*dto.TenderDTO, tendersOldAll []*dto.TenderDTO) (error, []*dto.TenderDTO, []*dto.TenderDTO, bool) {
	var orders []dto.OrderDTO
	pageStr := fmt.Sprintf("%d", page)
	response, err := session.Get("https://ezamowienia.gov.pl/mp-readmodels/api/Search/SearchTenders?SortingColumnName=InitiationDate&SortingDirection=DESC&PageNumber=" + pageStr + "&PageSize=50")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal([]byte(response.String()), &orders)
	if err != nil {
		println(err.Error())
		println("response:" + response.String())
		return err, nil, nil, false
	}
	for _, order := range orders {
		tender := order.GetTenderDTO()
		tendersIT, tenders = tender.AppendTo(tendersIT, tenders)
		if tender.IsIn(tendersOldAll) {
			fmt.Println("processGetOrderPage: old orders contains this:", tender)
			return err, tendersIT, tenders, true
		}
	}
	return err, tendersIT, tenders, false
}
