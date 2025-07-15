package login_trade_page

import (
	"fmt"
	"tender/dto"
	"tender/interfaces/data"

	"github.com/Noooste/azuretls-client"
	"github.com/gurkankaymak/gosoup"
)

const DEFAULT_URL_PREFIX = "portal,listaZapytaniaOfertowe.html?status_realizacji_zapytania[]=oczekiwanie_ofert&wojewodztwo=wszystkie&search=&search_sort=9&page="
const DEFAULT_URL_SUFIX = "&itemsperpage=100"

type GetHrefID func(string) string

func GetDefaultHrefID(value string) string {
	if len(value) < 39 {
		return "len err"
	}
	id := value[34 : len(value)-5]
	return id
}

func ProcessGetLoginTradePages(client string, url string, getHrefID GetHrefID, urlPrefix string, urlSuffix string, pages int, err error, tenders []data.Data, done bool, tendersOldAll []data.Data) (error, []data.Data) {
	session := azuretls.NewSession()
	for page := 1; page <= pages; page++ {
		fmt.Println("loginTrade page: ", page)
		err, tenders, done = ProcessGetLoginTradePage(client, url, getHrefID, urlPrefix, urlSuffix, page, pages, session, tenders, tendersOldAll)
		if done {
			fmt.Println("done")
			break
		}
	}
	return err, tenders
}

func ProcessGetLoginTradePage(client string, url string, getHrefID GetHrefID, urlPrefix string, urlSuffix string, page int, pages int, session *azuretls.Session, tenders []data.Data, tendersOldAll []data.Data) (error, []data.Data, bool) {
	pageStr := fmt.Sprintf("%d", page)

	getUrl := urlPrefix + pageStr + urlSuffix
	fmt.Println(getUrl)
	response, err := session.Get(getUrl)
	if err != nil {
		println(err.Error())
		return err, tenders, true
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
	expectedElementsSize := pages + 1
	expectedTdElementsSize := 6
	expectedAElementsSize := 2
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
				titleTdElement := tdElements[0]
				aElements := titleTdElement.FindAllByTag("a")
				if len(aElements) != expectedAElementsSize {
					fmt.Println("wrong number of <a> elements:", len(aElements))
				} else {
					nameElement := aElements[0]
					nameValue := ""
					name := nameElement.FirstChild
					if name != nil {
						nameValue = name.Data
					}

					numberElement := aElements[1]
					number := numberElement.FirstChild
					numberValue := ""
					if number != nil {
						numberValue = number.Data
					}

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

					ClientElement := tdElements[4]
					clientValue := ClientElement.FirstChild.Data

					statusElement := tdElements[5]
					spanElement := statusElement.FindByTag("span")
					statusValue, ok := spanElement.GetAttribute("title")
					if !ok {
						fmt.Printf("span attribute: does not exist")
					}
					//NewLoginTradeDTO(id , titleName , titleID , href , date , startDate , endDate , clientID , status)
					loginTradeDto := dto.NewLoginTradeDTO(url, getHrefID(hrefValue), nameValue, numberValue, hrefValue, dateValue, startValue, endValue, clientValue, statusValue)

					tender := loginTradeDto.GetDataDTO(client)
					tenders = append(tenders, tender)
					if data.IsIn(tendersOldAll, tender) {
						fmt.Println("processGetTenderPage: old tenders contains this", tender)
						return err, tenders, true
					}
				}
			}
		}
	}
	if len(elements) != expectedElementsSize {
		fmt.Printf("done: last page number of elements found: %d, expected number: %d\n", len(elements), expectedElementsSize)
		return err, tenders, true
	}
	return err, tenders, false
}
