package dto

import "flag"

type FlagDTO struct {
	SaveAll           bool
	TenderPages       int
	OrderPages        int
	AppendAll         bool
	TenderOldFileName string
	OrdersOldFileName string
}

func NewFlagDTO() *FlagDTO {
	saveAll := flag.Bool("saveAll", true, "save all tenders to excel")
	//max 1048567/12=87381, found 1090*12 =>1000
	tenderPages := flag.Int("tenderPages", 3, "number of tender Pages to get")
	//max 1000/50=200
	orderPages := flag.Int("orderPages", 3, "number of order Pages to get")
	appendAll := flag.Bool("appendAll", false, "append old tenders to new all")
	tenderOldFileName := flag.String("tenderOldFileName", "przetargi.xlsx", "tender old file name")
	ordersOldFileName := flag.String("orderOldFileName", "oferty.xlsx", "order old file name")
	flags := FlagDTO{
		SaveAll:           *saveAll,
		TenderPages:       *tenderPages,
		OrderPages:        *orderPages,
		AppendAll:         *appendAll,
		TenderOldFileName: *tenderOldFileName,
		OrdersOldFileName: *ordersOldFileName}
	return &flags
}
