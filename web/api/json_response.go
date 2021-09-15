package api

import (
	"github.com/APTrust/registry/common"
)

// JsonList provides the structure for an API response
// that contains a list of items.
type JsonList struct {
	// Count is the total number of items in the result set.
	Count int `json:"count"`
	// Next is the URL for the next page of results.
	Next string `json:"next"`
	// Previous is the URL for the previous page of results.
	Previous string `json:"previous"`
	// Results is the list of items on this page of the result set.
	Results interface{} `json:"results"`
}

func NewJsonList(items interface{}, pager *common.Pager) *JsonList {
	return &JsonList{
		Count:    pager.TotalItems,
		Next:     pager.NextLink,
		Previous: pager.PreviousLink,
		Results:  items,
	}
}
