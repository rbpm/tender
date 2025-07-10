package dto

import "flag"

type FlagDTO struct {
	SaveAll           bool
	TenderPages       int
	OrderPages        int
	BkPages           int
	FrogPages         int
	AppendAll         bool
	TenderOldFileName string
	OrdersOldFileName string
	BkOldFileName     string
}

func NewFlagDTO() *FlagDTO {
	saveAll := flag.Bool("saveAll", true, "save all tenders to excel")
	//max 1048567/12=87381, found 1090*12 =>1000
	tenderPages := flag.Int("tenderPages", 1000, "number of tender Pages to get")
	//max 1000/50=200
	orderPages := flag.Int("orderPages", 200, "number of order Pages to get")
	bkPages := flag.Int("bkPages", 0, "number of bk Pages to get, 0 when total")
	frogPages := flag.Int("frogPages", 1, "number of Frog Pages to get")
	appendAll := flag.Bool("appendAll", false, "append old tenders to new all")
	tenderOldFileName := flag.String("tenderOldFileName", "przetargi.xlsx", "tender old file name")
	ordersOldFileName := flag.String("orderOldFileName", "oferty.xlsx", "order old file name")
	bkOldFileName := flag.String("bkOldFileName", "bk.xlsx", "bk old file name")
	flags := FlagDTO{
		SaveAll:           *saveAll,
		TenderPages:       *tenderPages,
		OrderPages:        *orderPages,
		BkPages:           *bkPages,
		FrogPages:         *frogPages,
		AppendAll:         *appendAll,
		TenderOldFileName: *tenderOldFileName,
		OrdersOldFileName: *ordersOldFileName,
		BkOldFileName:     *bkOldFileName}
	return &flags
}
