package dto

import (
	"fmt"
	"time"
)

// {"results":15,"startFrom":1,"searchCategories":[],"searchedDemands":[],"searchOrgs":[],"latest":true}
type RequestDTO struct {
	Results          int   `json:"results"`
	StartFrom        int   `json:"startFrom"`
	SearchCategories []int `json:"searchCategories"`
	SearchedDemands  []int `json:"searchedDemands"`
	SearchOrgs       []int `json:"searchOrgs"`
	Latest           bool  `json:"latest"`
}

type OrlenDTO struct {
	Number               string      `json:"number"`
	EndDate              int64       `json:"endDate"`
	Org                  OrlenOrgDTO `json:"org"`
	Identity             int         `json:"identity"`
	Kind                 string      `json:"kind"`
	Name                 string      `json:"name"`
	SupplierExistInRound bool        `json:"supplierExistInRound"`
	Category             string      `json:"category"`
}

type OrlenOrgDTO struct {
	Country string `json:"country"`
	Name    string `json:"name"`
	Logo    string `json:"logo"`
	OrgId   string `json:"orgId"`
}

func (orlenDto OrlenDTO) GetDataDTO() *DataDTO {
	href := "https://connect.orlen.pl/app/outRfx/" + fmt.Sprintf("%d", orlenDto.Identity) + "/supplier/status"
	timeObj := time.UnixMilli(orlenDto.EndDate)
	return NewDataDTO("orlen", orlenDto.Name, href, &timeObj, fmt.Sprintf("%v", orlenDto.Identity))
}
