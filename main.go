package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/Noooste/azuretls-client"
	"github.com/gurkankaymak/gosoup"
	"github.com/tealeg/xlsx/v3"
	"strings"
	"time"
)

type flagDTO struct {
	saveAll           bool
	tenderPages       int
	orderPages        int
	appendAll         bool
	tenderOldFileName string
	ordersOldFileName string
}

func newFlagDTO() *flagDTO {
	saveAll := flag.Bool("saveAll", true, "save all tenders to excel")
	//max 1048567/12=87381, found 1090*12 =>1000
	tenderPages := flag.Int("tenderPages", 1000, "number of tender Pages to get")
	//max 1000/50=200!!!
	orderPages := flag.Int("orderPages", 200, "number of order Pages to get")
	appendAll := flag.Bool("appendAll", false, "append old tenders to new all")
	tenderOldFileName := flag.String("tenderOldFileName", "przetargi_all.xlsx", "tender old file name")
	ordersOldFileName := flag.String("orderOldFileName", "oferty_all.xlsx", "order old file name")
	flags := flagDTO{
		saveAll:           *saveAll,
		tenderPages:       *tenderPages,
		orderPages:        *orderPages,
		appendAll:         *appendAll,
		tenderOldFileName: *tenderOldFileName,
		ordersOldFileName: *ordersOldFileName}
	return &flags
}

type tenderDTO struct {
	name string
	href string
	date string
	id   string
	isIT bool
}

func newTenderDTO(name string, href string, date string, id string) *tenderDTO {
	lowerName := strings.ToLower(name)
	isIT := strings.Contains(lowerName, "oprogramowani") ||
		strings.Contains(lowerName, " it ") ||
		strings.Contains(lowerName, "rozwój i utrzymanie systemu") ||
		strings.Contains(lowerName, "aplikacj")
	p := tenderDTO{name: name, href: href, date: date, id: id, isIT: isIT}
	return &p
}

type orderDTO struct {
	ObjectId                string    `json:"objectId"`
	Title                   string    `json:"title"`
	OrganizationId          string    `json:"organizationId"`
	OrganizationName        string    `json:"organizationName"`
	OrganizationPartName    string    `json:"organizationPartName"`
	OrganizationCity        string    `json:"organizationCity"`
	OrganizationProvince    string    `json:"organizationProvince"`
	BzpNumber               string    `json:"bzpNumber"`
	TenderType              string    `json:"tenderType"`
	CompetitionType         string    `json:"competitionType"`
	ConcessionType          string    `json:"concessionType"`
	SubmissionOffersDate    time.Time `json:"submissionOffersDate"`
	TenderState             string    `json:"tenderState"`
	IsTenderAmountBelowEU   bool      `json:"isTenderAmountBelowEU"`
	TedContractNoticeNumber string    `json:"tedContractNoticeNumber"`
	InitiationDate          time.Time `json:"initiationDate"`
}

func (order orderDTO) getTenderDTO() *tenderDTO {
	href := "https://ezamowienia.gov.pl/mp-client/search/list/" + order.ObjectId
	return newTenderDTO(order.Title, href, order.SubmissionOffersDate.Format("2006.01.02"), order.ObjectId)
}

func main() {
	flags := newFlagDTO()
	processTenders(flags)
	processOrders(flags)
}

func fileDateStr() string {
	return time.Now().Format("20060102")
}

func processTenders(flags *flagDTO) {
	var err error
	var done bool
	var fileOldAll *xlsx.File

	tenders := make([]*tenderDTO, 0)
	tendersIT := make([]*tenderDTO, 0)
	tendersOldAll := make([]*tenderDTO, 0)

	fileOldAll, err = xlsx.OpenFile(flags.tenderOldFileName)
	if err == nil {
		tendersOldAll = readOldAll("przetargi", fileOldAll, tendersOldAll)
	} else {
		fmt.Printf("file %s was not found\n", flags.tenderOldFileName)
		tendersOldAll = make([]*tenderDTO, 0)
	}
	fmt.Printf("tendersOldAll len: %v\n", len(tendersOldAll))

	session := azuretls.NewSession()
	for page := 1; page <= flags.tenderPages; page++ {
		fmt.Println("tender page: ", page)
		err, tendersIT, tenders, done = processGetTenderPage(page, session, tendersIT, tenders, tendersOldAll)
		if done {
			fmt.Println("done")
			break
		}
	}
	processSaveDataToExcel("przetargi", err, tenders, tendersIT, tendersOldAll, flags)
	fmt.Println("tenders END")
}

func processOrders(flags *flagDTO) {
	var err error
	var done bool
	var fileOldAll *xlsx.File

	fmt.Println("orders START")
	tenders := make([]*tenderDTO, 0)
	tendersIT := make([]*tenderDTO, 0)
	tendersOldAll := make([]*tenderDTO, 0)

	fileOldAll, err = xlsx.OpenFile(flags.ordersOldFileName)
	if err == nil {
		tendersOldAll = readOldAll("oferty", fileOldAll, tendersOldAll)
	} else {
		fmt.Printf("file %s was not found\n", flags.ordersOldFileName)
		tendersOldAll = make([]*tenderDTO, 0)
	}
	fmt.Printf("ordersOldAll len: %v\n", len(tendersOldAll))

	session := azuretls.NewSession()
	for page := 1; page <= flags.orderPages; page++ {
		fmt.Println("order page: ", page)
		err, tendersIT, tenders, done = processGetOrderPage(page, session, tendersIT, tenders, tendersOldAll)
		if done {
			fmt.Println("done")
			break
		}
	}
	processSaveDataToExcel("oferty", err, tenders, tendersIT, tendersOldAll, flags)
	fmt.Println("orders END")
}

func processSaveDataToExcel(filename string, err error, tenders, tendersIT, tendersOldAll []*tenderDTO, flags *flagDTO) {
	var fileAll *xlsx.File
	var fileIT *xlsx.File

	fmt.Println("processSaveDataToExcel")

	fileIT = xlsx.NewFile()
	err = processSaveToExcel(filename+" IT", fileIT, tendersIT)

	err = fileIT.Save(filename + "_it_" + fileDateStr() + ".xlsx")
	if err != nil {
		fmt.Printf(err.Error())
	}

	if flags.saveAll {
		if flags.appendAll {
			for _, tender := range tendersOldAll {
				if !contains(tenders, tender) {
					tenders = append(tenders, tender)
				} else {
					fmt.Println("processSaveDataToExcel saveAll: there is already this old tender")
				}
			}
		}
		fileAll = xlsx.NewFile()
		err = processSaveAllToExcel(filename, tenders, fileAll)
		err = fileAll.Save(filename + "_all.xlsx")
		if err != nil {
			fmt.Printf(err.Error())
		}
		err = fileAll.Save(filename + "_all_" + fileDateStr() + ".xlsx")
		if err != nil {
			fmt.Printf(err.Error())
		}
	}
}

func readOldAll(sheetName string, fileOldAll *xlsx.File, tendersOldAll []*tenderDTO) []*tenderDTO {
	sheet, ok := fileOldAll.Sheet[sheetName]
	if !ok {
		panic(errors.New("Sheet tenders not found"))
	}
	fmt.Println("Max row is", sheet.MaxRow)
	for row := 1; row < sheet.MaxRow; row++ {
		r, err := sheet.Row(row)
		if err != nil {
			panic(err)
		}
		tendersOldAll = oldAllRowVisitor(r, tendersOldAll)
	}
	return tendersOldAll
}

func oldAllRowVisitor(r *xlsx.Row, tendersOldAll []*tenderDTO) []*tenderDTO {
	nr := 1
	idCell := r.GetCell(nr)
	idValue := idCell.Value
	dateCell := r.GetCell(nr + 1)
	dateValue := dateCell.Value
	hrefCell := r.GetCell(nr + 2)
	hrefValue := hrefCell.Value
	nameCell := r.GetCell(nr + 3)
	nameValue := nameCell.Value
	tender := newTenderDTO(nameValue, hrefValue, dateValue, idValue)
	tendersOldAll = append(tendersOldAll, tender)
	return tendersOldAll
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

func processSaveToExcel(sheetName string, file *xlsx.File, tendersIT []*tenderDTO) error {
	sheetIT, err := file.AddSheet(sheetName)
	setHeader(0, sheetIT)
	rowIT := 0
	for _, tendersT := range tendersIT {
		rowIT++
		setRowData(0, sheetIT, rowIT, tendersT)
	}
	return err
}

func processSaveAllToExcel(sheetName string, tenders []*tenderDTO, file *xlsx.File) error {
	sheet, err := file.AddSheet(sheetName)
	setAllHeader(sheet)
	rowOther := 0
	for _, tender := range tenders {
		rowOther++
		setRowData(2, sheet, rowOther, tender)
		if tender.isIT {
			cell, _ := sheet.Cell(rowOther, 0)
			cell.Value = "IT"
		}
		cell, _ := sheet.Cell(rowOther, 1)
		cell.Value = tender.id
	}
	return err
}

func setRowData(startCell int, sheet *xlsx.Sheet, r int, tender *tenderDTO) {
	nr := startCell
	cell, _ := sheet.Cell(r, nr)
	cell.Value = tender.date

	cell, _ = sheet.Cell(r, nr+1)
	cell.SetHyperlink(tender.href, tender.href, "")
	style := cell.GetStyle()
	style.Font.Underline = true
	style.Font.Color = "FF0000FF"
	cell.SetStyle(style)

	cell, _ = sheet.Cell(r, nr+2)
	cell.Value = tender.name
}

func setHeader(startCell int, sheet *xlsx.Sheet) {
	nr := startCell
	cell, _ := sheet.Cell(0, nr)
	cell.Value = "data"
	sheet.SetColWidth(nr+1, nr+1, 10)

	nr++
	cell, _ = sheet.Cell(0, nr)
	cell.Value = "link"
	sheet.SetColWidth(nr+1, nr+1, 18)

	nr++
	cell, _ = sheet.Cell(0, nr)
	cell.Value = "nazwa"
	sheet.SetColWidth(nr+1, nr+1, 100)
}

func setAllHeader(sheet *xlsx.Sheet) {
	setHeader(2, sheet)
	cell, _ := sheet.Cell(0, 0)
	cell.Value = "IT"
	sheet.SetColWidth(1, 1, 3)
	cell, _ = sheet.Cell(0, 1)
	cell.Value = "ID"
	sheet.SetColWidth(2, 2, 12.5)
}

func processGetTenderPage(page int, session *azuretls.Session, tendersIT []*tenderDTO, tenders []*tenderDTO, tendersOldAll []*tenderDTO) (error, []*tenderDTO, []*tenderDTO, bool) {
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
		tender := newTenderDTO(nameValue, hrefValue, dateValue, getHrefID(hrefValue))
		tendersIT, tenders = appendTender(tender, tendersIT, tenders)
		if contains(tendersOldAll, tender) {
			fmt.Println("processGetTenderPage: old tenders contains this", tender)
			return err, tendersIT, tenders, true
		}
	}
	return err, tendersIT, tenders, false
}

func appendTender(tender *tenderDTO, tendersIT, tenders []*tenderDTO) ([]*tenderDTO, []*tenderDTO) {
	if tender.isIT {
		tendersIT = append(tendersIT, tender)
		tenders = append(tenders, tender)
	} else {
		tenders = append(tenders, tender)
	}
	return tendersIT, tenders
}

func processGetOrderPage(page int, session *azuretls.Session, tendersIT []*tenderDTO, tenders []*tenderDTO, tendersOldAll []*tenderDTO) (error, []*tenderDTO, []*tenderDTO, bool) {
	var orders []orderDTO
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
		tender := order.getTenderDTO()
		tendersIT, tenders = appendTender(tender, tendersIT, tenders)
		if contains(tendersOldAll, tender) {
			fmt.Println("processGetOrderPage: old orders contains this:", tender)
			return err, tendersIT, tenders, true
		}
	}
	return err, tendersIT, tenders, false
}

func contains(tenders []*tenderDTO, tender *tenderDTO) bool {
	for _, p := range tenders {
		if p.id == tender.id && p.name == tender.name {
			fmt.Printf("id=%s;", p.id)
			return true
		}
	}
	return false
}
