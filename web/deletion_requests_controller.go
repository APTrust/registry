package web

import (
	"fmt"
	"net/http"

	"github.com/APTrust/registry/forms"
	"github.com/APTrust/registry/pgmodels"
	"github.com/gin-gonic/gin"
)

// DeletionRequestShow shows the deletion request with the specified id.
//
// Note that this shows a read-only view of the request. It does not include
// the Approve/Cancel buttons. This read-only view may be available to users
// who do not have permission to initiate, approve, or cancel deletion requests
// but who still need a read-only view of the requests that have been submitted.
//
// Deletions apply to files and/or intellectual objects. The methods for
// initiating, approving and rejecting deletion requests are in the
// Generic Files Controller (for files) and the Intellectual Objects Controller
// (for objects).
//
// GET /deletions/show/:id
func DeletionRequestShow(c *gin.Context) {
	req := NewRequest(c)
	deletionRequest, err := pgmodels.DeletionRequestByID(req.Auth.ResourceID)
	if AbortIfError(c, err) {
		return
	}
	req.TemplateData["deletionRequest"] = deletionRequest

	if deletionRequest.WorkItemID > 0 {
		req.TemplateData["workItemURL"] = fmt.Sprintf("%s/work_items/show/%d",
			req.BaseURL(),
			deletionRequest.WorkItemID)
	}

	c.HTML(http.StatusOK, "deletions/show.html", req.TemplateData)
}

// DeletionRequestIndex shows list of deletion requests.
// GET /deletions
func DeletionRequestIndex(c *gin.Context) {
	req := NewRequest(c)
	err := drIndexLoadDeletions(req)
	if AbortIfError(c, err) {
		return
	}
	c.HTML(http.StatusOK, "deletions/index.html", req.TemplateData)
}

func drIndexLoadDeletions(req *Request) error {
	filterCollection := req.GetFilterCollection()
	query, err := filterCollection.ToQuery()
	if err != nil {
		return err
	}
	if !req.CurrentUser.IsAdmin() {
		query.Where("institution_id", "=", req.CurrentUser.InstitutionID)
	}
	query.OrderBy("requested_at desc")
	baseURL := req.GinContext.Request.URL.Path + "?" + req.GinContext.Request.URL.RawQuery
	pager, err := NewPager(req.GinContext, baseURL, 20)
	if err != nil {
		return err
	}
	query.Offset(pager.QueryOffset).Limit(pager.PerPage)
	deletions, err := pgmodels.DeletionRequestViewSelect(query)
	if err != nil {
		return err
	}

	totalRecordCount, err := query.Count(&pgmodels.DeletionRequestView{})
	if err != nil {
		return err
	}
	pager.SetCounts(totalRecordCount, len(deletions))

	form, err := forms.NewDeletionRequestFilterForm(filterCollection, req.CurrentUser)

	req.TemplateData["deletions"] = deletions
	req.TemplateData["pager"] = pager
	req.TemplateData["filterForm"] = form

	return err
}
