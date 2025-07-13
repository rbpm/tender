package kghm_page

import (
	"encoding/json"
	"fmt"
	"github.com/Noooste/azuretls-client"
	"github.com/gurkankaymak/gosoup"
	"tender/dto"
	"tender/interfaces/data"
)

type GetHrefID func(string) string

func GetDefaultHrefID(value string) string {
	if len(value) < 39 {
		return "len err"
	}
	id := value[34 : len(value)-5]
	return id
}

func ProcessGetKghmPages(pages int, err error, tenders []data.Data, done bool, tendersOldAll []data.Data) (error, []data.Data) {
	session := azuretls.NewSession()
	for page := 1; page <= pages; page++ {
		fmt.Println("kghm page: ", page)
		err, tenders, done = ProcessGetKghmPage(page, session, tenders, tendersOldAll)
		if done {
			fmt.Println("done")
			break
		}
	}
	return err, tenders
}

func ProcessGetKghmPage(page int, session *azuretls.Session, tenders []data.Data, tendersOldAll []data.Data) (error, []data.Data, bool) {
	pageStr := fmt.Sprintf("%d", page)

	urlPrefix := "https://kghm.com/pl/views/ajax?_wrapper_format=drupal_ajax&view_name=tenders_normal&view_display_id=block_tenders_future&view_args=&view_path=/node/9217&view_base_path=&view_dom_id=0534474bbb5e8add417d226807501a98bef0682dc080a6e60d68f3c2ba52e181&pager_element=0&page="
	urlSuffix := "&_drupal_ajax=1&ajax_page_state[theme]=kghm&ajax_page_state[theme_token]=&ajax_page_state[libraries]=eJx1kOFugzAMhF8ImkeKTHKFFCdGsdnKnn4B2mqV1j_OdxfrZJtiNKGyOXrA5VqlWDfADNXjvogi-mviJtVFMiwpzI0_tYwoqMRdYFLd3ECKJ2eo0gh96iKx_UlFo5qJ0w86rD6IzAntyQsnKgHuP9NHXGll6-Zxym4kZtSt51RmfVgsA3F_e5dqW-sZT28S67PEldFXLEwBGeWRuLRJD_I7eRNhdUftKViScuZ6NQlzu0KYqIxwUrAPgS7LkBi-Ba4u7_XdObjf-RWnoBomT0vytLbUfVUY3Ae_000N-bywDS13pCPbprbG5Y_TfSV8qzvqhW50fzPOA_wCIY_LtQ"
	ajaxUrl := urlPrefix + pageStr + urlSuffix

	var kghmAjaxDto []dto.KghmAjaxDTO
	response, err := session.Get(ajaxUrl)
	if err != nil {
		fmt.Println("AJAX ERROR:", err.Error())
		return err, tenders, true
	}
	err = json.Unmarshal([]byte(response.String()), &kghmAjaxDto)
	if err != nil {
		println(err.Error())
		println("AJAX Unmarshal ERROR:" + response.String())
		return err, tenders, true
	}
	element, err := gosoup.ParseAsHTML(kghmAjaxDto[2].Data)
	if err != nil {
		fmt.Println("could not parse")
		return err, tenders, false
	}
	containerElement := element.FindAll("div", gosoup.Attributes{"class": "view-content"})
	if containerElement == nil {
		fmt.Println("could not find container element")
	}
	var tableElement *gosoup.Element
	if containerElement == nil || len(containerElement) == 0 {
		fmt.Println("could not find container \"view-content\" element")
		return err, tenders, true
	} else if len(containerElement) == 1 {
		// AJAX
		tableElement = containerElement[0].FindByTag("table")
	} else {
		// GET https://kghm.com/pl/przetargi-nieograniczone
		tableElement = containerElement[1].FindByTag("table")
	}

	if tableElement == nil {
		fmt.Println("could not find table element")
	}
	tableBodyElement := tableElement.FindByTag("tbody")
	if tableBodyElement == nil {
		fmt.Println("could not find table body element")
	}
	expectedTag := "tr"
	expectedTdElementsSize := 5
	elements := tableBodyElement.FindAllByTag(expectedTag)
	for _, element := range elements {
		if element.Data != expectedTag {
			fmt.Printf("wrong element tag, expected: %s, actual: %s", expectedTag, element.Data)
		}
		tdElements := element.FindAllByTag("td")
		if tdElements == nil {
			//fmt.Println("table header:\n", element)
		} else {
			if len(tdElements) != expectedTdElementsSize {
				fmt.Println("wrong number of td elements", len(tdElements))
			} else {
				startTimeTdElement := tdElements[0]
				startTimeElement := startTimeTdElement.FindByTag("time")
				publiahedDateValue, ok := startTimeElement.GetAttribute("datetime")
				if !ok {
					fmt.Printf("[0]datetime attribute: does not exist")
				}

				endTimeTdElement := tdElements[1]
				endTimeElement := endTimeTdElement.FindByTag("time")
				offerDateValue, ok := endTimeElement.GetAttribute("datetime")
				if !ok {
					fmt.Printf("[1]datetime attribute: does not exist")
				}

				titleTdElement := tdElements[3]
				titleElement := titleTdElement.FindByTag("a")

				titleValue := ""
				name := titleElement.FirstChild
				if name != nil {
					titleValue = name.Data
				}

				hrefValue, ok := titleElement.GetAttribute("href")
				if !ok {
					fmt.Printf("href attribute: does not exist")
				}

				updateTimeTdElement := tdElements[4]
				updateTimeElement := updateTimeTdElement.FindByTag("time")
				updatedDateValue, ok := updateTimeElement.GetAttribute("datetime")
				if !ok {
					fmt.Printf("[0]datetime attribute: does not exist")
				}

				_ = dto.NewKghmDTO("kghm", titleValue, "https://kghm.com/"+hrefValue, offerDateValue, publiahedDateValue, updatedDateValue, hrefValue)
				tender := dto.NewDataDTO("kghm", titleValue, "https://kghm.com/"+hrefValue, offerDateValue, hrefValue)

				tenders = append(tenders, tender)
				if data.IsIn(tendersOldAll, tender) {
					fmt.Println("processGetTenderPage: old tenders contains this", tender)
					return err, tenders, true
				}
			}
		}
	}
	return err, tenders, false
}
