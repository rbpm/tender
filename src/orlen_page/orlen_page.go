package orlen_page

import (
	"encoding/json"
	"fmt"
	"github.com/Noooste/azuretls-client"
	"net/http"
	"strconv"
	"tender/dto"
	"tender/interfaces/data"
)

func ProcessGetOrlenPages(flags *dto.FlagDTO, err error, tenders []data.Data, done bool, tendersOldAll []data.Data) (error, []data.Data) {
	session := azuretls.NewSession()
	for page := 1; page <= flags.OrlenPages; page++ {
		fmt.Println("orlen page: ", page)
		err, tenders, done = ProcessGetOrlenPage(page, session, tenders, tendersOldAll)
		if done {
			fmt.Println("done")
			break
		}
	}
	return err, tenders
}

func ProcessGetOrlenPage(page int, session *azuretls.Session, tenders []data.Data, tendersOldAll []data.Data) (error, []data.Data, bool) {
	var pkoDto []dto.OrlenDTO
	results := 15 // do not change this => it is always returned 15 records
	startFrom := (page-1)*results + 1

	searchedDemands := make([]int, 0)
	for _, tender := range tenders {
		tenderId, _ := strconv.Atoi(tender.Id())
		searchedDemands = append(searchedDemands, tenderId)
	}
	requestData := dto.RequestDTO{
		Results:          results,
		StartFrom:        startFrom,
		SearchCategories: make([]int, 0),
		SearchedDemands:  searchedDemands,
		SearchOrgs:       make([]int, 0),
		Latest:           true}

	//requestData := "{\"results\":" + fmt.Sprintf("%v", results) + ",\"startFrom\":" + fmt.Sprintf("%v", startFrom) + ",\"searchCategories\":[],\"searchedDemands\":[],\"searchOrgs\":[],\"latest\":true}"
	// fmt.Println(requestData)
	req := &azuretls.Request{
		Method: http.MethodPost,
		Url:    "https://connect.orlen.pl/app/main/recentDemands?specificOrgCode=", //"https://httpbin.org/post"
		Body:   requestData,
	}

	session.OrderedHeaders = azuretls.OrderedHeaders{
		{"Content-Type", "application/json; charset=utf-8"},
	}
	response, err := session.Do(req)

	if err != nil {
		fmt.Println("error:\n", err)
	}
	// fmt.Println("response:\n", response)
	err = json.Unmarshal([]byte(response.String()), &pkoDto)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("DONE with error response:" + response.String())
		return err, tenders, true
	}
	// fmt.Println("len:", len(pkoDto))
	for _, pko := range pkoDto {
		tender := pko.GetDataDTO()
		if data.IsIn(tenders, tender) {
			fmt.Println("TODO processGetOrlenPage: why there is repeated record: ", tender)
			return err, tenders, true
		}
		tenders = append(tenders, tender)
		if data.IsIn(tendersOldAll, tender) {
			fmt.Println("processGetOrlenPage: old orlen data contains this:", tender)
			return err, tenders, true
		}
	}
	return err, tenders, false
}
