package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AnyResult interface {
}

// Co2FtPrintParams : this helps determine - or allows us to calculate the co2 footprint
// annual emissions / yeild (fish + veggies)
// the form submitted can bind to this data type
type Co2FtPrintParams struct {
	ElectricKwh   float32 `form:"kwh"`
	FishFeedKgs   float32 `form:"feedkg"`
	PlantYeildKgs float32 `form:"yeildkg"`
	FishYeildKgs  float32 `form:"fishkg"`
	Emissions     float32
	Footprint     float32
}

type ResultOrErr struct {
	Err    HttpErr
	Result AnyResult
}

func (rore *ResultOrErr) HasError() (bool, HttpErr) {
	return (rore.Err != nil), rore.Err
}

type Page struct {
	Idx    int    // index of the page, index for pages start at 1
	HRef   string // link to the actual page, hence helps in ranging over pagination
	IsCurr bool   // flag to indicate that this is the current page
}

type PaginationResult struct {
	BlogList   ListOfBlogs
	TotalPages []Page
}

type PagePayload struct {
	OpnGraph *OpenGraph
	Content  interface{} // page specific content
	Foot     *Footer
}

type Footer struct {
	FBLink     string
	GmailLink  string
	LinkedLink string
	GitLink    string
	Phone      string
	Address    string
	Email      string
}

// OpenGraphPayload : payload that every page needs to carry to be able to posted on social web
// Find these fields being used on head partial template
type OpenGraph struct {
	Title   string
	OgImage string
	OgTitle string
	OgDesc  string
	OgUrl   string
}

// concrete Content under PagePayload
type IndexPagePaylod struct {
	// Any of the index pages - Splash, About Us, Journey, About aquaponics
	// Basically anything that has images and  text paragraphs side by side
	TemplData []MoreInfo
}

type HttpErr interface {
	ToHttpCtx(c *gin.Context)
}

/* Custom errors that are convertible to http context */
type InvalidQueryParam struct {
	Err error
}

func (iqp *InvalidQueryParam) Error() string {
	return fmt.Errorf("One or more query parameters in the request is invalid: %s", iqp.Err.Error()).Error()
}

func (iqp *InvalidQueryParam) ToHttpCtx(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
		"err_msg": iqp.Error(),
	})
}

type InvalidArgument struct {
	Err error
}

func (ia *InvalidArgument) Error() string {
	return fmt.Errorf("One or more arguemnts to operations is invalid : %s", ia.Err.Error()).Error()
}

func (ia *InvalidArgument) ToHttpCtx(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
		"err_msg": ia.Error(),
	})
}
