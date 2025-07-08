package main

import (
	"errors"
	"fmt"
	"github.com/tealeg/xlsx/v3"
	"tender/dto"
	"tender/interfaces/data"
	"tender/order_page"
	"tender/tender_page"
	"time"
)

func main() {
	flags := dto.NewFlagDTO()
	processOrders(flags)
	processTenders(flags)
}

func fileDateStr() string {
	return time.Now().Format("20060102")
}

func processTenders(flags *dto.FlagDTO) {
	var err error
	var done bool
	fmt.Println("tenders START")
	tenders := make([]data.Data, 0)
	tendersOldAll := make([]data.Data, 0)
	err, tendersOldAll = readOldAllFile(flags.TenderOldFileName, "przetargi", tendersOldAll)
	err, tenders = tender_page.ProcessGetTenderPages(flags, err, tenders, done, tendersOldAll)
	processSaveDataToExcel("przetargi", err, tenders, tendersOldAll, flags)
	fmt.Println("tenders END")
}

func processOrders(flags *dto.FlagDTO) {
	var err error
	var done bool
	fmt.Println("orders START")
	tenders := make([]data.Data, 0)
	tendersOldAll := make([]data.Data, 0)
	err, tendersOldAll = readOldAllFile(flags.OrdersOldFileName, "oferty", tendersOldAll)
	err, tenders = order_page.ProcessGetOrderPages(flags, err, tenders, done, tendersOldAll)
	processSaveDataToExcel("oferty", err, tenders, tendersOldAll, flags)
	fmt.Println("orders END")
}

func processSaveDataToExcel(filename string, err error, tenders, tendersOldAll []data.Data, flags *dto.FlagDTO) {
	var fileAll *xlsx.File
	var fileIT *xlsx.File

	fmt.Println("processSaveDataToExcel")

	fileIT = xlsx.NewFile()
	err = processSaveITDataToExcel(filename+" IT", fileIT, tenders)

	err = fileIT.Save(filename + "_it_" + fileDateStr() + ".xlsx")
	if err != nil {
		fmt.Printf(err.Error())
	}

	if flags.SaveAll {
		if flags.AppendAll {
			for _, tender := range tendersOldAll {
				if !data.IsIn(tenders, tender) {
					tenders = append(tenders, tender)
				} else {
					fmt.Println("processSaveDataToExcel saveAll: there is already this old tender")
				}
			}
		}
		fileAll = xlsx.NewFile()
		err = processSaveAllToExcel(filename, tenders, fileAll)
		err = fileAll.Save(filename + ".xlsx")
		if err != nil {
			fmt.Printf(err.Error())
		}
		err = fileAll.Save(filename + "_" + fileDateStr() + ".xlsx")
		if err != nil {
			fmt.Printf(err.Error())
		}
	}
}

func readOldAllFile(fileName string, sheetName string, tendersOldAll []data.Data) (error, []data.Data) {
	fileOldAll, err := xlsx.OpenFile(fileName)
	if err == nil {
		tendersOldAll = readOldAll(sheetName, fileOldAll, tendersOldAll)
	} else {
		fmt.Printf("file %s was not found\n", fileName)
		tendersOldAll = make([]data.Data, 0)
	}
	fmt.Printf("tendersOldAll len: %v\n", len(tendersOldAll))
	return err, tendersOldAll
}

func readOldAll(sheetName string, fileOldAll *xlsx.File, tendersOldAll []data.Data) []data.Data {
	sheet, ok := fileOldAll.Sheet[sheetName]
	if !ok {
		panic(errors.New("sheet " + sheetName + " not found"))
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

func oldAllRowVisitor(r *xlsx.Row, tendersOldAll []data.Data) []data.Data {
	nr := 1
	idCell := r.GetCell(nr)
	idValue := idCell.Value
	dateCell := r.GetCell(nr + 1)
	dateValue := dateCell.Value
	hrefCell := r.GetCell(nr + 2)
	hrefValue := hrefCell.Value
	nameCell := r.GetCell(nr + 3)
	nameValue := nameCell.Value
	tender := dto.NewDataDTO(nameValue, hrefValue, dateValue, idValue)
	tendersOldAll = append(tendersOldAll, tender)
	return tendersOldAll
}

func processSaveITDataToExcel(sheetName string, file *xlsx.File, tenders []data.Data) error {
	sheetIT, err := file.AddSheet(sheetName)
	setHeader(0, sheetIT)
	rowIT := 0
	for _, tender := range tenders {
		if tender.IsIT() {
			rowIT++
			setRowData(0, sheetIT, rowIT, tender)
		}
	}
	return err
}

func processSaveAllToExcel(sheetName string, tenders []data.Data, file *xlsx.File) error {
	sheet, err := file.AddSheet(sheetName)
	if err != nil {
		return err
	}
	setAllHeader(sheet)
	rowOther := 0
	for _, tender := range tenders {
		rowOther++
		setRowData(2, sheet, rowOther, tender)
		if tender.IsIT() {
			cell, _ := sheet.Cell(rowOther, 0)
			cell.Value = "IT"
		}
		cell, _ := sheet.Cell(rowOther, 1)
		cell.Value = tender.Id()
	}
	return err
}

func setRowData(startCell int, sheet *xlsx.Sheet, r int, tender data.Data) {
	nr := startCell
	cell, _ := sheet.Cell(r, nr)
	cell.Value = tender.Date()

	cell, _ = sheet.Cell(r, nr+1)
	cell.SetHyperlink(tender.Href(), tender.Href(), "")
	style := cell.GetStyle()
	style.Font.Underline = true
	style.Font.Color = "FF0000FF"
	cell.SetStyle(style)

	cell, _ = sheet.Cell(r, nr+2)
	cell.Value = tender.Name()
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
