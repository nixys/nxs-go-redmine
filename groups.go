package redmine

import (
	"net/http"
	"net/url"
	"strconv"
)

/* Get */

// GroupObject struct used for groups get operations
type GroupObject struct {
	ID          int64                   `json:"id"`
	Name        string                  `json:"name"`
	Users       []IDName                `json:"users"`       // used only: get single user
	Memberships []GroupMembershipObject `json:"memberships"` // used only: get single user
}

// GroupMembershipObject struct used for groups get operations
type GroupMembershipObject struct {
	ID      int64    `json:"id"`
	Project IDName   `json:"project"`
	Roles   []IDName `json:"roles"`
}

/* Create */

// GroupCreate struct used for groups create operations
type GroupCreate struct {
	Group GroupCreateObject `json:"group"`
}

type GroupCreateObject struct {
	Name    string  `json:"name"`
	UserIDs []int64 `json:"user_ids,omitempty"`
}

/* Update */

// GroupUpdate struct used for groups update operations
type GroupUpdate struct {
	Group GroupUpdateObject `json:"group"`
}

type GroupUpdateObject struct {
	Name    string  `json:"name,omitempty"`
	UserIDs []int64 `json:"user_ids,omitempty"`
}

/* Add user */

// GroupAddUserObject struct used for add new user into group
type GroupAddUserObject struct {
	UserID int64 `json:"user_id"`
}

/* Requests */

// GroupMultiGetRequest contains data for making request to get limited groups count
type GroupMultiGetRequest struct {
	Offset int64
	Limit  int64
}

// GroupSingleGetRequest contains data for making request to get specified group
type GroupSingleGetRequest struct {
	Includes []string
}

/* Results */

// GroupResult stores groups requests processing result
type GroupResult struct {
	Groups     []GroupObject `json:"groups"`
	TotalCount int64         `json:"total_count"`
	Offset     int64         `json:"offset"`
	Limit      int64         `json:"limit"`
}

/* Internal types */

type groupSingleResult struct {
	Group GroupObject `json:"group"`
}

// GroupAllGet gets info for all groups
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Groups#GET
func (r *Context) GroupAllGet() (GroupResult, StatusCode, error) {

	var (
		groups GroupResult
		offset int64
		status StatusCode
	)

	m := GroupMultiGetRequest{
		Limit: limitDefault,
	}

	for {

		m.Offset = offset

		g, s, err := r.GroupMultiGet(m)
		if err != nil {
			return groups, s, err
		}

		status = s

		groups.Groups = append(groups.Groups, g.Groups...)

		if offset+g.Limit >= g.TotalCount {
			groups.TotalCount = g.TotalCount
			groups.Limit = g.TotalCount

			break
		}

		offset += g.Limit
	}

	return groups, status, nil
}

// GroupMultiGet gets info for multiple groups
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Groups#GET
func (r *Context) GroupMultiGet(request GroupMultiGetRequest) (GroupResult, StatusCode, error) {

	var g GroupResult

	urlParams := url.Values{}
	urlParams.Add("offset", strconv.FormatInt(request.Offset, 10))
	urlParams.Add("limit", strconv.FormatInt(request.Limit, 10))

	ur := url.URL{
		Path:     "/groups.json",
		RawQuery: urlParams.Encode(),
	}

	s, err := r.Get(&g, ur, http.StatusOK)

	return g, s, err
}

// GroupSingleGet gets single group info by specific ID
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Groups#GET-2
//
// Available includes:
// * users
// * memberships
func (r *Context) GroupSingleGet(id int64, request GroupSingleGetRequest) (GroupObject, StatusCode, error) {

	var g groupSingleResult

	urlParams := url.Values{}

	// Preparing includes
	urlIncludes(&urlParams, request.Includes)

	ur := url.URL{
		Path:     "/groups/" + strconv.FormatInt(id, 10) + ".json",
		RawQuery: urlParams.Encode(),
	}

	status, err := r.Get(&g, ur, http.StatusOK)

	return g.Group, status, err
}

// GroupCreate creates new group
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Groups#POST
func (r *Context) GroupCreate(group GroupCreate) (GroupObject, StatusCode, error) {

	var g groupSingleResult

	ur := url.URL{
		Path: "/groups.json",
	}

	status, err := r.Post(group, &g, ur, http.StatusCreated)

	return g.Group, status, err
}

// GroupUpdate updates group with specified ID
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Groups#PUT
func (r *Context) GroupUpdate(id int64, group GroupUpdate) (StatusCode, error) {

	ur := url.URL{
		Path: "/groups/" + strconv.FormatInt(id, 10) + ".json",
	}

	status, err := r.Put(group, nil, ur, http.StatusNoContent)

	return status, err
}

// GroupDelete deletes group with specified ID
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Groups#DELETE
func (r *Context) GroupDelete(id int64) (StatusCode, error) {

	ur := url.URL{
		Path: "/groups/" + strconv.FormatInt(id, 10) + ".json",
	}

	status, err := r.Del(nil, nil, ur, http.StatusNoContent)

	return status, err
}

// GroupAddUser adds new user into group with specified ID
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Groups#POST-2
func (r *Context) GroupAddUser(id int64, group GroupAddUserObject) (StatusCode, error) {

	ur := url.URL{
		Path: "/groups/" + strconv.FormatInt(id, 10) + "/users.json",
	}

	status, err := r.Post(group, nil, ur, http.StatusNoContent)

	return status, err
}

// GroupDeleteUser deletes user from group with specified ID
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Groups#DELETE-2
func (r *Context) GroupDeleteUser(id int64, userID int64) (StatusCode, error) {

	ur := url.URL{
		Path: "/groups/" + strconv.FormatInt(id, 10) + "/users/" + strconv.FormatInt(userID, 10) + ".json",
	}

	status, err := r.Del(nil, nil, ur, http.StatusNoContent)

	return status, err
}
