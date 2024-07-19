package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AnyResult interface {
}

type ResultOrErr struct {
	Err    HttpErr
	Result AnyResult
}

func (rore *ResultOrErr) HasError() (bool, HttpErr) {
	return (rore.Err != nil), rore.Err
}

type Page struct {
	Idx  int    // index of the page, index for pages start at 1
	HRef string // link to the actual page, hence helps in ranging over pagination i
}

type PaginationResult struct {
	BlogList   ListOfBlogs
	TotalPages []Page
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
