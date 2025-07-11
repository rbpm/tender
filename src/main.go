package main

import (
	"errors"
	"fmt"
	"os"
	"tender/bk_page"
	"tender/dto"
	"tender/interfaces/data"
	"tender/login_trade_page"
	"tender/order_page"
	"tender/pko_page"
	"tender/process"
	"tender/tender_page"
)

func main() {
	flags := dto.NewFlagDTO()
	mkDirIfNotExist(flags.ExcelDir)
	common := make([]data.Data, 0)
	common = append(common, processTenders(flags)...)
	common = append(common, processOrders(flags)...)
	common = append(common, processBK(flags)...)
	common = append(common, processFrog(flags)...)
	common = append(common, processAnimex(flags)...)
	common = append(common, processBosbank(flags)...)
	common = append(common, processGemetica(flags)...)
	common = append(common, processPko(flags)...)
	processCommon(flags, common)
}

func processPko(flags *dto.FlagDTO) []data.Data {
	var err error
	var done bool
	fmt.Println("pko START")
	tenders := make([]data.Data, 0)
	tendersOldAll := make([]data.Data, 0)
	err, tendersOldAll = process.ReadOldAllFile(flags.ExcelDir+flags.PkoOldFileName, "pko", tendersOldAll)
	err, tenders = pko_page.ProcessGetPkoPages(flags, err, tenders, done, tendersOldAll)
	process.ProcessSaveDataToExcel("pko", err, tenders, tendersOldAll, flags)
	fmt.Println("pko END")
	return tenders
}

func processGemetica(flags *dto.FlagDTO) []data.Data {
	url := "https://platforma-qemetica.logintrade.net/"
	return processLoginTrade("qemetica", url, flags, "qemetica.xlsx")
}

func processBosbank(flags *dto.FlagDTO) []data.Data {
	url := "https://bosbank.logintrade.net/"
	return processLoginTrade("bosbank", url, flags, "bosbank.xlsx")
}

//func processLotams(flags *dto.FlagDTO) []data.Data {
//	url := "https://lotams.logintrade.net/"
//	return processLoginTradeNET("lotams", url, flags)
//}

// no such page => NET version
// https://cersanit.logintrade.net/portal,listaZapytaniaOfertowe.html?status_realizacji_zapytania[]=oczekiwanie_ofert&wojewodztwo=wszystkie&search=&search_sort=9&page=1&itemsperpage=100
//func processCersanit(flags *dto.FlagDTO) []data.Data {
//	url := "https://cersanit.logintrade.net/"
//	return processLoginTradeNET("cersanit", url, flags)
//}

func processAnimex(flags *dto.FlagDTO) []data.Data {
	url := "https://grupasmithfield.logintrade.net/"
	return processLoginTrade("animex", url, flags, "animex.xlsx")
}

func processFrog(flags *dto.FlagDTO) []data.Data {
	url := "https://zabka.logintrade.net/"
	return processLoginTrade("frog", url, flags, "frog.xlsx")
}

func processLoginTrade(client string, url string, flags *dto.FlagDTO, oldFileName string) []data.Data {
	var err error
	var done bool
	fmt.Println(client + " login trade START ***")
	tenders := make([]data.Data, 0)
	tendersOldAll := make([]data.Data, 0)
	err, tendersOldAll = process.ReadOldAllFile(flags.ExcelDir+oldFileName, client, tendersOldAll)
	urlPrefix := url + login_trade_page.DEFAULT_URL_PREFIX
	urlSuffix := login_trade_page.DEFAULT_URL_SUFIX
	err, tenders = login_trade_page.ProcessGetLoginTradePages(client, url, login_trade_page.GetDefaultHrefID, urlPrefix, urlSuffix, flags.LoginTradePages, err, tenders, done, tendersOldAll)
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
	err, tendersOldAll = process.ReadOldAllFile(flags.ExcelDir+flags.BkOldFileName, "bk", tendersOldAll)
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
	err, tendersOldAll = process.ReadOldAllFile(flags.ExcelDir+flags.OrdersOldFileName, "oferty", tendersOldAll)
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
	err, tendersOldAll = process.ReadOldAllFile(flags.ExcelDir+flags.TenderOldFileName, "przetargi", tendersOldAll)
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

func mkDirIfNotExist(path string) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			fmt.Println(err)
		}
	}
}
