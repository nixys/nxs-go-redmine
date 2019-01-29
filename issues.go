package redmine

import (
	"net/url"
	"strconv"
	"strings"
)

/* Get */

// IssueObject struct used for issues get operations
type IssueObject struct {
	ID             int                    `json:"id"`
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
	DoneRatio      int                    `json:"done_ratio"`
	IsPrivate      int                    `json:"is_private"`
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
	ID int `json:"id"`
}

// IssueChildrenObject struct used for issues get operations
type IssueChildrenObject struct {
	ID       int                   `json:"id"`
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
	ID           int    `json:"id"`
	IssueID      int    `json:"issue_id"`
	IssueToID    int    `json:"issue_to_id"`
	RelationType string `json:"relation_type"`
	Delay        int    `json:"delay"`
}

// IssueJournalObject struct used for issues get operations
type IssueJournalObject struct {
	ID        int                        `json:"id"`
	User      IDName                     `json:"user"`
	Notes     string                     `json:"notes"`
	CreatedOn string                     `json:"created_on"`
	Details   []IssueJournalDetailObject `json:"details"`
}

// IssueJournalDetailObject struct used for issues get operations
type IssueJournalDetailObject struct {
	Property string `json:"property"`
	Name     string `json:"name"`
	OldValue string `json:"old_value"`
	NewValue string `json:"new_value"`
}

/* Create */

// IssueCreateObject struct used for issues create operations
type IssueCreateObject struct {
	ProjectID      int                       `json:"project_id"`
	TrackerID      int                       `json:"tracker_id,omitempty"`
	StatusID       int                       `json:"status_id,omitempty"`
	PriorityID     int                       `json:"priority_id,omitempty"`
	Subject        string                    `json:"subject"`
	Description    string                    `json:"description,omitempty"`
	CategoryID     int                       `json:"category_id,omitempty"`
	FixedVersionID int                       `json:"fixed_version_id,omitempty"`
	AssignedToID   int                       `json:"assigned_to_id,omitempty"`
	ParentIssueID  int                       `json:"parent_issue_id,omitempty"`
	WatcherUserIDs []int                     `json:"watcher_user_ids,omitempty"`
	IsPrivate      bool                      `json:"is_private,omitempty"`
	EstimatedHours float64                   `json:"estimated_hours,omitempty"`
	CustomFields   []CustomFieldUpdateObject `json:"custom_fields,omitempty"`
	Uploads        []AttachmentUploadObject  `json:"uploads,omitempty"`
}

/* Update */

// IssueUpdateObject struct used for issues update operations
type IssueUpdateObject struct {
	ProjectID      int                       `json:"project_id,omitempty"`
	TrackerID      int                       `json:"tracker_id,omitempty"`
	StatusID       int                       `json:"status_id,omitempty"`
	PriorityID     int                       `json:"priority_id,omitempty"`
	Subject        string                    `json:"subject,omitempty"`
	Description    string                    `json:"description,omitempty"`
	CategoryID     int                       `json:"category_id,omitempty"`
	FixedVersionID int                       `json:"fixed_version_id,omitempty"`
	AssignedToID   int                       `json:"assigned_to_id,omitempty"`
	ParentIssueID  int                       `json:"parent_issue_id,omitempty"`
	IsPrivate      bool                      `json:"is_private,omitempty"`
	EstimatedHours float64                   `json:"estimated_hours,omitempty"`
	CustomFields   []CustomFieldUpdateObject `json:"custom_fields,omitempty"`
	Uploads        []AttachmentUploadObject  `json:"uploads,omitempty"`
	Notes          string                    `json:"notes,omitempty"`
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
	Offset   int
	Limit    int
}

// IssueGetRequestFilters contains data for making issues get request
type IssueGetRequestFilters struct {
	Fields map[string][]string
	Cf     []IssueGetRequestFiltersCf
}

// IssueGetRequestFiltersCf contains data for making issues get request
type IssueGetRequestFiltersCf struct {
	ID    int
	Value string
}

/* Results */

// IssueResult stores issues requests processing result
type IssueResult struct {
	Issues     []IssueObject `json:"issues"`
	TotalCount int           `json:"total_count"`
	Offset     int           `json:"offset"`
	Limit      int           `json:"limit"`
}

/* Internal types */

type issueSingleResult struct {
	Issue IssueObject `json:"issue"`
}

type issueCreate struct {
	Issue IssueCreateObject `json:"issue"`
}

type issueUpdate struct {
	Issue IssueUpdateObject `json:"issue"`
}

type issueWatcherAdd struct {
	UserID int `json:"user_id"`
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
func (r *Context) IssuesAllGet(request IssueAllGetRequest) (IssueResult, int, error) {

	var (
		issues         IssueResult
		offset, status int
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
func (r *Context) IssuesMultiGet(request IssueMultiGetRequest) (IssueResult, int, error) {

	var i IssueResult

	urlParams := url.Values{}
	urlParams.Add("offset", strconv.Itoa(request.Offset))
	urlParams.Add("limit", strconv.Itoa(request.Limit))

	// Preparing includes
	urlIncludes(&urlParams, request.Includes)

	// Preparing filters
	issueURLFilters(&urlParams, request.Filters)

	ur := url.URL{
		Path:     "/issues.json",
		RawQuery: urlParams.Encode(),
	}

	s, err := r.get(&i, ur, 200)

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
func (r *Context) IssueSingleGet(id int, includes []string) (IssueObject, int, error) {

	var i issueSingleResult

	urlParams := url.Values{}

	// Preparing includes
	urlIncludes(&urlParams, includes)

	ur := url.URL{
		Path:     "/issues/" + strconv.Itoa(id) + ".json",
		RawQuery: urlParams.Encode(),
	}

	status, err := r.get(&i, ur, 200)

	return i.Issue, status, err
}

// IssueCreate creates new issue
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Issues#Creating-an-issue
func (r *Context) IssueCreate(issue IssueCreateObject) (IssueObject, int, error) {

	var i issueSingleResult

	ur := url.URL{
		Path: "/issues.json",
	}

	status, err := r.post(issueCreate{Issue: issue}, &i, ur, 201)

	return i.Issue, status, err
}

// IssueUpdate updates issue with specified ID
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Projects#Updating-a-project
func (r *Context) IssueUpdate(id int, issue IssueUpdateObject) (int, error) {

	ur := url.URL{
		Path: "/issues/" + strconv.Itoa(id) + ".json",
	}

	status, err := r.put(issueUpdate{Issue: issue}, nil, ur, 200)

	return status, err
}

// IssueDelete deletes issue with specified ID
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Issues#Deleting-an-issue
func (r *Context) IssueDelete(id int) (int, error) {

	ur := url.URL{
		Path: "/issues/" + strconv.Itoa(id) + ".json",
	}

	status, err := r.del(nil, nil, ur, 200)

	return status, err
}

// IssueWatcherAdd adds watcher into issue with specified ID
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Issues#Adding-a-watcher
func (r *Context) IssueWatcherAdd(id int, userID int) (int, error) {

	ur := url.URL{
		Path: "/issues/" + strconv.Itoa(id) + "/watchers.json",
	}

	status, err := r.post(issueWatcherAdd{
		UserID: userID,
	}, nil, ur, 200)

	return status, err
}

// IssueWatcherDelete deletes watcher from issue with specified ID
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Issues#Removing-a-watcher
func (r *Context) IssueWatcherDelete(id int, userID int) (int, error) {

	ur := url.URL{
		Path: "/issues/" + strconv.Itoa(id) + "/watchers/" + strconv.Itoa(userID) + ".json",
	}

	status, err := r.del(nil, nil, ur, 200)

	return status, err
}

func issueURLFilters(urlParams *url.Values, filters IssueGetRequestFilters) {

	// Filter fields (e.g. `issue_id`, `tracker_id`, etc)
	for n, s := range filters.Fields {
		urlParams.Add(n, strings.Join(s, ","))
	}

	// Custom fields
	for _, c := range filters.Cf {
		urlParams.Add("cf_"+strconv.Itoa(c.ID), c.Value)
	}
}
