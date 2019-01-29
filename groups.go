package redmine

import (
	"strconv"
	"strings"
)

/* Get */

// GroupObject struct used for groups get operations
type GroupObject struct {
	ID          int                     `json:"id"`
	Name        string                  `json:"name"`
	Users       []IDName                `json:"users"`       // used only: get single user
	Memberships []GroupMembershipObject `json:"memberships"` // used only: get single user
}

// GroupMembershipObject struct used for groups get operations
type GroupMembershipObject struct {
	ID      int      `json:"id"`
	Project IDName   `json:"project"`
	Roles   []IDName `json:"roles"`
}

/* Create */

// GroupCreateObject struct used for groups create operations
type GroupCreateObject struct {
	Name    string `json:"name"`
	UserIDs []int  `json:"user_ids,omitempty"`
}

/* Update */

// GroupUpdateObject struct used for groups update operations
type GroupUpdateObject struct {
	Name    string `json:"name,omitempty"`
	UserIDs []int  `json:"user_ids,omitempty"`
}

/* Add user */

// GroupAddUserObject struct used for add new user into group
type GroupAddUserObject struct {
	UserID int `json:"user_id"`
}

/* Results */

// GroupResult stores groups requests processing result
type GroupResult struct {
	Groups     []GroupObject `json:"groups"`
	TotalCount int           `json:"total_count"`
	Offset     int           `json:"offset"`
	Limit      int           `json:"limit"`
}

/* Internal types */

type groupSingleResult struct {
	Group GroupObject `json:"group"`
}

type groupCreate struct {
	Group GroupCreateObject `json:"group"`
}

type groupUpdate struct {
	Group GroupUpdateObject `json:"group"`
}

// GroupAllGet gets info for all groups
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Groups#GET
func (r *Context) GroupAllGet() (GroupResult, int, error) {

	var (
		groups         GroupResult
		offset, status int
	)

	for {

		g, s, err := r.GroupMultiGet(offset, limitDefault)
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
func (r *Context) GroupMultiGet(offset, limit int) (GroupResult, int, error) {

	var g GroupResult

	uri := "/groups.json?limit=" + strconv.Itoa(limit) + "&offset=" + strconv.Itoa(offset)

	s, err := r.get(&g, uri, 200)

	return g, s, err
}

// GroupSingleGet gets single group info by specific ID
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Groups#GET-2
//
// Available includes:
// * users
// * memberships
func (r *Context) GroupSingleGet(id int, includes []string) (GroupObject, int, error) {

	var g groupSingleResult
	var i string

	if len(includes) != 0 {
		i = "?include=" + strings.Join(includes, ",")
	}

	uri := "/groups/" + strconv.Itoa(id) + ".json" + i

	status, err := r.get(&g, uri, 200)

	return g.Group, status, err
}

// GroupCreate creates new group
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Groups#POST
func (r *Context) GroupCreate(group GroupCreateObject) (GroupObject, int, error) {

	var g groupSingleResult

	uri := "/groups.json"

	status, err := r.post(groupCreate{Group: group}, &g, uri, 201)

	return g.Group, status, err
}

// GroupUpdate updates group with specified ID
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Groups#PUT
func (r *Context) GroupUpdate(id int, group GroupUpdateObject) (int, error) {

	uri := "/groups/" + strconv.Itoa(id) + ".json"

	status, err := r.put(groupUpdate{Group: group}, nil, uri, 200)

	return status, err
}

// GroupDelete deletes group with specified ID
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Groups#DELETE
func (r *Context) GroupDelete(id int) (int, error) {

	uri := "/groups/" + strconv.Itoa(id) + ".json"

	status, err := r.del(nil, nil, uri, 200)

	return status, err
}

// GroupAddUser adds new user into group with specified ID
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Groups#POST-2
func (r *Context) GroupAddUser(id int, group GroupAddUserObject) (int, error) {

	uri := "/groups/" + strconv.Itoa(id) + "/users.json"

	status, err := r.post(group, nil, uri, 200)

	return status, err
}

// GroupDeleteUser deletes user from group with specified ID
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Groups#DELETE-2
func (r *Context) GroupDeleteUser(id int, userID int) (int, error) {

	uri := "/groups/" + strconv.Itoa(id) + "/users/" + strconv.Itoa(userID) + ".json"

	status, err := r.del(nil, nil, uri, 200)

	return status, err
}
