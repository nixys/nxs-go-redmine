package redmine

import (
	"net/http"
	"net/url"
)

/* Get */

// EnumerationPriorityObject struct used for priorities get operations
type EnumerationPriorityObject struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	IsDefault bool   `json:"is_default"`
}

// EnumerationTimeEntryActivityObject struct used for time entry activities get operations
type EnumerationTimeEntryActivityObject struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	IsDefault bool   `json:"is_default"`
}

// EnumerationDocumentCategoryObject struct used for document categories get operations
type EnumerationDocumentCategoryObject struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	IsDefault bool   `json:"is_default"`
}

/* Internal types */

type enumerationPrioritiesAllResult struct {
	Priorities []EnumerationPriorityObject `json:"issue_priorities"`
}

type enumerationTimeEntryActivitiesAllResult struct {
	TimeEntryActivities []EnumerationTimeEntryActivityObject `json:"time_entry_activities"`
}

type enumerationDocumentCategoriesAllResult struct {
	DocumentCategories []EnumerationDocumentCategoryObject `json:"document_categories"`
}

// EnumerationPrioritiesAllGet gets info for all priority enumerations
//
// see: https://www.redmine.org/projects/redmine/wiki/Rest_Enumerations#GET
func (r *Context) EnumerationPrioritiesAllGet() ([]EnumerationPriorityObject, int, error) {

	var e enumerationPrioritiesAllResult

	ur := url.URL{
		Path: "/enumerations/issue_priorities.json",
	}

	status, err := r.Get(&e, ur, http.StatusOK)

	return e.Priorities, status, err
}

// EnumerationTimeEntryActivitiesAllGet gets info for all time entry activity enumerations
//
// see: https://www.redmine.org/projects/redmine/wiki/Rest_Enumerations#GET-2
func (r *Context) EnumerationTimeEntryActivitiesAllGet() ([]EnumerationTimeEntryActivityObject, int, error) {

	var e enumerationTimeEntryActivitiesAllResult

	ur := url.URL{
		Path: "/enumerations/time_entry_activities.json",
	}

	status, err := r.Get(&e, ur, http.StatusOK)

	return e.TimeEntryActivities, status, err
}

// EnumerationDocumentCategoriesAllGet gets info for all document category enumerations
//
// see: https://www.redmine.org/projects/redmine/wiki/Rest_Enumerations#GET-3
func (r *Context) EnumerationDocumentCategoriesAllGet() ([]EnumerationDocumentCategoryObject, int, error) {

	var e enumerationDocumentCategoriesAllResult

	ur := url.URL{
		Path: "/enumerations/document_categories.json",
	}

	status, err := r.Get(&e, ur, http.StatusOK)

	return e.DocumentCategories, status, err
}
