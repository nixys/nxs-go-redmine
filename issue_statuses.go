package redmine

import (
	"net/http"
	"net/url"
)

/* Get */

// IssueStatusObject struct used for issue_statuses get operations
type IssueStatusObject struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	IsDefault bool   `json:"is_default"`
	IsClosed  bool   `json:"is_closed"`
}

/* Internal types */

type issueStatusAllResult struct {
	IssueStatuses []IssueStatusObject `json:"issue_statuses"`
}

// IssueStatusAllGet gets info for all issue statuses
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_IssueStatuses#GET
func (r *Context) IssueStatusAllGet() ([]IssueStatusObject, int, error) {

	var i issueStatusAllResult

	ur := url.URL{
		Path: "/issue_statuses.json",
	}

	status, err := r.Get(&i, ur, http.StatusOK)

	return i.IssueStatuses, status, err
}
