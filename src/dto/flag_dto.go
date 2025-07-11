package dto

import "flag"

type FlagDTO struct {
	SaveAll           bool
	TenderPages       int
	OrderPages        int
	BkPages           int
	LoginTradePages   int
	PkoPages          int
	AppendAll         bool
	TenderOldFileName string
	OrdersOldFileName string
	BkOldFileName     string
	FrogOldFileName   string
	AnimexOldFileName string
	PkoOldFileName    string
}

func NewFlagDTO() *FlagDTO {
	saveAll := flag.Bool("saveAll", true, "save all tenders to excel")
	//max 1048567/12=87381, found 1090*12 =>1000
	tenderPages := flag.Int("tenderPages", 1000, "number of tender Pages to get")
	//max 1000/50=200
	orderPages := flag.Int("orderPages", 200, "number of order Pages to get")
	bkPages := flag.Int("bkPages", 0, "number of bk Pages to get, 0 when total")
	loginTradePages := flag.Int("loginTradePages", 100, "number of Login Trade Pages to get")
	pkoPages := flag.Int("pkoPages", 100, "number of PKO Pages to get")
	appendAll := flag.Bool("appendAll", false, "append old tenders to new all")
	tenderOldFileName := flag.String("tenderOldFileName", "przetargi.xlsx", "tender old file name")
	ordersOldFileName := flag.String("orderOldFileName", "oferty.xlsx", "order old file name")
	bkOldFileName := flag.String("bkOldFileName", "bk.xlsx", "bk old file name")
	frogOldFileName := flag.String("frogOldFileName", "frog.xlsx", "frog old file name")
	animexOldFileName := flag.String("animexOldFileName", "animex.xlsx", "animex old file name")
	pkoOldFileName := flag.String("pkoOldFileName", "pko.xlsx", "pko old file name")
	flags := FlagDTO{
		SaveAll:           *saveAll,
		TenderPages:       *tenderPages,
		OrderPages:        *orderPages,
		BkPages:           *bkPages,
		LoginTradePages:   *loginTradePages,
		PkoPages:          *pkoPages,
		AppendAll:         *appendAll,
		TenderOldFileName: *tenderOldFileName,
		OrdersOldFileName: *ordersOldFileName,
		BkOldFileName:     *bkOldFileName,
		FrogOldFileName:   *frogOldFileName,
		AnimexOldFileName: *animexOldFileName,
		PkoOldFileName:    *pkoOldFileName}
	return &flags
}
