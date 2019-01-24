package redmine

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

/* Internal types */

type customFieldMultiResult struct {
	CustomFields []CustomFieldObject `json:"custom_fields"`
}

// CustomFieldMultiGet gets multiple custom fields
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_CustomFields#GET
func (r *Context) CustomFieldMultiGet() ([]CustomFieldObject, int, error) {

	var c customFieldMultiResult

	uri := "/custom_fields.json"

	status, err := r.get(&c, uri, 200)

	return c.CustomFields, status, err
}
