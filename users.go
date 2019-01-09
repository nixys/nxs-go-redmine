package redmine

import (
	"strconv"
	"strings"
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
	ID           int                     `json:"id"`
	Login        string                  `json:"login"`
	FirstName    string                  `json:"firstname"`
	LastName     string                  `json:"lastname"`
	Mail         string                  `json:"mail"`
	CreatedOn    string                  `json:"created_on"`
	LastLoginOn  string                  `json:"last_login_on"`
	APIKey       string                  `json:"api_key"` // used only: get single user
	Status       int                     `json:"status"`  // used only: get single user
	CustomFields []UserCustomFieldObject `json:"custom_fields"`
	Groups       []IDName                `json:"groups"`      // used only: get single user
	Memberships  []UserMembershipObject  `json:"memberships"` // used only: get single user
}

// UserCustomFieldObject struct used for users get operations
type UserCustomFieldObject struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
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
	Login            string `json:"login"`
	FirstName        string `json:"firstname"`
	LastName         string `json:"lastname"`
	Mail             string `json:"mail"`
	Password         string `json:"password,omitempty"`
	AuthSourceID     int    `json:"auth_source_id,omitempty"`
	MailNotification string `json:"mail_notification,omitempty"`
	MustChangePasswd bool   `json:"must_change_passwd,omitempty"`
	GeneratePassword bool   `json:"generate_password,omitempty"`
	SendInformation  bool   `json:"send_information,omitempty"`
}

/* Update */

// UserUpdateObject struct used for users update operations
type UserUpdateObject struct {
	Login            string `json:"login,omitempty"`
	FirstName        string `json:"firstname,omitempty"`
	LastName         string `json:"lastname,omitempty"`
	Mail             string `json:"mail,omitempty"`
	Password         string `json:"password,omitempty"`
	AuthSourceID     int    `json:"auth_source_id,omitempty"`
	MailNotification string `json:"mail_notification,omitempty"`
	MustChangePasswd bool   `json:"must_change_passwd,omitempty"`
	GeneratePassword bool   `json:"generate_password,omitempty"`
	SendInformation  bool   `json:"send_information,omitempty"`
}

/* Internal types */

type userMultiResult struct {
	Users      []UserObject `json:"users"`
	TotalCount int          `json:"total_count"`
	Offset     int          `json:"offset"`
	Limit      int          `json:"limit"`
}

type userSingleResult struct {
	User UserObject `json:"user"`
}

type userCreate struct {
	User UserCreateObject `json:"user"`
}

type userUpdate struct {
	User UserUpdateObject `json:"user"`
}

// UserMultiGet gets multiple users info by specific filters
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Users#GET
//
// * If `statusFilter` == 0 default users status filter will be used (show active users only)
// * Use `groupIDFilter` == 0 to disable this filter
func (r *Redmine) UserMultiGet(statusFilter int, nameFilter string, groupIDFilter int) ([]UserObject, int, error) {

	var u userMultiResult
	var status int

	if statusFilter == 0 {
		statusFilter = 1
	}

	filters := "&status=" + strconv.Itoa(statusFilter)

	if len(nameFilter) > 0 {
		filters += "&name=" + nameFilter
	}

	if groupIDFilter > 0 {
		filters += "&group_id=" + strconv.Itoa(groupIDFilter)
	}

	offset := 0

	for {

		uri := "/users.json?limit=" + strconv.Itoa(r.limit) + "&offset=" + strconv.Itoa(offset) + filters

		ut := userMultiResult{}

		s, err := r.get(&ut, uri, 200)
		if err != nil {
			return u.Users, s, err
		}

		status = s

		for _, e := range ut.Users {
			u.Users = append(u.Users, e)
		}

		if offset+ut.Limit >= ut.TotalCount {
			u.TotalCount = ut.TotalCount
			u.Limit = ut.TotalCount
			break
		}

		offset += ut.Limit
	}

	return u.Users, status, nil
}

// UserSingleGet gets single user info by specific ID
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Users#GET-2
//
// Available includes:
// * groups
// * memberships
func (r *Redmine) UserSingleGet(id int, includes []string) (UserObject, int, error) {

	var u userSingleResult
	var i string

	if len(includes) != 0 {
		i = "?include=" + strings.Join(includes, ",")
	}

	uri := "/users/" + strconv.Itoa(id) + ".json" + i

	status, err := r.get(&u, uri, 200)

	return u.User, status, err
}

// UserCreate creates new user
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Users#POST
func (r *Redmine) UserCreate(user UserCreateObject) (UserObject, int, error) {

	var u userSingleResult

	uri := "/users.json"

	status, err := r.post(userCreate{User: user}, &u, uri, 201)

	return u.User, status, err
}

// UserUpdate updates user with specified ID
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Users#PUT
func (r *Redmine) UserUpdate(id int, user UserUpdateObject) (int, error) {

	uri := "/users/" + strconv.Itoa(id) + ".json"

	status, err := r.put(userUpdate{User: user}, nil, uri, 200)

	return status, err
}

// UserDelete deletes user with specified ID
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Users#DELETE
func (r *Redmine) UserDelete(id int) (int, error) {

	uri := "/users/" + strconv.Itoa(id) + ".json"

	status, err := r.del(nil, nil, uri, 200)

	return status, err
}
