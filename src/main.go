package main

import (
	"errors"
	"fmt"
	"github.com/Noooste/azuretls-client"
	"github.com/tealeg/xlsx/v3"
	"tender/dto"
	"tender/order_page"
	"tender/tender_page"
	"time"
)

func main() {
	flags := dto.NewFlagDTO()
	processTenders(flags)
	processOrders(flags)
}

func fileDateStr() string {
	return time.Now().Format("20060102")
}

func processTenders(flags *dto.FlagDTO) {
	var err error
	var done bool
	var fileOldAll *xlsx.File

	tenders := make([]*dto.TenderDTO, 0)
	tendersIT := make([]*dto.TenderDTO, 0)
	tendersOldAll := make([]*dto.TenderDTO, 0)

	fileOldAll, err = xlsx.OpenFile(flags.TenderOldFileName)
	if err == nil {
		tendersOldAll = readOldAll("przetargi", fileOldAll, tendersOldAll)
	} else {
		fmt.Printf("file %s was not found\n", flags.TenderOldFileName)
		tendersOldAll = make([]*dto.TenderDTO, 0)
	}
	fmt.Printf("tendersOldAll len: %v\n", len(tendersOldAll))

	session := azuretls.NewSession()
	for page := 1; page <= flags.TenderPages; page++ {
		fmt.Println("tender page: ", page)
		err, tendersIT, tenders, done = tender_page.ProcessGetTenderPage(page, session, tendersIT, tenders, tendersOldAll)
		if done {
			fmt.Println("done")
			break
		}
	}
	processSaveDataToExcel("przetargi", err, tenders, tendersIT, tendersOldAll, flags)
	fmt.Println("tenders END")
}

func processOrders(flags *dto.FlagDTO) {
	var err error
	var done bool
	var fileOldAll *xlsx.File

	fmt.Println("orders START")
	tenders := make([]*dto.TenderDTO, 0)
	tendersIT := make([]*dto.TenderDTO, 0)
	tendersOldAll := make([]*dto.TenderDTO, 0)

	fileOldAll, err = xlsx.OpenFile(flags.OrdersOldFileName)
	if err == nil {
		tendersOldAll = readOldAll("oferty", fileOldAll, tendersOldAll)
	} else {
		fmt.Printf("file %s was not found\n", flags.OrdersOldFileName)
		tendersOldAll = make([]*dto.TenderDTO, 0)
	}
	fmt.Printf("ordersOldAll len: %v\n", len(tendersOldAll))

	session := azuretls.NewSession()
	for page := 1; page <= flags.OrderPages; page++ {
		fmt.Println("order page: ", page)
		err, tendersIT, tenders, done = order_page.ProcessGetOrderPage(page, session, tendersIT, tenders, tendersOldAll)
		if done {
			fmt.Println("done")
			break
		}
	}
	processSaveDataToExcel("oferty", err, tenders, tendersIT, tendersOldAll, flags)
	fmt.Println("orders END")
}

func processSaveDataToExcel(filename string, err error, tenders, tendersIT, tendersOldAll []*dto.TenderDTO, flags *dto.FlagDTO) {
	var fileAll *xlsx.File
	var fileIT *xlsx.File

	fmt.Println("processSaveDataToExcel")

	fileIT = xlsx.NewFile()
	err = processSaveToExcel(filename+" IT", fileIT, tendersIT)

	err = fileIT.Save(filename + "_it_" + fileDateStr() + ".xlsx")
	if err != nil {
		fmt.Printf(err.Error())
	}

	if flags.SaveAll {
		if flags.AppendAll {
			for _, tender := range tendersOldAll {
				if !tender.IsIn(tenders) {
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

func readOldAll(sheetName string, fileOldAll *xlsx.File, tendersOldAll []*dto.TenderDTO) []*dto.TenderDTO {
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

func oldAllRowVisitor(r *xlsx.Row, tendersOldAll []*dto.TenderDTO) []*dto.TenderDTO {
	nr := 1
	idCell := r.GetCell(nr)
	idValue := idCell.Value
	dateCell := r.GetCell(nr + 1)
	dateValue := dateCell.Value
	hrefCell := r.GetCell(nr + 2)
	hrefValue := hrefCell.Value
	nameCell := r.GetCell(nr + 3)
	nameValue := nameCell.Value
	tender := dto.NewTenderDTO(nameValue, hrefValue, dateValue, idValue)
	tendersOldAll = append(tendersOldAll, tender)
	return tendersOldAll
}

func processSaveToExcel(sheetName string, file *xlsx.File, tendersIT []*dto.TenderDTO) error {
	sheetIT, err := file.AddSheet(sheetName)
	setHeader(0, sheetIT)
	rowIT := 0
	for _, tendersT := range tendersIT {
		rowIT++
		setRowData(0, sheetIT, rowIT, tendersT)
	}
	return err
}

func processSaveAllToExcel(sheetName string, tenders []*dto.TenderDTO, file *xlsx.File) error {
	sheet, err := file.AddSheet(sheetName)
	setAllHeader(sheet)
	rowOther := 0
	for _, tender := range tenders {
		rowOther++
		setRowData(2, sheet, rowOther, tender)
		if tender.IsIT {
			cell, _ := sheet.Cell(rowOther, 0)
			cell.Value = "IT"
		}
		cell, _ := sheet.Cell(rowOther, 1)
		cell.Value = tender.Id
	}
	return err
}

func setRowData(startCell int, sheet *xlsx.Sheet, r int, tender *dto.TenderDTO) {
	nr := startCell
	cell, _ := sheet.Cell(r, nr)
	cell.Value = tender.Date

	cell, _ = sheet.Cell(r, nr+1)
	cell.SetHyperlink(tender.Href, tender.Href, "")
	style := cell.GetStyle()
	style.Font.Underline = true
	style.Font.Color = "FF0000FF"
	cell.SetStyle(style)

	cell, _ = sheet.Cell(r, nr+2)
	cell.Value = tender.Name
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
