package api

import (
	"encoding/json"

	"github.com/APTrust/registry/common"
	"github.com/APTrust/registry/pgmodels"
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

// NewJsonList creates a new json list response structure.
func NewJsonList(items interface{}, pager *common.Pager) *JsonList {
	return &JsonList{
		Count:    pager.TotalItems,
		Next:     pager.NextLink,
		Previous: pager.PreviousLink,
		Results:  items,
	}
}

// NewListFromJson converts a json string to a JsonList object.
// This is used primarily in API testing.
func NewListFromJson(jsonStr string) (*JsonList, error) {
	jsonList := &JsonList{}
	err := json.Unmarshal([]byte(jsonStr), jsonList)
	return jsonList, err
}

// AlertViewList is used into testing to covert a generic
// JsonList into a typed list that we can test with assertions.
type AlertViewList struct {
	Count    int                   `json:"count"`
	Next     string                `json:"next"`
	Previous string                `json:"previous"`
	Results  []*pgmodels.AlertView `json:"results"`
}

// DeletionRequestViewList is used into testing to covert a generic
// JsonList into a typed list that we can test with assertions.
type DeletionRequestViewList struct {
	Count    int                             `json:"count"`
	Next     string                          `json:"next"`
	Previous string                          `json:"previous"`
	Results  []*pgmodels.DeletionRequestView `json:"results"`
}

// GenericFileList is used into testing to covert a generic
// JsonList into a typed list that we can test with assertions.
type GenericFileList struct {
	Count    int                     `json:"count"`
	Next     string                  `json:"next"`
	Previous string                  `json:"previous"`
	Results  []*pgmodels.GenericFile `json:"results"`
}
