package redmine

import (
	"net/http"
	"net/url"
)

/* Get */

// CustomFieldObject struct used for custom fields get operations
type CustomFieldObject struct {
	ID             int                              `json:"id"`
	Name           string                           `json:"name"`
	CustomizedType string                           `json:"customized_type"`
	FieldFormat    string                           `json:"field_format"`
	Regexp         string                           `json:"regexp"`
	MinLength      int                              `json:"min_length"`
	MaxLength      int                              `json:"max_length"`
	IsRequired     bool                             `json:"is_required"`
	IsFilter       bool                             `json:"is_filter"`
	Searchable     bool                             `json:"searchable"`
	Multiple       bool                             `json:"multiple"`
	DefaultValue   string                           `json:"default_value"`
	Visible        bool                             `json:"visible"`
	PossibleValues []CustomFieldPossibleValueObject `json:"possible_values"`
	Trackers       []IDName                         `json:"trackers"`
	Roles          []IDName                         `json:"roles"`
}

// CustomFieldPossibleValueObject struct used for custom fields get operations
type CustomFieldPossibleValueObject struct {
	Value string `json:"value"`
}

// CustomFieldGetObject struct used for custom fields get operations in other methods
type CustomFieldGetObject struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Multiple bool     `json:"multiple"`
	Value    []string `json:"value"`
}

/* Update */

// CustomFieldUpdateObject struct used for custom fields insert and update operations in other methods
type CustomFieldUpdateObject struct {
	ID    int         `json:"id"`
	Value interface{} `json:"value"` // can be a string or strings slice
}

/* Internal types */

type customFieldAllResult struct {
	CustomFields []CustomFieldObject `json:"custom_fields"`
}

// CustomFieldAllGet gets info for all custom fields
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_CustomFields#GET
func (r *Context) CustomFieldAllGet() ([]CustomFieldObject, int, error) {

	var c customFieldAllResult

	ur := url.URL{
		Path: "/custom_fields.json",
	}

	status, err := r.get(&c, ur, http.StatusOK)

	return c.CustomFields, status, err
}
