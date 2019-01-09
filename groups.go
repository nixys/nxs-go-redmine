package redmine

import (
	"strconv"
	"strings"
)

/* Get */

type GroupObject struct {
	ID          int                     `json:"id"`
	Name        string                  `json:"name"`
	Users       []IDName                `json:"users"`       /* used only: get single user*/
	Memberships []GroupMembershipObject `json:"memberships"` /* used only: get single user*/
}

type GroupMembershipObject struct {
	ID      int      `json:"id"`
	Project IDName   `json:"project"`
	Roles   []IDName `json:"roles"`
}

/* Create */

type GroupCreateObject struct {
	Name    string `json:"name"`
	UserIDs []int  `json:"user_ids,omitempty"`
}

/* Update */

type GroupUpdateObject struct {
	Name    string `json:"name,omitempty"`
	UserIDs []int  `json:"user_ids,omitempty"`
}

/* Add user */

type GroupAddUserObject struct {
	UserID int `json:"user_id"`
}

/* Internal types */

type groupMultiResult struct {
	Groups     []GroupObject `json:"groups"`
	TotalCount int           `json:"total_count"`
	Offset     int           `json:"offset"`
	Limit      int           `json:"limit"`
}

type groupSingleResult struct {
	Group GroupObject `json:"group"`
}

type groupCreate struct {
	Group GroupCreateObject `json:"group"`
}

type groupUpdate struct {
	Group GroupUpdateObject `json:"group"`
}

/* see: http://www.redmine.org/projects/redmine/wiki/Rest_Groups#GET */
func (r *Redmine) GroupMultiGet() ([]GroupObject, int, error) {

	var g groupMultiResult
	var status int

	offset := 0

	for {

		uri := "/groups.json?limit=" + strconv.Itoa(r.limit) + "&offset=" + strconv.Itoa(offset)

		gt := groupMultiResult{}

		s, err := r.get(&gt, uri, 200)
		if err != nil {
			return g.Groups, s, err
		}

		status = s

		for _, e := range gt.Groups {
			g.Groups = append(g.Groups, e)
		}

		if offset+gt.Limit >= gt.TotalCount {
			g.TotalCount = gt.TotalCount
			g.Limit = gt.TotalCount
			break
		}

		offset += gt.Limit
	}

	return g.Groups, status, nil
}

/*
see: http://www.redmine.org/projects/redmine/wiki/Rest_Groups#GET-2

Available includes:
* users
* memberships
*/
func (r *Redmine) GroupSingleGet(id int, includes []string) (GroupObject, int, error) {

	var g groupSingleResult
	var i string

	if len(includes) != 0 {
		i = "?include=" + strings.Join(includes, ",")
	}

	uri := "/groups/" + strconv.Itoa(id) + ".json" + i

	status, err := r.get(&g, uri, 200)

	return g.Group, status, err
}

/* see: http://www.redmine.org/projects/redmine/wiki/Rest_Groups#POST */
func (r *Redmine) GroupCreate(group GroupCreateObject) (GroupObject, int, error) {

	var g groupSingleResult

	uri := "/groups.json"

	status, err := r.post(groupCreate{Group: group}, &g, uri, 201)

	return g.Group, status, err
}

/* see: http://www.redmine.org/projects/redmine/wiki/Rest_Groups#PUT */
func (r *Redmine) GroupUpdate(id int, group GroupUpdateObject) (int, error) {

	uri := "/groups/" + strconv.Itoa(id) + ".json"

	status, err := r.put(groupUpdate{Group: group}, nil, uri, 200)

	return status, err
}

/* see: http://www.redmine.org/projects/redmine/wiki/Rest_Groups#DELETE */
func (r *Redmine) GroupDelete(id int) (int, error) {

	uri := "/groups/" + strconv.Itoa(id) + ".json"

	status, err := r.del(nil, nil, uri, 200)

	return status, err
}

/* see: http://www.redmine.org/projects/redmine/wiki/Rest_Groups#POST-2 */
func (r *Redmine) GroupAddUser(id int, group GroupAddUserObject) (int, error) {

	uri := "/groups/" + strconv.Itoa(id) + "/users.json"

	status, err := r.post(group, nil, uri, 200)

	return status, err
}

/* see: http://www.redmine.org/projects/redmine/wiki/Rest_Groups#DELETE-2 */
func (r *Redmine) GroupDeleteUser(id int, userID int) (int, error) {

	uri := "/groups/" + strconv.Itoa(id) + "/users/" + strconv.Itoa(userID) + ".json"

	status, err := r.del(nil, nil, uri, 200)

	return status, err
}
