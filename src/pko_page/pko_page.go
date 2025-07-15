package pko_page

import (
	"encoding/json"
	"fmt"
	"tender/dto"
	"tender/interfaces/data"

	"github.com/Noooste/azuretls-client"
)

func ProcessGetPkoPages(flags *dto.FlagDTO, err error, tenders []data.Data, done bool, tendersOldAll []data.Data) (error, []data.Data) {
	session := azuretls.NewSession()
	for page := 1; page <= flags.PkoPages; page++ {
		fmt.Println("pko page: ", page)
		err, tenders, done = ProcessGetPkoPage(page, session, tenders, tendersOldAll)
		if done {
			fmt.Println("done")
			break
		}
	}
	return err, tenders
}

func ProcessGetPkoPage(page int, session *azuretls.Session, tenders []data.Data, tendersOldAll []data.Data) (error, []data.Data, bool) {
	var pkoDto dto.PkoDTO
	pageStr := fmt.Sprintf("%d", page)
	response, err := session.Get("https://www.pkobp.pl/api/news/items?page=" + pageStr + "&page_size=8&page_id=649&categories=8&variant=contents")
	if err != nil {
		println(err.Error())
		return err, tenders, true
	}
	err = json.Unmarshal([]byte(response.String()), &pkoDto)
	if err != nil {
		println(err.Error())
		println("DONE with error response:" + response.String())
		return err, tenders, true
	}
	println(pkoDto.Count)
	for _, pko := range pkoDto.Results {
		if len(pko.Snippet.Lead) > 0 {
			tender := pko.GetDataDTO()
			tenders = append(tenders, tender)
			if data.IsIn(tendersOldAll, tender) {
				fmt.Println("processGetPkoPage: old pko contains this:", tender)
				return err, tenders, true
			}
		}
	}
	return err, tenders, false
}
