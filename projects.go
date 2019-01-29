package redmine

import (
	"strconv"
	"strings"
)

// ProjectStatus const
const (
	ProjectStatusActive   = 1
	ProjectStatusClosed   = 5
	ProjectStatusArchived = 9
)

// ProjectStatus map
var ProjectStatus = map[int]string{
	ProjectStatusActive:   "active",
	ProjectStatusClosed:   "closed",
	ProjectStatusArchived: "archived",
}

/* Get */

// ProjectObject struct used for projects get operations
type ProjectObject struct {
	ID              int                    `json:"id"`
	Name            string                 `json:"name"`
	Identifier      string                 `json:"identifier"`
	Description     string                 `json:"description"`
	Homepage        string                 `json:"homepage"` // used only: get single project
	Parent          IDName                 `json:"parent"`
	Status          int                    `json:"status"`
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

// ProjectAllGet gets info for all projects
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Projects#Listing-projects
//
// Available includes:
// * trackers
// * issue_categories
// * enabled_modules
func (r *Context) ProjectAllGet(includes []string) (ProjectResult, int, error) {

	var (
		projects       ProjectResult
		offset, status int
	)

	for {

		p, s, err := r.ProjectMultiGet(includes, offset, limitDefault)
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
func (r *Context) ProjectMultiGet(includes []string, offset, limit int) (ProjectResult, int, error) {

	var p ProjectResult
	var i string

	if len(includes) != 0 {
		i = "&include=" + strings.Join(includes, ",")
	}

	uri := "/projects.json?limit=" + strconv.Itoa(limit) + "&offset=" + strconv.Itoa(offset) + i

	s, err := r.get(&p, uri, 200)

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
func (r *Context) ProjectSingleGet(id int, includes []string) (ProjectObject, int, error) {

	var p projectSingleResult
	var i string

	if len(includes) != 0 {
		i = "?include=" + strings.Join(includes, ",")
	}

	uri := "/projects/" + strconv.Itoa(id) + ".json" + i

	status, err := r.get(&p, uri, 200)

	return p.Project, status, err
}

// ProjectCreate creates new project
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Projects#Creating-a-project
func (r *Context) ProjectCreate(project ProjectCreateObject) (ProjectObject, int, error) {

	var p projectSingleResult

	uri := "/projects.json"

	status, err := r.post(projectCreate{Project: project}, &p, uri, 201)

	return p.Project, status, err
}

// ProjectUpdate updates project with specified ID
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Projects#Updating-a-project
func (r *Context) ProjectUpdate(id int, project ProjectUpdateObject) (int, error) {

	uri := "/projects/" + strconv.Itoa(id) + ".json"

	status, err := r.put(projectUpdate{Project: project}, nil, uri, 200)

	return status, err
}

// ProjectDelete deletes project with specified ID
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Projects#Deleting-a-project
func (r *Context) ProjectDelete(id int) (int, error) {

	uri := "/projects/" + strconv.Itoa(id) + ".json"

	status, err := r.del(nil, nil, uri, 200)

	return status, err
}
