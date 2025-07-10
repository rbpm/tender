package main

import (
	"fmt"
	"tender/bk_page"
	"tender/dto"
	"tender/frog_page"
	"tender/interfaces/data"
	"tender/order_page"
	"tender/process"
	"tender/tender_page"
)

func main() {
	flags := dto.NewFlagDTO()
	processFrog(flags)
	processBK(flags)
	processOrders(flags)
	processTenders(flags)
}

func processFrog(flags *dto.FlagDTO) {
	var err error
	var done bool
	fmt.Println("frog START")
	tenders := make([]data.Data, 0)
	tendersOldAll := make([]data.Data, 0)
	//err, tendersOldAll = process.ReadOldAllFile(flags.TenderOldFileName, "frog", tendersOldAll)
	err, tenders = frog_page.ProcessGetFrogPages(flags, err, tenders, done, tendersOldAll)
	process.ProcessSaveDataToExcel("frog", err, tenders, tendersOldAll, flags)
	fmt.Println("frog END")
}

func processTenders(flags *dto.FlagDTO) {
	var err error
	var done bool
	fmt.Println("tenders START")
	tenders := make([]data.Data, 0)
	tendersOldAll := make([]data.Data, 0)
	err, tendersOldAll = process.ReadOldAllFile(flags.TenderOldFileName, "przetargi", tendersOldAll)
	err, tenders = tender_page.ProcessGetTenderPages(flags, err, tenders, done, tendersOldAll)
	process.ProcessSaveDataToExcel("przetargi", err, tenders, tendersOldAll, flags)
	fmt.Println("tenders END")
}

func processOrders(flags *dto.FlagDTO) {
	var err error
	var done bool
	fmt.Println("orders START")
	tenders := make([]data.Data, 0)
	tendersOldAll := make([]data.Data, 0)
	err, tendersOldAll = process.ReadOldAllFile(flags.OrdersOldFileName, "oferty", tendersOldAll)
	err, tenders = order_page.ProcessGetOrderPages(flags, err, tenders, done, tendersOldAll)
	process.ProcessSaveDataToExcel("oferty", err, tenders, tendersOldAll, flags)
	fmt.Println("orders END")
}

func processBK(flags *dto.FlagDTO) {
	var err error
	var done bool
	fmt.Println("bks START")
	tenders := make([]data.Data, 0)
	tendersOldAll := make([]data.Data, 0)
	err, tendersOldAll = process.ReadOldAllFile(flags.BkOldFileName, "bk", tendersOldAll)
	err, tenders = bk_page.ProcessGetBkPages(flags, err, tenders, done, tendersOldAll)
	process.ProcessSaveDataToExcel("bk", err, tenders, tendersOldAll, flags)
	fmt.Println("bks END")
}
