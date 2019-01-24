package redmine

/* Get */

// IssueStatusObject struct used for issue_statuses get operations
type IssueStatusObject struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	IsDefault bool   `json:"is_default"`
	IsClosed  bool   `json:"is_closed"`
}

/* Internal types */

type issueStatusMultiResult struct {
	IssueStatuses []IssueStatusObject `json:"issue_statuses"`
}

// IssueStatusMultiGet gets multiple issue statuses
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_IssueStatuses#GET
func (r *Context) IssueStatusMultiGet() ([]IssueStatusObject, int, error) {

	var i issueStatusMultiResult

	uri := "/issue_statuses.json"

	status, err := r.get(&i, uri, 200)

	return i.IssueStatuses, status, err
}
