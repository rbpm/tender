package bk_page

import (
	"encoding/json"
	"fmt"
	"tender/dto"
	"tender/interfaces/data"

	"github.com/Noooste/azuretls-client"
)

func ProcessGetBkPages(flags *dto.FlagDTO, err error, tenders []data.Data, done bool, tendersOldAll []data.Data) (error, []data.Data) {
	session := azuretls.NewSession()
	if flags.BkPages > 0 {
		for page := 1; page <= flags.BkPages; page++ {
			fmt.Println("bk page: ", page)
			err, tenders, done = processGetBkPage(page, session, tenders, tendersOldAll)
			if done {
				fmt.Println("done")
				break
			}
		}
	} else {
		total := 0
		err, total = processGetBkTotal(session)
		err, tenders, done = processGetBkPage(total, session, tenders, tendersOldAll)
		if done {
			fmt.Println("bk done")
		}
	}
	return err, tenders
}

func processGetBkTotal(session *azuretls.Session) (error, int) {
	var bks dto.BkDTO
	response, err := session.Get("https://bazakonkurencyjnosci.funduszeeuropejskie.gov.pl/api/announcements/search?page=1&limit=20&sort=publicationDate&status%5B0%5D=PUBLISHED")
	if err != nil {
		return err, 0
	}
	err = json.Unmarshal([]byte(response.String()), &bks)
	if err != nil {
		println(err.Error())
		println("response:" + response.String())
		return err, 0
	}
	total := bks.Data.Meta.Total
	fmt.Println("BK total:", total)
	return err, total
}

func processGetBkPage(limit int, session *azuretls.Session, tenders []data.Data, tendersOldAll []data.Data) (error, []data.Data, bool) {
	var bks dto.BkDTO
	limitStr := fmt.Sprintf("%d", limit)
	response, err := session.Get("https://bazakonkurencyjnosci.funduszeeuropejskie.gov.pl/api/announcements/search?page=1&limit=" + limitStr + "&sort=publicationDate&status%5B0%5D=PUBLISHED")
	if err != nil {
		println(err.Error())
		return err, tenders, true
	}
	err = json.Unmarshal([]byte(response.String()), &bks)
	if err != nil {
		println(err.Error())
		println("response:" + response.String())
		return err, tenders, true
	}
	total := bks.Data.Meta.Total
	fmt.Println("BK total:", total, " for limit:", limit)
	for _, bk := range bks.Data.Advertisements {
		tender := bk.GetDataDTO()
		tenders = append(tenders, tender)
		if data.IsIn(tendersOldAll, tender) {
			fmt.Println("processGetBkPage: old bk contains this:", tender)
			return err, tenders, true
		}
	}
	return err, tenders, false
}
