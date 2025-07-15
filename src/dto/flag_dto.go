package dto

import "flag"

type FlagDTO struct {
	SaveAll           bool
	OneplacePages     int
	EzamowieniaPages  int
	BkPages           int
	LoginTradePages   int
	PkoPages          int
	OrlenPages        int
	AppendAll         bool
	ExcelDir          string
	TenderOldFileName string
	OrdersOldFileName string
	BkOldFileName     string
	FrogOldFileName   string
	AnimexOldFileName string
	PkoOldFileName    string
	OrlenOldFileName  string
}

func NewFlagDTO() *FlagDTO {
	saveAll := flag.Bool("saveAll", true, "save all tenders to excel")
	//max 1048567/12=87381, found 1090*12 =>1000
	oneplacePages := flag.Int("oneplacePages", 1000, "number of oneplace Pages to get")
	//max 1000/50=200
	ezamowieniaPages := flag.Int("ezamowieniaPages", 200, "number of ezamowienia Pages to get")
	bkPages := flag.Int("bkPages", 0, "number of bk Pages to get, 0 when total")
	loginTradePages := flag.Int("loginTradePages", 100, "number of Login Trade Pages to get")
	pkoPages := flag.Int("pkoPages", 100, "number of PKO Pages to get")
	orlenPages := flag.Int("orlenPages", 100, "number of Orlen Pages to get")
	appendAll := flag.Bool("appendAll", false, "append old tenders to new all")
	excelDir := flag.String("excelDir", "excel/", "excel directory")
	oneplaceOldFileName := flag.String("oneplaceOldFileName", "oneplace.xlsx", "oneplace old file name")
	ordersOldFileName := flag.String("ezamowieniaOldFileName", "ezamowienia.xlsx", "ezamowienia old file name")
	bkOldFileName := flag.String("bkOldFileName", "bk.xlsx", "bk old file name")
	frogOldFileName := flag.String("frogOldFileName", "frog.xlsx", "frog old file name")
	animexOldFileName := flag.String("animexOldFileName", "animex.xlsx", "animex old file name")
	pkoOldFileName := flag.String("pkoOldFileName", "pko.xlsx", "pko old file name")
	orlenOldFileName := flag.String("orlenOldFileName", "orlen.xlsx", "Orlen old file name")
	flags := FlagDTO{
		SaveAll:           *saveAll,
		OneplacePages:     *oneplacePages,
		EzamowieniaPages:  *ezamowieniaPages,
		BkPages:           *bkPages,
		LoginTradePages:   *loginTradePages,
		PkoPages:          *pkoPages,
		OrlenPages:        *orlenPages,
		AppendAll:         *appendAll,
		ExcelDir:          *excelDir,
		TenderOldFileName: *oneplaceOldFileName,
		OrdersOldFileName: *ordersOldFileName,
		BkOldFileName:     *bkOldFileName,
		FrogOldFileName:   *frogOldFileName,
		AnimexOldFileName: *animexOldFileName,
		PkoOldFileName:    *pkoOldFileName,
		OrlenOldFileName:  *orlenOldFileName}
	return &flags
}
