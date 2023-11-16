package redmine

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

/* Get */

// IssueObject struct used for issues get operations
type IssueObject struct {
	ID             int64                  `json:"id"`
	Project        IDName                 `json:"project"`
	Tracker        IDName                 `json:"tracker"`
	Status         IDName                 `json:"status"`
	Priority       IDName                 `json:"priority"`
	Author         IDName                 `json:"author"`
	AssignedTo     IDName                 `json:"assigned_to"`
	Category       IDName                 `json:"category"`
	FixedVersion   IDName                 `json:"fixed_version"`
	Parent         IssueParentObject      `json:"parent"`
	Subject        string                 `json:"subject"`
	Description    string                 `json:"description"`
	StartDate      string                 `json:"start_date"`
	DueDate        string                 `json:"due_date"`
	DoneRatio      int64                  `json:"done_ratio"`
	IsPrivate      int64                  `json:"is_private"`
	EstimatedHours float64                `json:"estimated_hours"`
	SpentHours     float64                `json:"spent_hours"` // used only: get single issue
	CustomFields   []CustomFieldGetObject `json:"custom_fields"`
	CreatedOn      string                 `json:"created_on"`
	UpdatedOn      string                 `json:"updated_on"`
	ClosedOn       string                 `json:"closed_on"`
	Children       []IssueChildrenObject  `json:"children"`
	Attachments    []AttachmentObject     `json:"attachments"` // used only: get single issue
	Relations      []IssueRelationObject  `json:"relations"`
	Changesets     []IssueChangesetObject `json:"changesets"` // used only: get single issue
	Journals       []IssueJournalObject   `json:"journals"`   // used only: get single issue
	Watchers       []IDName               `json:"watchers"`   // used only: get single issue
}

// IssueParentObject struct used for issues get operations
type IssueParentObject struct {
	ID int64 `json:"id"`
}

// IssueChildrenObject struct used for issues get operations
type IssueChildrenObject struct {
	ID       int64                 `json:"id"`
	Tracker  IDName                `json:"tracker"`
	Subject  string                `json:"subject"`
	Children []IssueChildrenObject `json:"children"`
}

// IssueChangesetObject struct used for issues get operations
type IssueChangesetObject struct {
	Revision    string `json:"revision"`
	User        IDName `json:"user"`
	Comments    string `json:"comments"`
	CommittedOn string `json:"committed_on"`
}

// IssueRelationObject struct used for issues get operations
type IssueRelationObject struct {
	ID           int64  `json:"id"`
	IssueID      int64  `json:"issue_id"`
	IssueToID    int64  `json:"issue_to_id"`
	RelationType string `json:"relation_type"`
	Delay        int64  `json:"delay"`
}

// IssueJournalObject struct used for issues get operations
type IssueJournalObject struct {
	ID           int64                      `json:"id"`
	User         IDName                     `json:"user"`
	Notes        string                     `json:"notes"`
	CreatedOn    string                     `json:"created_on"`
	PrivateNotes bool                       `json:"private_notes"`
	Details      []IssueJournalDetailObject `json:"details"`
}

// IssueJournalDetailObject struct used for issues get operations
type IssueJournalDetailObject struct {
	Property string `json:"property"`
	Name     string `json:"name"`
	OldValue string `json:"old_value"`
	NewValue string `json:"new_value"`
}

/* Create */

// IssueCreate struct used for issues create operations
type IssueCreate struct {
	Issue IssueCreateObject `json:"issue"`
}

type IssueCreateObject struct {
	ProjectID      int64                     `json:"project_id"`
	TrackerID      int64                     `json:"tracker_id,omitempty"`
	StatusID       int64                     `json:"status_id,omitempty"`
	PriorityID     int64                     `json:"priority_id,omitempty"`
	Subject        string                    `json:"subject"`
	Description    string                    `json:"description,omitempty"`
	StartDate      string                    `json:"start_date,omitempty"`
	DueDate        string                    `json:"due_date,omitempty"`
	CategoryID     int64                     `json:"category_id,omitempty"`
	FixedVersionID int64                     `json:"fixed_version_id,omitempty"`
	AssignedToID   int64                     `json:"assigned_to_id,omitempty"`
	ParentIssueID  int64                     `json:"parent_issue_id,omitempty"`
	WatcherUserIDs []int64                   `json:"watcher_user_ids,omitempty"`
	IsPrivate      bool                      `json:"is_private,omitempty"`
	EstimatedHours float64                   `json:"estimated_hours,omitempty"`
	CustomFields   []CustomFieldUpdateObject `json:"custom_fields,omitempty"`
	Uploads        []AttachmentUploadObject  `json:"uploads,omitempty"`
}

/* Update */

// IssueUpdate struct used for issues update operations
type IssueUpdate struct {
	Issue IssueUpdateObject `json:"issue"`
}

type IssueUpdateObject struct {
	ProjectID      int64                     `json:"project_id,omitempty"`
	TrackerID      int64                     `json:"tracker_id,omitempty"`
	StatusID       int64                     `json:"status_id,omitempty"`
	PriorityID     int64                     `json:"priority_id,omitempty"`
	Subject        string                    `json:"subject,omitempty"`
	Description    string                    `json:"description,omitempty"`
	StartDate      *string                   `json:"start_date,omitempty"`
	DueDate        *string                   `json:"due_date,omitempty"`
	CategoryID     int64                     `json:"category_id,omitempty"`
	FixedVersionID int64                     `json:"fixed_version_id,omitempty"`
	AssignedToID   int64                     `json:"assigned_to_id,omitempty"`
	ParentIssueID  int64                     `json:"parent_issue_id,omitempty"`
	IsPrivate      bool                      `json:"is_private,omitempty"`
	EstimatedHours float64                   `json:"estimated_hours,omitempty"`
	CustomFields   []CustomFieldUpdateObject `json:"custom_fields,omitempty"`
	Uploads        []AttachmentUploadObject  `json:"uploads,omitempty"`
	Notes          string                    `json:"notes,omitempty"`
	PrivateNotes   bool                      `json:"private_notes,omitempty"`
}

/* Requests */

// IssueAllGetRequest contains data for making request to get all issues satisfying specified filters
type IssueAllGetRequest struct {
	Includes []string
	Filters  IssueGetRequestFilters
}

// IssueMultiGetRequest contains data for making request to get limited issues count satisfying specified filters
type IssueMultiGetRequest struct {
	Includes []string
	Filters  IssueGetRequestFilters
	Offset   int64
	Limit    int64
}

// IssueSingleGetRequest contains data for making request to get specified issue
type IssueSingleGetRequest struct {
	Includes []string
}

// IssueGetRequestFilters contains data for making issues get request
type IssueGetRequestFilters struct {
	Fields map[string][]string
	Cf     []IssueGetRequestFiltersCf
}

// IssueGetRequestFiltersCf contains data for making issues get request
type IssueGetRequestFiltersCf struct {
	ID    int64
	Value string
}

/* Results */

// IssueResult stores issues requests processing result
type IssueResult struct {
	Issues     []IssueObject `json:"issues"`
	TotalCount int64         `json:"total_count"`
	Offset     int64         `json:"offset"`
	Limit      int64         `json:"limit"`
}

/* Internal types */

type issueSingleResult struct {
	Issue IssueObject `json:"issue"`
}

type issueWatcherAdd struct {
	UserID int64 `json:"user_id"`
}

// IssuesAllGet gets info for all issues satisfying specified filters
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Issues#Listing-issues
//
// Available includes:
// * attachments - Since 3.4.0
// * relations
// * journals
// * children
func (r *Context) IssuesAllGet(request IssueAllGetRequest) (IssueResult, StatusCode, error) {

	var (
		issues IssueResult
		offset int64
		status StatusCode
	)

	m := IssueMultiGetRequest{
		Filters:  request.Filters,
		Includes: request.Includes,
		Limit:    limitDefault,
	}

	for {

		m.Offset = offset

		i, s, err := r.IssuesMultiGet(m)
		if err != nil {
			return issues, s, err
		}

		status = s

		issues.Issues = append(issues.Issues, i.Issues...)

		if offset+i.Limit >= i.TotalCount {
			issues.TotalCount = i.TotalCount
			issues.Limit = i.TotalCount

			break
		}

		offset += i.Limit
	}

	return issues, status, nil
}

// IssuesMultiGet gets info for multiple issues satisfying specified filters
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Issues#Listing-issues
//
// Available includes:
// * attachments - Since 3.4.0
// * relations
// * journals
// * children
func (r *Context) IssuesMultiGet(request IssueMultiGetRequest) (IssueResult, StatusCode, error) {

	var i IssueResult

	urlParams := url.Values{}
	urlParams.Add("offset", strconv.FormatInt(request.Offset, 10))
	urlParams.Add("limit", strconv.FormatInt(request.Limit, 10))

	// Preparing includes
	urlIncludes(&urlParams, request.Includes)

	// Preparing filters
	issueURLFilters(&urlParams, request.Filters)

	ur := url.URL{
		Path:     "/issues.json",
		RawQuery: urlParams.Encode(),
	}

	s, err := r.Get(&i, ur, http.StatusOK)

	return i, s, err
}

// IssueSingleGet gets single issue info
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Issues#Showing-an-issue
//
// Available includes:
// * children
// * attachments
// * relations
// * changesets
// * journals
// * watchers - Since 2.3.0
func (r *Context) IssueSingleGet(id int64, request IssueSingleGetRequest) (IssueObject, StatusCode, error) {

	var i issueSingleResult

	urlParams := url.Values{}

	// Preparing includes
	urlIncludes(&urlParams, request.Includes)

	ur := url.URL{
		Path:     "/issues/" + strconv.FormatInt(id, 10) + ".json",
		RawQuery: urlParams.Encode(),
	}

	status, err := r.Get(&i, ur, http.StatusOK)

	return i.Issue, status, err
}

// IssueCreate creates new issue
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Issues#Creating-an-issue
func (r *Context) IssueCreate(issue IssueCreate) (IssueObject, StatusCode, error) {

	var i issueSingleResult

	ur := url.URL{
		Path: "/issues.json",
	}

	status, err := r.Post(issue, &i, ur, http.StatusCreated)

	return i.Issue, status, err
}

// IssueUpdate updates issue with specified ID
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Projects#Updating-a-project
func (r *Context) IssueUpdate(id int64, issue IssueUpdate) (StatusCode, error) {

	ur := url.URL{
		Path: "/issues/" + strconv.FormatInt(id, 10) + ".json",
	}

	status, err := r.Put(issue, nil, ur, http.StatusNoContent)

	return status, err
}

// IssueDelete deletes issue with specified ID
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Issues#Deleting-an-issue
func (r *Context) IssueDelete(id int64) (StatusCode, error) {

	ur := url.URL{
		Path: "/issues/" + strconv.FormatInt(id, 10) + ".json",
	}

	status, err := r.Del(nil, nil, ur, http.StatusNoContent)

	return status, err
}

// IssueWatcherAdd adds watcher into issue with specified ID
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Issues#Adding-a-watcher
func (r *Context) IssueWatcherAdd(id int64, userID int64) (StatusCode, error) {

	ur := url.URL{
		Path: "/issues/" + strconv.FormatInt(id, 10) + "/watchers.json",
	}

	status, err := r.Post(issueWatcherAdd{
		UserID: userID,
	}, nil, ur, http.StatusNoContent)

	return status, err
}

// IssueWatcherDelete deletes watcher from issue with specified ID
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Issues#Removing-a-watcher
func (r *Context) IssueWatcherDelete(id int64, userID int64) (StatusCode, error) {

	ur := url.URL{
		Path: "/issues/" + strconv.FormatInt(id, 10) + "/watchers/" + strconv.FormatInt(userID, 10) + ".json",
	}

	status, err := r.Del(nil, nil, ur, http.StatusNoContent)

	return status, err
}

func issueURLFilters(urlParams *url.Values, filters IssueGetRequestFilters) {

	// Filter fields (e.g. `issue_id`, `tracker_id`, etc)
	for n, s := range filters.Fields {
		urlParams.Add(n, strings.Join(s, ","))
	}

	// Custom fields
	for _, c := range filters.Cf {
		urlParams.Add("cf_"+strconv.FormatInt(c.ID, 10), c.Value)
	}
}
