package redmine

import (
	"net/http"
	"net/url"
	"strconv"
)

// ProjectStatus defines project status type
type ProjectStatus int

// ProjectStatus const
const (
	ProjectStatusActive   ProjectStatus = 1
	ProjectStatusClosed   ProjectStatus = 5
	ProjectStatusArchived ProjectStatus = 9
)

/* Get */

// ProjectObject struct used for projects get operations
type ProjectObject struct {
	ID              int                    `json:"id"`
	Name            string                 `json:"name"`
	Identifier      string                 `json:"identifier"`
	Description     string                 `json:"description"`
	Homepage        string                 `json:"homepage"` // used only: get single project
	Parent          IDName                 `json:"parent"`
	Status          ProjectStatus          `json:"status"`
	CustomFields    []CustomFieldGetObject `json:"custom_fields"`
	Trackers        []IDName               `json:"trackers"`
	IssueCategories []IDName               `json:"issue_categories"`
	EnabledModules  []IDName               `json:"enabled_modules"`
	CreatedOn       string                 `json:"created_on"`
	UpdatedOn       string                 `json:"updated_on"`
}

/* Create */

// ProjectCreateObject struct used for projects create operations
type ProjectCreateObject struct {
	Name                string                    `json:"name"`
	Identifier          string                    `json:"identifier"`
	Description         string                    `json:"description,omitempty"`
	Homepage            string                    `json:"homepage,omitempty"`
	IsPublic            bool                      `json:"is_public,omitempty"`
	ParentID            int                       `json:"parent_id,omitempty"`
	InheritMembers      bool                      `json:"inherit_members,omitempty"`
	TrackerIDs          []int                     `json:"tracker_ids,omitempty"`
	EnabledModuleNames  []string                  `json:"enabled_module_names,omitempty"`
	IssueCustomFieldIDs []int                     `json:"issue_custom_field_ids,omitempty"`
	CustomFields        []CustomFieldUpdateObject `json:"custom_fields,omitempty"`
}

/* Update */

// ProjectUpdateObject struct used for projects update operations
type ProjectUpdateObject struct {
	Name                string                    `json:"name,omitempty"`
	Description         string                    `json:"description,omitempty"`
	Homepage            string                    `json:"homepage,omitempty"`
	IsPublic            bool                      `json:"is_public,omitempty"`
	ParentID            int                       `json:"parent_id,omitempty"`
	InheritMembers      bool                      `json:"inherit_members,omitempty"`
	TrackerIDs          []int                     `json:"tracker_ids,omitempty"`
	EnabledModuleNames  []string                  `json:"enabled_module_names,omitempty"`
	IssueCustomFieldIDs []int                     `json:"issue_custom_field_ids,omitempty"`
	CustomFields        []CustomFieldUpdateObject `json:"custom_fields,omitempty"`
}

/* Requests */

// ProjectAllGetRequest contains data for making request to get all projects satisfying specified filters
type ProjectAllGetRequest struct {
	Includes []string
	Filters  ProjectGetRequestFilters
}

// ProjectMultiGetRequest contains data for making request to get limited projects count satisfying specified filters
type ProjectMultiGetRequest struct {
	Includes []string
	Filters  ProjectGetRequestFilters
	Offset   int
	Limit    int
}

// ProjectSingleGetRequest contains data for making request to get specified project
type ProjectSingleGetRequest struct {
	Includes []string
}

// ProjectGetRequestFilters contains data for making projects get request
type ProjectGetRequestFilters struct {
	Status ProjectStatus
}

/* Results */

// ProjectResult stores projects requests processing result
type ProjectResult struct {
	Projects   []ProjectObject `json:"projects"`
	TotalCount int             `json:"total_count"`
	Offset     int             `json:"offset"`
	Limit      int             `json:"limit"`
}

/* Internal types */

type projectSingleResult struct {
	Project ProjectObject `json:"project"`
}

type projectCreate struct {
	Project ProjectCreateObject `json:"project"`
}

type projectUpdate struct {
	Project ProjectUpdateObject `json:"project"`
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

// ProjectAllGet gets info for all projects
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Projects#Listing-projects
//
// Available includes:
// * trackers
// * issue_categories
// * enabled_modules
func (r *Context) ProjectAllGet(request ProjectAllGetRequest) (ProjectResult, int, error) {

	var (
		projects       ProjectResult
		offset, status int
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
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Projects#Listing-projects
//
// Available includes:
// * trackers
// * issue_categories
// * enabled_modules
func (r *Context) ProjectMultiGet(request ProjectMultiGetRequest) (ProjectResult, int, error) {

	var p ProjectResult

	status := ProjectStatusActive
	if request.Filters.Status != 0 {
		status = request.Filters.Status
	}

	urlParams := url.Values{}
	urlParams.Add("offset", strconv.Itoa(request.Offset))
	urlParams.Add("limit", strconv.Itoa(request.Limit))
	urlParams.Add("status", strconv.Itoa(int(status)))

	// Preparing includes
	urlIncludes(&urlParams, request.Includes)

	ur := url.URL{
		Path:     "/projects.json",
		RawQuery: urlParams.Encode(),
	}

	s, err := r.get(&p, ur, http.StatusOK)

	return p, s, err
}

// ProjectSingleGet gets single project info with specified ID
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Projects#Showing-a-project
//
// Available includes:
// * trackers
// * issue_categories
// * enabled_modules
// * time_entry_activities (since 3.4.0)
func (r *Context) ProjectSingleGet(id int, request ProjectSingleGetRequest) (ProjectObject, int, error) {

	var p projectSingleResult

	urlParams := url.Values{}

	// Preparing includes
	urlIncludes(&urlParams, request.Includes)

	ur := url.URL{
		Path:     "/projects/" + strconv.Itoa(id) + ".json",
		RawQuery: urlParams.Encode(),
	}

	status, err := r.get(&p, ur, http.StatusOK)

	return p.Project, status, err
}

// ProjectCreate creates new project
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Projects#Creating-a-project
func (r *Context) ProjectCreate(project ProjectCreateObject) (ProjectObject, int, error) {

	var p projectSingleResult

	ur := url.URL{
		Path: "/projects.json",
	}

	status, err := r.post(projectCreate{Project: project}, &p, ur, http.StatusCreated)

	return p.Project, status, err
}

// ProjectUpdate updates project with specified ID
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Projects#Updating-a-project
func (r *Context) ProjectUpdate(id int, project ProjectUpdateObject) (int, error) {

	ur := url.URL{
		Path: "/projects/" + strconv.Itoa(id) + ".json",
	}

	status, err := r.put(projectUpdate{Project: project}, nil, ur, http.StatusNoContent)

	return status, err
}

// ProjectDelete deletes project with specified ID
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Projects#Deleting-a-project
func (r *Context) ProjectDelete(id int) (int, error) {

	ur := url.URL{
		Path: "/projects/" + strconv.Itoa(id) + ".json",
	}

	status, err := r.del(nil, nil, ur, http.StatusNoContent)

	return status, err
}
