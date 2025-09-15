package ezamowienia_page

import (
	"encoding/json"
	"fmt"
	"tender/dto"
	"tender/interfaces/data"
	"time"

	"github.com/Noooste/azuretls-client"
)

func ProcessGetEzamowieniaPages(flags *dto.FlagDTO, err error, tenders []data.Data, done bool, tendersOldAll []data.Data) (error, []data.Data) {
	session := azuretls.NewSession()
	session.SetTimeout(5 * time.Minute)
	for page := 1; page <= flags.EzamowieniaPages; page++ {
		fmt.Println("order page: ", page)
		err, tenders, done = ProcessGetEzamowieniaPage(page, session, tenders, tendersOldAll)
		if done {
			fmt.Println("done")
			break
		}
	}
	return err, tenders
}

func ProcessGetEzamowieniaPage(page int, session *azuretls.Session, tenders []data.Data, tendersOldAll []data.Data) (error, []data.Data, bool) {
	var orders []dto.EzamowieniaDTO
	pageStr := fmt.Sprintf("%d", page)
	response, err := session.Get("https://ezamowienia.gov.pl/mp-readmodels/api/Search/SearchTenders?SortingColumnName=InitiationDate&SortingDirection=DESC&PageNumber=" + pageStr + "&PageSize=50")
	if err != nil {
		println(err.Error())
		return err, tenders, true
	}
	err = json.Unmarshal([]byte(response.String()), &orders)
	if err != nil {
		println(err.Error())
		println("response:" + response.String())
		return err, nil, false
	}
	for _, order := range orders {
		tender := order.GetDataDTO()
		tenders = append(tenders, tender)
		if data.IsIn(tendersOldAll, tender) {
			fmt.Println("processGetOrderPage: old ezamowienia contains this:", tender)
			return err, tenders, true
		}
	}
	return err, tenders, false
}
