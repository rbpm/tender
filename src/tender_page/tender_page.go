package tender_page

import (
	"fmt"
	"github.com/Noooste/azuretls-client"
	"github.com/gurkankaymak/gosoup"
	"strings"
	"tender/dto"
	"time"
)

func ProcessGetTenderPage(page int, session *azuretls.Session, tendersIT []*dto.TenderDTO, tenders []*dto.TenderDTO, tendersOldAll []*dto.TenderDTO) (error, []*dto.TenderDTO, []*dto.TenderDTO, bool) {
	pageStr := fmt.Sprintf("%d", page)

	//https://oneplace.marketplanet.pl/zapytania-ofertowe-przetargi/-/rfp/cat?_7_WAR_organizationnoticeportlet_order=createDateDesc
	//why not this?        https://oneplace.marketplanet.pl/zapytania-ofertowe-przetargi/-/rfp/cat?_7_WAR_organizationnoticeportlet_cur=1_order=createDateDesc
	//only this form is on www for page n: https://oneplace.marketplanet.pl/zapytania-ofertowe-przetargi/-/rfp/cat?_7_WAR_organizationnoticeportlet_cur=1
	response, err := session.Get("https://oneplace.marketplanet.pl/zapytania-ofertowe-przetargi/-/rfp/cat?_7_WAR_organizationnoticeportlet_cur=" + pageStr)
	if err != nil {
		panic(err)
	}

	element, err := gosoup.ParseAsHTML(response.String())
	if err != nil {
		fmt.Println("could not parse")
		return err, tendersIT, tenders, false
	}
	containerElement := element.Find("div", gosoup.Attributes{"id": "_7_WAR_organizationnoticeportlet_selectNoticesSearchContainer"})
	if containerElement == nil {
		fmt.Println("could not find container element")
	}
	subContainer := containerElement.Find("div", gosoup.Attributes{"class": "lfr-search-container-list"})
	if subContainer == nil {
		fmt.Println("could not find subContainer element")
	}
	group := subContainer.FindByTag("dl")
	if group == nil {
		fmt.Println("could not find group element")
	}
	expectedTag := "dd"
	expectedAttrKey := "data-qa-id"
	expectedAttrVal := "row"
	expectedElementsSize := 12
	elements := group.FindAll(expectedTag, gosoup.Attributes{expectedAttrKey: expectedAttrVal})
	if len(elements) != expectedElementsSize {
		fmt.Printf("wrong number of elements found: %q, expected number: %q", len(elements), expectedElementsSize)
	}
	for _, element := range elements {
		if element.Data != expectedTag {
			fmt.Printf("wrong element tag, expected: %q, actual: %q", expectedTag, element.Data)
		}
		attributeValue, ok := element.GetAttribute(expectedAttrKey)
		if !ok || attributeValue != expectedAttrVal {
			fmt.Printf("expected attribute: %q: %q does not exist", expectedAttrKey, expectedAttrVal)
		}
		//TODO add app parameter --debug
		if false {
			fmt.Println("dd element:", element)
		}
		aTag := element.FindByTag("a")
		if aTag == nil {
			fmt.Println("could not find aTag element")
		}
		hrefValue, ok := aTag.GetAttribute("href")
		if !ok {
			fmt.Printf("href attribute: does not exist")
		}

		nameValue := aTag.FirstChild.Data
		//class="notice-date"
		//dateDiv := element.Find("div", gosoup.Attributes{"class": "notice-date"})
		dateSpan := element.Find("span", gosoup.Attributes{"title": "Termin składania ofert"})
		dateTimeValue := strings.TrimSpace(dateSpan.FirstChild.Data)
		//t, err := time.Parse(time.RFC3339, "2023-05-02T09:34:01Z")
		//Mon Jun 23 09:00:00 GMT 2025: example value
		//Mon Jan _2 15:04:05 GMT 2006: layout form
		const longForm = "Mon Jan _2 15:04:05 GMT 2006"
		dateTime, _ := time.Parse(longForm, dateTimeValue)
		dateValue := dateTime.Format("2006.01.02")
		tender := dto.NewTenderDTO(nameValue, hrefValue, dateValue, getHrefID(hrefValue))
		tendersIT, tenders = tender.AppendTo(tendersIT, tenders)
		if tender.IsIn(tendersOldAll) {
			fmt.Println("processGetTenderPage: old tenders contains this", tender)
			return err, tendersIT, tenders, true
		}
	}
	return err, tendersIT, tenders, false
}

func getHrefID(value string) string {
	if len(value) < 10 {
		return "len err"
	}
	pos := strings.Index(value, "_noticeId=")
	if pos == -1 {
		return "index err"
	}
	id := value[pos+10:]
	return id
}
