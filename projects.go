package redmine

import (
	"net/http"
	"net/url"
	"strconv"
)

// ProjectStatus defines project status type
type ProjectStatus int64

// ProjectStatus const
const (
	ProjectStatusActive   ProjectStatus = 1
	ProjectStatusClosed   ProjectStatus = 5
	ProjectStatusArchived ProjectStatus = 9
)

type ProjectInclude string

const (
	ProjectIncludeTrackers            ProjectInclude = "trackers"
	ProjectIncludeIssueCategories     ProjectInclude = "issue_categories"
	ProjectIncludeEnabledModules      ProjectInclude = "enabled_modules"       // (since 2.6.0)
	ProjectIncludeTimeEntryActivities ProjectInclude = "time_entry_activities" // (since 3.4.0)
	ProjectIncludeIssueCustomFields   ProjectInclude = "issue_custom_fields"   // (since 4.2.0)
)

/* Get */

// ProjectObject struct used for projects get operations
type ProjectObject struct {
	ID                  int64                  `json:"id"`
	Name                string                 `json:"name"`
	Identifier          string                 `json:"identifier"`
	Description         string                 `json:"description"`
	Homepage            *string                `json:"homepage"` // used only: get single project
	Parent              IDName                 `json:"parent"`
	Status              ProjectStatus          `json:"status"`
	IsPublic            bool                   `json:"is_public"`
	InheritMembers      bool                   `json:"inherit_members"`
	DefaultVersion      *IDName                `json:"default_version"`  // used only: get single project and if set for project
	DefaultAssignee     *IDName                `json:"default_assignee"` // used only: get single project and if set for project
	CustomFields        []CustomFieldGetObject `json:"custom_fields"`
	Trackers            *[]IDName              `json:"trackers"`              // used only: include specified
	IssueCategories     *[]IDName              `json:"issue_categories"`      // used only: include specified
	TimeEntryActivities *[]IDName              `json:"time_entry_activities"` // used only: include specified
	EnabledModules      *[]IDName              `json:"enabled_modules"`       // used only: include specified
	IssueCustomFields   *[]IDName              `json:"issue_custom_fields"`   // used only: include specified
	CreatedOn           string                 `json:"created_on"`
	UpdatedOn           string                 `json:"updated_on"`
}

/* Create */

// ProjectCreate struct used for projects create operations
type ProjectCreate struct {
	Project ProjectCreateObject `json:"project"`
}

type ProjectCreateObject struct {
	Name                string                     `json:"name"`
	Identifier          string                     `json:"identifier"`
	Description         *string                    `json:"description,omitempty"`
	Homepage            *string                    `json:"homepage,omitempty"`
	IsPublic            *bool                      `json:"is_public,omitempty"`
	ParentID            *int64                     `json:"parent_id,omitempty"`
	InheritMembers      *bool                      `json:"inherit_members,omitempty"`
	DefaultAssignedToID *int64                     `json:"default_assigned_to_id,omitempty"`
	DefaultVersionID    *int64                     `json:"default_version_id,omitempty"`
	TrackerIDs          *[]int64                   `json:"tracker_ids,omitempty"`
	EnabledModuleNames  *[]string                  `json:"enabled_module_names,omitempty"`
	IssueCustomFieldIDs *[]int64                   `json:"issue_custom_field_ids,omitempty"`
	CustomFields        *[]CustomFieldUpdateObject `json:"custom_fields,omitempty"`
}

/* Update */

// ProjectUpdate struct used for projects update operations
type ProjectUpdate struct {
	Project ProjectUpdateObject `json:"project"`
}

type ProjectUpdateObject struct {
	Name                *string                    `json:"name,omitempty"`
	Description         *string                    `json:"description,omitempty"`
	Homepage            *string                    `json:"homepage,omitempty"`
	IsPublic            *bool                      `json:"is_public,omitempty"`
	ParentID            *int64                     `json:"parent_id,omitempty"`
	InheritMembers      *bool                      `json:"inherit_members,omitempty"`
	DefaultAssignedToID *int64                     `json:"default_assigned_to_id,omitempty"`
	DefaultVersionID    *int64                     `json:"default_version_id,omitempty"`
	TrackerIDs          *[]int64                   `json:"tracker_ids,omitempty"`
	EnabledModuleNames  *[]string                  `json:"enabled_module_names,omitempty"`
	IssueCustomFieldIDs *[]int64                   `json:"issue_custom_field_ids,omitempty"`
	CustomFields        *[]CustomFieldUpdateObject `json:"custom_fields,omitempty"`
}

/* Requests */

// ProjectAllGetRequest contains data for making request to get all projects satisfying specified filters
type ProjectAllGetRequest struct {
	Includes []ProjectInclude
	Filters  ProjectGetRequestFilters
}

// ProjectMultiGetRequest contains data for making request to get limited projects count satisfying specified filters
type ProjectMultiGetRequest struct {
	Includes []ProjectInclude
	Filters  ProjectGetRequestFilters
	Offset   int64
	Limit    int64
}

// ProjectSingleGetRequest contains data for making request to get specified project
type ProjectSingleGetRequest struct {
	Includes []ProjectInclude
}

// ProjectGetRequestFilters contains data for making projects get request
type ProjectGetRequestFilters struct {
	Status ProjectStatus
}

/* Results */

// ProjectResult stores projects requests processing result
type ProjectResult struct {
	Projects   []ProjectObject `json:"projects"`
	TotalCount int64           `json:"total_count"`
	Offset     int64           `json:"offset"`
	Limit      int64           `json:"limit"`
}

/* Internal types */

type projectSingleResult struct {
	Project ProjectObject `json:"project"`
}

func (p ProjectStatus) String() string {

	status := map[ProjectStatus]string{
		ProjectStatusActive:   "active",
		ProjectStatusClosed:   "closed",
		ProjectStatusArchived: "archived",
	}

	s, b := status[p]
	if b == false {
		return "unknown"
	}

	return s
}

func (pi ProjectInclude) String() string {
	return string(pi)
}

// ProjectAllGet gets info for all projects
//
// see: https://www.redmine.org/projects/redmine/wiki/Rest_Projects#Listing-projects
func (r *Context) ProjectAllGet(request ProjectAllGetRequest) (ProjectResult, StatusCode, error) {

	var (
		projects ProjectResult
		offset   int64
		status   StatusCode
	)

	m := ProjectMultiGetRequest{
		Filters:  request.Filters,
		Includes: request.Includes,
		Limit:    limitDefault,
	}

	for {

		m.Offset = offset

		p, s, err := r.ProjectMultiGet(m)
		if err != nil {
			return projects, s, err
		}

		status = s

		projects.Projects = append(projects.Projects, p.Projects...)

		if offset+p.Limit >= p.TotalCount {
			projects.TotalCount = p.TotalCount
			projects.Limit = p.TotalCount

			break
		}

		offset += p.Limit
	}

	return projects, status, nil
}

// ProjectMultiGet gets info for multiple projects
//
// see: https://www.redmine.org/projects/redmine/wiki/Rest_Projects#Listing-projects
func (r *Context) ProjectMultiGet(request ProjectMultiGetRequest) (ProjectResult, StatusCode, error) {

	var p ProjectResult

	status := ProjectStatusActive
	if request.Filters.Status != 0 {
		status = request.Filters.Status
	}

	urlParams := url.Values{}
	urlParams.Add("offset", strconv.FormatInt(request.Offset, 10))
	urlParams.Add("limit", strconv.FormatInt(request.Limit, 10))
	urlParams.Add("status", strconv.FormatInt(int64(status), 10))

	// Preparing includes
	urlIncludes(
		&urlParams,
		func() []string {
			var is []string
			for _, i := range request.Includes {
				is = append(is, i.String())
			}
			return is
		}(),
	)

	ur := url.URL{
		Path:     "/projects.json",
		RawQuery: urlParams.Encode(),
	}

	s, err := r.Get(&p, ur, http.StatusOK)

	return p, s, err
}

// ProjectSingleGet gets single project info with specified ID
//
// see: https://www.redmine.org/projects/redmine/wiki/Rest_Projects#Showing-a-project
func (r *Context) ProjectSingleGet(id string, request ProjectSingleGetRequest) (ProjectObject, StatusCode, error) {

	var p projectSingleResult

	urlParams := url.Values{}

	// Preparing includes
	urlIncludes(
		&urlParams,
		func() []string {
			var is []string
			for _, i := range request.Includes {
				is = append(is, i.String())
			}
			return is
		}(),
	)

	ur := url.URL{
		Path:     "/projects/" + id + ".json",
		RawQuery: urlParams.Encode(),
	}

	status, err := r.Get(&p, ur, http.StatusOK)

	return p.Project, status, err
}

// ProjectCreate creates new project
//
// see: https://www.redmine.org/projects/redmine/wiki/Rest_Projects#Creating-a-project
func (r *Context) ProjectCreate(project ProjectCreate) (ProjectObject, StatusCode, error) {

	var p projectSingleResult

	ur := url.URL{
		Path: "/projects.json",
	}

	status, err := r.Post(project, &p, ur, http.StatusCreated)

	return p.Project, status, err
}

// ProjectUpdate updates project with specified ID
//
// see: https://www.redmine.org/projects/redmine/wiki/Rest_Projects#Updating-a-project
func (r *Context) ProjectUpdate(id string, project ProjectUpdate) (StatusCode, error) {

	ur := url.URL{
		Path: "/projects/" + id + ".json",
	}

	status, err := r.Put(project, nil, ur, http.StatusNoContent)

	return status, err
}

// ProjectArchive archives a project (available since Redmine 5.0)
//
// see: https://www.redmine.org/projects/redmine/wiki/Rest_Projects#Archiving-a-project
func (r *Context) ProjectArchive(id string) (StatusCode, error) {

	ur := url.URL{
		Path: "/projects/" + id + "/archive.json",
	}

	status, err := r.Put(nil, nil, ur, http.StatusNoContent)

	return status, err
}

// ProjectUnarchive unarchives a project (available since Redmine 5.0)
//
// see: https://www.redmine.org/projects/redmine/wiki/Rest_Projects#Unarchiving-a-project
func (r *Context) ProjectUnarchive(id string) (StatusCode, error) {

	ur := url.URL{
		Path: "/projects/" + id + "/unarchive.json",
	}

	status, err := r.Put(nil, nil, ur, http.StatusNoContent)

	return status, err
}

// ProjectDelete deletes project with specified ID
//
// see: https://www.redmine.org/projects/redmine/wiki/Rest_Projects#Deleting-a-project
func (r *Context) ProjectDelete(id string) (StatusCode, error) {

	ur := url.URL{
		Path: "/projects/" + id + ".json",
	}

	status, err := r.Del(nil, nil, ur, http.StatusNoContent)

	return status, err
}
