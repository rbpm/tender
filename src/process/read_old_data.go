package process

import (
	"errors"
	"fmt"
	"tender/dto"
	"tender/interfaces/data"
	"time"

	"github.com/tealeg/xlsx/v3"
)

func ReadOldAllFile(fileName string, sheetName string, tendersOldAll []data.Data) (error, []data.Data) {
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

	nameCell := r.GetCell(nr + 1)
	nameValue := nameCell.Value

	srcCell := r.GetCell(nr + 4)
	srcValue := srcCell.Value

	hrefCell := r.GetCell(nr + 5)
	hrefValue := hrefCell.Value

	dateCell := r.GetCell(nr + 6)
	datePtr := getCellTime(dateCell)

	tender := dto.NewDataDTO(srcValue, nameValue, hrefValue, datePtr, idValue)
	tendersOldAll = append(tendersOldAll, tender)
	return tendersOldAll
}

func getCellTime(dateCell *xlsx.Cell) *time.Time {
	if dateValue, err := dateCell.GetTime(false); err == nil {
		return &dateValue
	} else {
		fmt.Println("process.getCellTime() error:", err.Error())
		return nil
	}
}
