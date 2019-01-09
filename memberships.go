package redmine

import (
	"strconv"
)

/* Get */

// MembershipObject struct used for project memberships get operations
type MembershipObject struct {
	ID      int                    `json:"id"`
	Project IDName                 `json:"project"`
	User    IDName                 `json:"user"`
	Group   IDName                 `json:"group"`
	Roles   []MembershipRoleObject `json:"roles"`
}

// MembershipRoleObject struct used for project memberships get operations
type MembershipRoleObject struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Inherited bool   `json:"inherited"`
}

/* Add */

// MembershipAddObject struct used for project memberships add operations
type MembershipAddObject struct {
	UserID  int   `json:"user_id"`
	RoleIDs []int `json:"role_ids"`
}

/* Update */

// MembershipUpdateObject struct used for project memberships update operations
type MembershipUpdateObject struct {
	RoleIDs []int `json:"role_ids"`
}

/* Internal types */

type membershipMultiResult struct {
	Memberships []MembershipObject `json:"memberships"`
	TotalCount  int                `json:"total_count"`
	Offset      int                `json:"offset"`
	Limit       int                `json:"limit"`
}

type membershipSingleResult struct {
	Membership MembershipObject `json:"membership"`
}

type membershipAdd struct {
	Membership MembershipAddObject `json:"membership"`
}

type membershipUpdate struct {
	Membership MembershipUpdateObject `json:"membership"`
}

// MembershipMultiGet gets multiple memberships info for project with specified ID
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Memberships#GET
func (r *Redmine) MembershipMultiGet(projectID int) ([]MembershipObject, int, error) {

	var m membershipMultiResult
	var status int

	offset := 0

	for {
		uri := "/projects/" + strconv.Itoa(projectID) + "/memberships.json?limit=" + strconv.Itoa(r.limit) + "&offset=" + strconv.Itoa(offset)

		mt := membershipMultiResult{}

		s, err := r.get(&mt, uri, 200)
		if err != nil {
			return m.Memberships, s, err
		}

		status = s

		for _, e := range mt.Memberships {
			m.Memberships = append(m.Memberships, e)
		}

		if offset+mt.Limit >= mt.TotalCount {
			m.TotalCount = mt.TotalCount
			m.Limit = mt.TotalCount
			break
		}

		offset += mt.Limit
	}

	return m.Memberships, status, nil
}

// MembershipSingleGet gets single project membership info with specified ID
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Memberships#GET-2
func (r *Redmine) MembershipSingleGet(membershipID int) (MembershipObject, int, error) {

	var m membershipSingleResult

	uri := "/memberships/" + strconv.Itoa(membershipID) + ".json"

	status, err := r.get(&m, uri, 200)

	return m.Membership, status, err
}

// MembershipAdd adds new member to project with specified ID
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Memberships#POST
func (r *Redmine) MembershipAdd(projectID int, membership MembershipAddObject) (MembershipObject, int, error) {

	var m membershipSingleResult

	uri := "/projects/" + strconv.Itoa(projectID) + "/memberships.json"

	status, err := r.post(membershipAdd{Membership: membership}, &m, uri, 201)

	return m.Membership, status, err
}

// MembershipUpdate updates project membership with specified ID
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Memberships#PUT
func (r *Redmine) MembershipUpdate(membershipID int, membership MembershipUpdateObject) (int, error) {

	uri := "/memberships/" + strconv.Itoa(membershipID) + ".json"

	status, err := r.put(membershipUpdate{Membership: membership}, nil, uri, 200)

	return status, err
}

// MembershipDelete deletes project membership with specified ID
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Memberships#DELETE
func (r *Redmine) MembershipDelete(membershipID int) (int, error) {

	uri := "/memberships/" + strconv.Itoa(membershipID) + ".json"

	status, err := r.del(nil, nil, uri, 200)

	return status, err
}
