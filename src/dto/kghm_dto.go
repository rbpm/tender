package dto

import (
	"tender/interfaces/data"
	"tender/tools"
)

const KGHM_TIME_LAYOUT = "2006-01-02T15:04:05Z"

type KghmDTO struct {
	Src           string
	Name          string
	Href          string
	OfferDate     string
	PublishedDate string
	UpdatedDate   string
	Id            string
	IsIT          bool
}

func NewKghmDTO(src string, name string, href string, offerDate string,
	publishedDate string, updatedDate string, id string) *KghmDTO {
	p := KghmDTO{Src: src, Name: name, Href: href, OfferDate: offerDate,
		PublishedDate: publishedDate, UpdatedDate: updatedDate, Id: id, IsIT: data.IsIT(name)}
	return &p
}

func (kghmDto KghmDTO) GetDataDTO() *DataDTO {
	timePtr := tools.ParseDate(KGHM_TIME_LAYOUT, kghmDto.OfferDate)
	return NewDataDTO(kghmDto.Src, kghmDto.Name, kghmDto.Href, timePtr, kghmDto.Id)
}

type KghmAjaxDTO struct {
	Command  string `json:"command"`
	Settings *struct {
		AjaxPageState struct {
			Theme     string `json:"theme"`
			Libraries string `json:"libraries"`
		} `json:"ajaxPageState"`
		Views struct {
			AjaxPath  string `json:"ajax_path"`
			AjaxViews struct {
				ViewsDomId0534474Bbb5E8Add417D226807501A98Bef0682Dc080A6E60D68F3C2Ba52E181 struct {
					ViewName      string      `json:"view_name"`
					ViewDisplayId string      `json:"view_display_id"`
					ViewArgs      string      `json:"view_args"`
					ViewPath      string      `json:"view_path"`
					ViewBasePath  interface{} `json:"view_base_path"`
					ViewDomId     string      `json:"view_dom_id"`
					PagerElement  int         `json:"pager_element"`
				} `json:"views_dom_id:0534474bbb5e8add417d226807501a98bef0682dc080a6e60d68f3c2ba52e181"`
			} `json:"ajaxViews"`
		} `json:"views"`
		AjaxTrustedUrl struct {
			PlPrzetargiNieograniczone bool `json:"/pl/przetargi-nieograniczone"`
		} `json:"ajaxTrustedUrl"`
		BetterExposedFilters struct {
			Datepicker        bool          `json:"datepicker"`
			DatepickerOptions []interface{} `json:"datepicker_options"`
		} `json:"better_exposed_filters"`
		PluralDelimiter string `json:"pluralDelimiter"`
		User            struct {
			Uid             int    `json:"uid"`
			PermissionsHash string `json:"permissionsHash"`
		} `json:"user"`
	} `json:"settings,omitempty"`
	Merge    bool   `json:"merge,omitempty"`
	Selector string `json:"selector,omitempty"`
	Method   string `json:"method,omitempty"`
	Data     string `json:"data,omitempty"`
}
