package main

import (
	"fmt"
	"tender/bk_page"
	"tender/dto"
	"tender/interfaces/data"
	"tender/login_trade_page"
	"tender/order_page"
	"tender/process"
	"tender/tender_page"
)

func main() {
	flags := dto.NewFlagDTO()
	common := make([]data.Data, 0)
	common = append(common, processTenders(flags)...)
	common = append(common, processOrders(flags)...)
	common = append(common, processBK(flags)...)
	common = append(common, processFrog(flags)...)
	common = append(common, processAnimex(flags)...)
	processCommon(flags, common)
}
func processAnimex(flags *dto.FlagDTO) []data.Data {
	url := "https://grupasmithfield.logintrade.net/"
	return processLoginTrade("animex", url, flags)
}

func processFrog(flags *dto.FlagDTO) []data.Data {
	url := "https://zabka.logintrade.net/"
	return processLoginTrade("frog", url, flags)
}

func processLoginTrade(client string, url string, flags *dto.FlagDTO) []data.Data {
	var err error
	var done bool
	fmt.Println(client + " login trade START ***")
	tenders := make([]data.Data, 0)
	tendersOldAll := make([]data.Data, 0)
	err, tendersOldAll = process.ReadOldAllFile(flags.AnimexOldFileName, client, tendersOldAll)
	urlPrefix := url + login_trade_page.DEFAULT_URL_PREFIX
	urlSuffix := login_trade_page.DEFAULT_URL_SUFIX
	err, tenders = login_trade_page.ProcessGetLoginTradePages(client, url, login_trade_page.GetDefaultHrefID, urlPrefix, urlSuffix, flags.AnimexPages, err, tenders, done, tendersOldAll)
	process.ProcessSaveDataToExcel(client, err, tenders, tendersOldAll, flags)
	fmt.Println(client + " login trade END ***")
	return tenders
}

func processBK(flags *dto.FlagDTO) []data.Data {
	var err error
	var done bool
	fmt.Println("bks START")
	tenders := make([]data.Data, 0)
	tendersOldAll := make([]data.Data, 0)
	err, tendersOldAll = process.ReadOldAllFile(flags.BkOldFileName, "bk", tendersOldAll)
	err, tenders = bk_page.ProcessGetBkPages(flags, err, tenders, done, tendersOldAll)
	process.ProcessSaveDataToExcel("bk", err, tenders, tendersOldAll, flags)
	fmt.Println("bks END")
	return tenders
}

func processOrders(flags *dto.FlagDTO) []data.Data {
	var err error
	var done bool
	fmt.Println("orders START")
	tenders := make([]data.Data, 0)
	tendersOldAll := make([]data.Data, 0)
	err, tendersOldAll = process.ReadOldAllFile(flags.OrdersOldFileName, "oferty", tendersOldAll)
	err, tenders = order_page.ProcessGetOrderPages(flags, err, tenders, done, tendersOldAll)
	process.ProcessSaveDataToExcel("oferty", err, tenders, tendersOldAll, flags)
	fmt.Println("orders END")
	return tenders
}

func processTenders(flags *dto.FlagDTO) []data.Data {
	var err error
	var done bool
	fmt.Println("tenders START")
	tenders := make([]data.Data, 0)
	tendersOldAll := make([]data.Data, 0)
	err, tendersOldAll = process.ReadOldAllFile(flags.TenderOldFileName, "przetargi", tendersOldAll)
	err, tenders = tender_page.ProcessGetTenderPages(flags, err, tenders, done, tendersOldAll)
	process.ProcessSaveDataToExcel("przetargi", err, tenders, tendersOldAll, flags)
	fmt.Println("tenders END")
	return tenders
}

func processCommon(flags *dto.FlagDTO, common []data.Data) {
	var err error
	fmt.Println("common START")
	tendersOldAll := make([]data.Data, 0)
	process.ProcessSaveDataToExcel("common", err, common, tendersOldAll, flags)
	fmt.Println("common END")
}
