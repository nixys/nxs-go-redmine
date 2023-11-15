package redmine

import (
	"net/http"
	"net/url"
	"strconv"
)

// UserStatus defines user status type
type UserStatus int

// UserNotification defines user notification type
type UserNotification string

// UserStatus const
const (
	UserStatusAnonymous  UserStatus = 0
	UserStatusActive     UserStatus = 1
	UserStatusRegistered UserStatus = 2
	UserStatusLocked     UserStatus = 3
)

// UserNotification const
const (
	UserNotificationAll          UserNotification = "all"
	UserNotificationSelected     UserNotification = "selected"
	UserNotificationOnlyMyEvents UserNotification = "only_my_events"
	UserNotificationOnlyAssigned UserNotification = "only_assigned"
	UserNotificationOnlyOwner    UserNotification = "only_owner"
	UserNotificationOnlyNone     UserNotification = "none"
)

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
	Status       UserStatus             `json:"status"`  // used only: get single user
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

// UserCreate struct used for users create operations
type UserCreate struct {
	User            UserCreateObject `json:"user"`
	SendInformation bool             `json:"send_information,omitempty"`
}

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
	CustomFields     []CustomFieldUpdateObject `json:"custom_fields,omitempty"`
}

/* Update */

// UserUpdate struct used for users update operations
type UserUpdate struct {
	User            UserUpdateObject `json:"user"`
	SendInformation bool             `json:"send_information,omitempty"`
}

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

// UserSingleGetRequest contains data for making request to get specified user
type UserSingleGetRequest struct {
	Includes []string
}

// UserCurrentGetRequest contains data for making request to get current user
type UserCurrentGetRequest struct {
	Includes []string
}

// UserGetRequestFilters contains data for making users get request
type UserGetRequestFilters struct {
	Status  UserStatus
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

func (u UserStatus) String() string {

	status := map[UserStatus]string{
		UserStatusAnonymous:  "anonymous",
		UserStatusActive:     "active",
		UserStatusRegistered: "registered",
		UserStatusLocked:     "locked",
	}

	s, b := status[u]
	if b == false {
		return "unknown"
	}

	return s
}

func (u UserNotification) String() string {
	return string(u)
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

	s, err := r.Get(&u, ur, http.StatusOK)

	return u, s, err
}

// UserSingleGet gets single user info by specific ID
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Users#GET-2
//
// Available includes:
// * groups
// * memberships
func (r *Context) UserSingleGet(id int, request UserSingleGetRequest) (UserObject, int, error) {

	var u userSingleResult

	urlParams := url.Values{}

	// Preparing includes
	urlIncludes(&urlParams, request.Includes)

	ur := url.URL{
		Path:     "/users/" + strconv.Itoa(id) + ".json",
		RawQuery: urlParams.Encode(),
	}

	status, err := r.Get(&u, ur, http.StatusOK)

	return u.User, status, err
}

// UserCurrentGet gets current user info
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Users#GET-2
//
// Available includes:
// * groups
// * memberships
func (r *Context) UserCurrentGet(request UserCurrentGetRequest) (UserObject, int, error) {

	var u userSingleResult

	urlParams := url.Values{}

	// Preparing includes
	urlIncludes(&urlParams, request.Includes)

	ur := url.URL{
		Path:     "/users/current.json",
		RawQuery: urlParams.Encode(),
	}

	status, err := r.Get(&u, ur, http.StatusOK)

	return u.User, status, err
}

// UserCreate creates new user
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Users#POST
func (r *Context) UserCreate(user UserCreate) (UserObject, int, error) {

	var u userSingleResult

	ur := url.URL{
		Path: "/users.json",
	}

	status, err := r.Post(user, &u, ur, http.StatusCreated)

	return u.User, status, err
}

// UserUpdate updates user with specified ID
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Users#PUT
func (r *Context) UserUpdate(id int, user UserUpdate) (int, error) {

	ur := url.URL{
		Path: "/users/" + strconv.Itoa(id) + ".json",
	}

	status, err := r.Put(user, nil, ur, http.StatusNoContent)

	return status, err
}

// UserDelete deletes user with specified ID
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Users#DELETE
func (r *Context) UserDelete(id int) (int, error) {

	ur := url.URL{
		Path: "/users/" + strconv.Itoa(id) + ".json",
	}

	status, err := r.Del(nil, nil, ur, http.StatusNoContent)

	return status, err
}

func userURLFilters(urlParams *url.Values, filters UserGetRequestFilters) {

	if filters.Status == 0 {
		filters.Status = 1
	}

	urlParams.Add("status", strconv.Itoa(int(filters.Status)))

	if len(filters.Name) > 0 {
		urlParams.Add("name", filters.Name)
	}

	if filters.GroupID > 0 {
		urlParams.Add("group_id", strconv.Itoa(filters.GroupID))
	}
}
