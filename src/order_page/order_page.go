package order_page

import (
	"encoding/json"
	"fmt"
	"github.com/Noooste/azuretls-client"
	"tender/dto"
	"tender/interfaces/data"
)

func ProcessGetOrderPages(flags *dto.FlagDTO, err error, tenders []data.Data, done bool, tendersOldAll []data.Data) (error, []data.Data) {
	session := azuretls.NewSession()
	for page := 1; page <= flags.OrderPages; page++ {
		fmt.Println("order page: ", page)
		err, tenders, done = ProcessGetOrderPage(page, session, tenders, tendersOldAll)
		if done {
			fmt.Println("done")
			break
		}
	}
	return err, tenders
}

func ProcessGetOrderPage(page int, session *azuretls.Session, tenders []data.Data, tendersOldAll []data.Data) (error, []data.Data, bool) {
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
		return err, nil, false
	}
	for _, order := range orders {
		tender := order.GetTenderDTO()
		tenders = append(tenders, tender)
		if data.IsIn(tender, tendersOldAll) {
			fmt.Println("processGetOrderPage: old orders contains this:", tender)
			return err, tenders, true
		}
	}
	return err, tenders, false
}
