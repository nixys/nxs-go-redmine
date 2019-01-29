package redmine

import (
	"net/url"
	"strconv"
)

// UserStatus const
const (
	UserStatusActive     = 1
	UserStatusRegistered = 2
	UserStatusLocked     = 3
)

// UserNotification const
const (
	UserNotificationAll          = 1
	UserNotificationSelected     = 2
	UserNotificationOnlyMyEvents = 3
	UserNotificationOnlyAssigned = 4
	UserNotificationOnlyOwner    = 5
	UserNotificationOnlyNone     = 6
)

// UserStatus names in Redmine
var UserStatus = map[int]string{
	UserStatusActive:     "active",
	UserStatusRegistered: "registered",
	UserStatusLocked:     "locked",
}

// UserNotification names in Redmine
var UserNotification = map[int]string{
	UserNotificationAll:          "all",
	UserNotificationSelected:     "selected",
	UserNotificationOnlyMyEvents: "only_my_events",
	UserNotificationOnlyAssigned: "only_assigned",
	UserNotificationOnlyOwner:    "only_owner",
	UserNotificationOnlyNone:     "none",
}

/* Get */

// UserObject struct used for users get operations
type UserObject struct {
	ID           int                    `json:"id"`
	Login        string                 `json:"login"`
	FirstName    string                 `json:"firstname"`
	LastName     string                 `json:"lastname"`
	Mail         string                 `json:"mail"`
	CreatedOn    string                 `json:"created_on"`
	LastLoginOn  string                 `json:"last_login_on"`
	APIKey       string                 `json:"api_key"` // used only: get single user
	Status       int                    `json:"status"`  // used only: get single user
	CustomFields []CustomFieldGetObject `json:"custom_fields"`
	Groups       []IDName               `json:"groups"`      // used only: get single user
	Memberships  []UserMembershipObject `json:"memberships"` // used only: get single user
}

// UserMembershipObject struct used for users get operations
type UserMembershipObject struct {
	ID      int      `json:"id"`
	Project IDName   `json:"project"`
	Roles   []IDName `json:"roles"`
}

/* Create */

// UserCreateObject struct used for users create operations
type UserCreateObject struct {
	Login            string                    `json:"login"`
	FirstName        string                    `json:"firstname"`
	LastName         string                    `json:"lastname"`
	Mail             string                    `json:"mail"`
	Password         string                    `json:"password,omitempty"`
	AuthSourceID     int                       `json:"auth_source_id,omitempty"`
	MailNotification string                    `json:"mail_notification,omitempty"`
	MustChangePasswd bool                      `json:"must_change_passwd,omitempty"`
	GeneratePassword bool                      `json:"generate_password,omitempty"`
	SendInformation  bool                      `json:"send_information,omitempty"`
	CustomFields     []CustomFieldUpdateObject `json:"custom_fields,omitempty"`
}

/* Update */

// UserUpdateObject struct used for users update operations
type UserUpdateObject struct {
	Login            string                    `json:"login,omitempty"`
	FirstName        string                    `json:"firstname,omitempty"`
	LastName         string                    `json:"lastname,omitempty"`
	Mail             string                    `json:"mail,omitempty"`
	Password         string                    `json:"password,omitempty"`
	AuthSourceID     int                       `json:"auth_source_id,omitempty"`
	MailNotification string                    `json:"mail_notification,omitempty"`
	MustChangePasswd bool                      `json:"must_change_passwd,omitempty"`
	GeneratePassword bool                      `json:"generate_password,omitempty"`
	SendInformation  bool                      `json:"send_information,omitempty"`
	CustomFields     []CustomFieldUpdateObject `json:"custom_fields,omitempty"`
}

/* Requests */

// UserAllGetRequest contains data for making request to get all users satisfying specified filters
type UserAllGetRequest struct {
	Filters UserGetRequestFilters
}

// UserMultiGetRequest contains data for making request to get limited users count satisfying specified filters
type UserMultiGetRequest struct {
	Filters UserGetRequestFilters
	Offset  int
	Limit   int
}

// UserGetRequestFilters contains data for making users get request
type UserGetRequestFilters struct {
	Status  int
	Name    string
	GroupID int
}

/* Results */

// UserResult stores users requests processing result
type UserResult struct {
	Users      []UserObject `json:"users"`
	TotalCount int          `json:"total_count"`
	Offset     int          `json:"offset"`
	Limit      int          `json:"limit"`
}

/* Internal types */

type userSingleResult struct {
	User UserObject `json:"user"`
}

type userCreate struct {
	User UserCreateObject `json:"user"`
}

type userUpdate struct {
	User UserUpdateObject `json:"user"`
}

// UserAllGet gets info for all users satisfying specified filters
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Users#GET
//
// * If `statusFilter` == 0 default users status filter will be used (show active users only)
// * Use `groupIDFilter` == 0 to disable this filter
func (r *Context) UserAllGet(request UserAllGetRequest) (UserResult, int, error) {

	var (
		users          UserResult
		offset, status int
	)

	m := UserMultiGetRequest{
		Filters: request.Filters,
		Limit:   limitDefault,
	}

	for {

		m.Offset = offset

		u, s, err := r.UserMultiGet(m)
		if err != nil {
			return users, s, err
		}

		status = s

		users.Users = append(users.Users, u.Users...)

		if offset+u.Limit >= u.TotalCount {
			users.TotalCount = u.TotalCount
			users.Limit = u.TotalCount

			break
		}

		offset += u.Limit
	}

	return users, status, nil
}

// UserMultiGet gets info for multiple users satisfying specified filters
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Users#GET
//
// * If `statusFilter` == 0 default users status filter will be used (show active users only)
// * Use `groupIDFilter` == 0 to disable this filter
func (r *Context) UserMultiGet(request UserMultiGetRequest) (UserResult, int, error) {

	var u UserResult

	urlParams := url.Values{}
	urlParams.Add("offset", strconv.Itoa(request.Offset))
	urlParams.Add("limit", strconv.Itoa(request.Limit))

	// Preparing filters
	userURLFilters(&urlParams, request.Filters)

	ur := url.URL{
		Path:     "/users.json",
		RawQuery: urlParams.Encode(),
	}

	s, err := r.get(&u, ur, 200)

	return u, s, err
}

// UserSingleGet gets single user info by specific ID
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Users#GET-2
//
// Available includes:
// * groups
// * memberships
func (r *Context) UserSingleGet(id int, includes []string) (UserObject, int, error) {

	var u userSingleResult

	urlParams := url.Values{}

	// Preparing includes
	urlIncludes(&urlParams, includes)

	ur := url.URL{
		Path:     "/users/" + strconv.Itoa(id) + ".json",
		RawQuery: urlParams.Encode(),
	}

	status, err := r.get(&u, ur, 200)

	return u.User, status, err
}

// UserCreate creates new user
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Users#POST
func (r *Context) UserCreate(user UserCreateObject) (UserObject, int, error) {

	var u userSingleResult

	ur := url.URL{
		Path: "/users.json",
	}

	status, err := r.post(userCreate{User: user}, &u, ur, 201)

	return u.User, status, err
}

// UserUpdate updates user with specified ID
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Users#PUT
func (r *Context) UserUpdate(id int, user UserUpdateObject) (int, error) {

	ur := url.URL{
		Path: "/users/" + strconv.Itoa(id) + ".json",
	}

	status, err := r.put(userUpdate{User: user}, nil, ur, 200)

	return status, err
}

// UserDelete deletes user with specified ID
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Users#DELETE
func (r *Context) UserDelete(id int) (int, error) {

	ur := url.URL{
		Path: "/users/" + strconv.Itoa(id) + ".json",
	}

	status, err := r.del(nil, nil, ur, 200)

	return status, err
}

func userURLFilters(urlParams *url.Values, filters UserGetRequestFilters) {

	if filters.Status == 0 {
		filters.Status = 1
	}

	urlParams.Add("status", strconv.Itoa(filters.Status))

	if len(filters.Name) > 0 {
		urlParams.Add("name", filters.Name)
	}

	if filters.GroupID > 0 {
		urlParams.Add("group_id", strconv.Itoa(filters.GroupID))
	}
}
