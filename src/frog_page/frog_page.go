package frog_page

import (
	"fmt"
	"github.com/Noooste/azuretls-client"
	"github.com/gurkankaymak/gosoup"
	"tender/dto"
	"tender/interfaces/data"
)

func ProcessGetFrogPages(flags *dto.FlagDTO, err error, tenders []data.Data, done bool, tendersOldAll []data.Data) (error, []data.Data) {
	session := azuretls.NewSession()
	for page := 1; page <= flags.FrogPages; page++ {
		fmt.Println("frog page: ", page)
		err, tenders, done = ProcessGetFrogPage(page, session, tenders, tendersOldAll)
		if done {
			fmt.Println("done")
			break
		}
	}
	return err, tenders
}

func ProcessGetFrogPage(page int, session *azuretls.Session, tenders []data.Data, tendersOldAll []data.Data) (error, []data.Data, bool) {
	pageStr := fmt.Sprintf("%d", page)

	response, err := session.Get("https://zabka.logintrade.net/portal,listaZapytaniaOfertowe.html?page=" + pageStr + "&itemsperpage=100&search_sort=9")
	if err != nil {
		panic(err)
	}

	element, err := gosoup.ParseAsHTML(response.String())
	if err != nil {
		fmt.Println("could not parse")
		return err, tenders, false
	}
	containerElement := element.Find("div", gosoup.Attributes{"class": "dataTableContent"})
	if containerElement == nil {
		fmt.Println("could not find container element")
	}
	tableElement := containerElement.FindByTag("table")
	if tableElement == nil {
		fmt.Println("could not find table element")
	}
	tableBodyElement := containerElement.FindByTag("tbody")
	if tableBodyElement == nil {
		fmt.Println("could not find table body element")
	}
	expectedTag := "tr"
	expectedElementsSize := 101
	expectedTdElementsSize := 6
	expectedAElementsSize := 2
	elements := tableBodyElement.FindAllByTag(expectedTag)
	if len(elements) != expectedElementsSize {
		fmt.Printf("wrong number of elements found: %d, expected number: %d\n", len(elements), expectedElementsSize)
	}
	for _, element := range elements {
		if element.Data != expectedTag {
			fmt.Printf("wrong element tag, expected: %s, actual: %s", expectedTag, element.Data)
		}

		tdElements := element.FindAllByTag("td")

		if tdElements == nil {
			fmt.Println("table header:\n", element)
		} else {
			if len(tdElements) != expectedTdElementsSize {
				fmt.Println("wrong number of td elements", len(tdElements))
			} else {
				titleTdElement := tdElements[0]
				aElements := titleTdElement.FindAllByTag("a")
				if len(aElements) != expectedAElementsSize {
					fmt.Println("wrong number of <a> elements:", len(aElements))
				} else {
					nameElement := aElements[0]
					nameValue := nameElement.FirstChild.Data

					numberElement := aElements[1]
					numberValue := numberElement.FirstChild.Data

					hrefValue, ok := nameElement.GetAttribute("href")
					if !ok {
						fmt.Printf("href attribute: does not exist")
					}

					dateElement := tdElements[1]
					dateValue := dateElement.FirstChild.Data

					startElement := tdElements[2]
					startValue := startElement.FirstChild.Data

					endElement := tdElements[3]
					endValue := endElement.FirstChild.Data

					frogElement := tdElements[4]
					frogValue := frogElement.FirstChild.Data

					statusElement := tdElements[5]
					spanElement := statusElement.FindByTag("span")
					statusValue, ok := spanElement.GetAttribute("title")
					if !ok {
						fmt.Printf("span attribute: does not exist")
					}
					//NewFrogDTO(id , titleName , titleID , href , date , startDate , endDate , frogID , status)
					frog := dto.NewFrogDTO(getHrefID(hrefValue), nameValue, numberValue, hrefValue, dateValue, startValue, endValue, frogValue, statusValue)
					fmt.Println(frog)
					tender := frog.GetDataDTO()
					tenders = append(tenders, tender)
					if data.IsIn(tendersOldAll, tender) {
						fmt.Println("processGetTenderPage: old tenders contains this", tender)
						return err, tenders, true
					}
				}
			}
		}
	}
	return err, tenders, false
}

func getHrefID(value string) string {
	if len(value) < 39 {
		return "len err"
	}
	id := value[34 : len(value)-5]
	return id
}
