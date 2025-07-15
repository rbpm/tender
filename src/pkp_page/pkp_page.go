package pkp_page

import (
	"fmt"
	"net/http"
	"tender/dto"
	"tender/interfaces/data"
	"tender/tools"

	"github.com/Noooste/azuretls-client"
	"github.com/gurkankaymak/gosoup"
)

func ProcessGetPkpPages(err error, client string, tenders []data.Data, done bool, tendersOldAll []data.Data) (error, []data.Data) {
	session := azuretls.NewSession()
	for page := 1; page <= 1; page++ {
		fmt.Println("pkp page: ", page)
		err, tenders, done = ProcessGetPkpPage(page, client, session, tenders, tendersOldAll)
		if done {
			fmt.Println("done")
			break
		}
	}
	return err, tenders
}

func ProcessGetPkpPage(page int, client string, session *azuretls.Session, tenders []data.Data, tendersOldAll []data.Data) (error, []data.Data, bool) {
	// 99999999 is set on page when selected "all"
	// other 100
	pageSize := 99999999
	pageStart := page*pageSize - 1
	requestData := `{
		"CSRFToken": "",
		"UNIQUE_ID": "1752382931294Tv0IwqqxfmcxMD6dBHT",
		"GD_sortOrder": "desc",
		"GD_sortCol": "5",
		"isFiredBySearchButton": "false",
		"demandName": "",
		"demandNumber": "",
		"offerDeadlineFrom": "",
		"offerDeadlineTo": "",
		"procedureStartDateFrom": "",
		"procedureStartDateTo": "",
		"demandOrganization_orgName": "",
		"demandOrganization_withSuborgs": "true",
		"demandDealingCategoryItemCpvHierarchicalId": "",
		"demandCategoryItemHierarchicalId": "",
		"ta_jsp.searchform.header.demand.notice": "0",
		"clearGridSearchTagSearchFormScript": "",
		"GD_pagesize": "%d",
		"GD_pagestart": "%d"
    }`
	requestData = fmt.Sprintf(requestData, pageSize, pageStart)

	req := &azuretls.Request{
		Method: http.MethodPost,
		Url:    "https://platformazakupowa.plk-sa.pl/app/demand/notice/public/current/list?USER_MENU_HOVER=currentNoticeList",
		Body:   requestData,
	}
	response, err := session.Do(req)

	element, err := gosoup.ParseAsHTML(response.String())
	tableElements := element.FindAll("table", gosoup.Attributes{"id": "publicList"})

	if len(tableElements) < 1 {
		fmt.Println("could not find table element")
		return err, tenders, true
	}
	tableElement := tableElements[0]

	tableBodyElement := tableElement.FindByTag("tbody")
	if tableBodyElement == nil {
		fmt.Println("could not find table body element")
		return err, tenders, true
	}
	expectedTag := "tr"
	expectedTdElementsSize := 12
	elements := tableBodyElement.FindAllByTag(expectedTag)
	for _, element := range elements {
		if element.Data != expectedTag {
			fmt.Printf("wrong element tag, expected: %s, actual: %s", expectedTag, element.Data)
		}
		id, _ := element.GetAttribute("id")
		tdElements := element.FindAllByTag("td")
		if tdElements == nil {
			//fmt.Println("table headers:\n", element)
		} else if len(tdElements) != expectedTdElementsSize {
			fmt.Println("wrong number of td elements", len(tdElements))
		} else {
			// <th>...</th>
			// 0 Numer postępowania
			// 1 Nazwa postępowania
			// 2 Podstawa prawna
			// 3 Tryb postępowania
			// 4 Rodzaj zamówienia
			// 5 Kategoria zakupowa / grupa materiałowa
			// 6 Data publikacji
			// 7 Termin składania<
			// 8 Osoba kontaktowa<
			// 9 Jednostka organizacyjna
			// 10 Kod CPV
			// 11 Data wszczęcia
			// idPostepowania := tdElements[0].FirstChild.Data
			title := tdElements[1].FirstChild.Data
			// publishedDate := tdElements[6].FirstChild.Data
			offerDateTime := tdElements[7].FirstChild.Data
			const longForm = "2006-01-02 15:04"
			timePtr := tools.ParseDate(longForm, offerDateTime)
			// createdDate := tdElements[11].FirstChild.Data
			href := "https://platformazakupowa.plk-sa.pl/app/demand/notice/public/" + id + "/details"
			tender := dto.NewDataDTO(client, title, href, timePtr, id)

			tenders = append(tenders, tender)
			if data.IsIn(tendersOldAll, tender) {
				fmt.Println("processGetPzPage: old pkp contains this", tender)
				return err, tenders, true
			}
		}
	}

	return err, tenders, false
}
