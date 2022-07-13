package redmine

import (
	"net/http"
	"net/url"
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

/* Requests */

// MembershipMultiGetRequest contains data for making request to get limited memberships count
type MembershipMultiGetRequest struct {
	Offset int
	Limit  int
}

/* Results */

// MembershipResult stores project memberships requests processing result
type MembershipResult struct {
	Memberships []MembershipObject `json:"memberships"`
	TotalCount  int                `json:"total_count"`
	Offset      int                `json:"offset"`
	Limit       int                `json:"limit"`
}

/* Internal types */

type membershipSingleResult struct {
	Membership MembershipObject `json:"membership"`
}

type membershipAdd struct {
	Membership MembershipAddObject `json:"membership"`
}

type membershipUpdate struct {
	Membership MembershipUpdateObject `json:"membership"`
}

// MembershipAllGet gets info for all memberships for project with specified ID
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Memberships#GET
func (r *Context) MembershipAllGet(projectID string) (MembershipResult, int, error) {

	var (
		membership     MembershipResult
		offset, status int
	)

	m := MembershipMultiGetRequest{
		Limit: limitDefault,
	}

	for {

		m.Offset = offset

		m, s, err := r.MembershipMultiGet(projectID, m)
		if err != nil {
			return membership, s, err
		}

		status = s

		membership.Memberships = append(membership.Memberships, m.Memberships...)

		if offset+m.Limit >= m.TotalCount {
			membership.TotalCount = m.TotalCount
			membership.Limit = m.TotalCount

			break
		}

		offset += m.Limit
	}

	return membership, status, nil
}

// MembershipMultiGet gets info for multiple memberships for project with specified ID
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Memberships#GET
func (r *Context) MembershipMultiGet(projectID string, request MembershipMultiGetRequest) (MembershipResult, int, error) {

	var m MembershipResult

	urlParams := url.Values{}
	urlParams.Add("offset", strconv.Itoa(request.Offset))
	urlParams.Add("limit", strconv.Itoa(request.Limit))

	ur := url.URL{
		Path:     "/projects/" + projectID + "/memberships.json",
		RawQuery: urlParams.Encode(),
	}

	s, err := r.Get(&m, ur, http.StatusOK)

	return m, s, err
}

// MembershipSingleGet gets single project membership info with specified ID
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Memberships#GET-2
func (r *Context) MembershipSingleGet(membershipID int) (MembershipObject, int, error) {

	var m membershipSingleResult

	ur := url.URL{
		Path: "/memberships/" + strconv.Itoa(membershipID) + ".json",
	}

	status, err := r.Get(&m, ur, http.StatusOK)

	return m.Membership, status, err
}

// MembershipAdd adds new member to project with specified ID
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Memberships#POST
func (r *Context) MembershipAdd(projectID string, membership MembershipAddObject) (MembershipObject, int, error) {

	var m membershipSingleResult

	ur := url.URL{
		Path: "/projects/" + projectID + "/memberships.json",
	}

	status, err := r.Post(membershipAdd{Membership: membership}, &m, ur, http.StatusCreated)

	return m.Membership, status, err
}

// MembershipUpdate updates project membership with specified ID
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Memberships#PUT
func (r *Context) MembershipUpdate(membershipID int, membership MembershipUpdateObject) (int, error) {

	ur := url.URL{
		Path: "/memberships/" + strconv.Itoa(membershipID) + ".json",
	}

	status, err := r.Put(membershipUpdate{Membership: membership}, nil, ur, http.StatusNoContent)

	return status, err
}

// MembershipDelete deletes project membership with specified ID
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Memberships#DELETE
func (r *Context) MembershipDelete(membershipID int) (int, error) {

	ur := url.URL{
		Path: "/memberships/" + strconv.Itoa(membershipID) + ".json",
	}

	status, err := r.Del(nil, nil, ur, http.StatusNoContent)

	return status, err
}
