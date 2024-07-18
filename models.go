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

type PaginationResult struct {
	BlogList   ListOfBlogs
	TotalPages int
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
