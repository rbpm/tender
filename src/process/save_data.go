package process

import (
	"fmt"
	"github.com/tealeg/xlsx/v3"
	"tender/dto"
	"tender/interfaces/data"
	"time"
)

func ProcessSaveDataToExcel(filename string, err error, tenders, tendersOldAll []data.Data, flags *dto.FlagDTO) {
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

func fileDateStr() string {
	return time.Now().Format("20060102")
}

func dateStr() string {
	return time.Now().Format("2006-01-02")
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
	cell.Value = tender.Name()

	cell, _ = sheet.Cell(r, nr+1)
	cell.Value = "WB"

	cell, _ = sheet.Cell(r, nr+2)
	cell.Value = dateStr()

	cell, _ = sheet.Cell(r, nr+3)
	cell.Value = tender.Src()

	cell, _ = sheet.Cell(r, nr+4)
	cell.SetHyperlink(tender.Href(), tender.Href(), "")
	style := cell.GetStyle()
	style.Font.Underline = true
	style.Font.Color = "FF0000FF"
	cell.SetStyle(style)

	cell, _ = sheet.Cell(r, nr+5)
	cell.Value = tender.Date()

}

func setHeader(startCell int, sheet *xlsx.Sheet) {
	nr := startCell
	cell, _ := sheet.Cell(0, nr)
	cell.Value = "Przetarg"
	sheet.SetColWidth(nr+1, nr+1, 75)

	nr++
	cell, _ = sheet.Cell(0, nr)
	cell.Value = "Osoba, która zgłosiła"
	sheet.SetColWidth(nr+1, nr+1, 25)

	nr++
	cell, _ = sheet.Cell(0, nr)
	cell.Value = "Data dodania"
	sheet.SetColWidth(nr+1, nr+1, 25)

	nr++
	cell, _ = sheet.Cell(0, nr)
	cell.Value = "Klient"
	sheet.SetColWidth(nr+1, nr+1, 25)

	nr++
	cell, _ = sheet.Cell(0, nr)
	cell.Value = "Źródło przetargu"
	sheet.SetColWidth(nr+1, nr+1, 25)

	nr++
	cell, _ = sheet.Cell(0, nr)
	cell.Value = "Deadline na złożenie oferty"
	sheet.SetColWidth(nr+1, nr+1, 25)
}

func setAllHeader(sheet *xlsx.Sheet) {
	setHeader(2, sheet)
	cell, _ := sheet.Cell(0, 0)
	cell.Value = "IT"
	sheet.SetColWidth(1, 1, 3)
	cell, _ = sheet.Cell(0, 1)
	cell.Value = "ID"
	sheet.SetColWidth(2, 2, 40)
}
