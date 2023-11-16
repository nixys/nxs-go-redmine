package redmine

import (
	"net/http"
	"net/url"
	"strconv"
)

/* Get */

// MembershipObject struct used for project memberships get operations
type MembershipObject struct {
	ID      int64                  `json:"id"`
	Project IDName                 `json:"project"`
	User    IDName                 `json:"user"`
	Group   IDName                 `json:"group"`
	Roles   []MembershipRoleObject `json:"roles"`
}

// MembershipRoleObject struct used for project memberships get operations
type MembershipRoleObject struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Inherited bool   `json:"inherited"`
}

/* Add */

// MembershipAdd struct used for project memberships add operations
type MembershipAdd struct {
	Membership MembershipAddObject `json:"membership"`
}

type MembershipAddObject struct {
	UserID  int64   `json:"user_id"`
	RoleIDs []int64 `json:"role_ids"`
}

/* Update */

// MembershipUpdate struct used for project memberships update operations
type MembershipUpdate struct {
	Membership MembershipUpdateObject `json:"membership"`
}

type MembershipUpdateObject struct {
	RoleIDs []int64 `json:"role_ids"`
}

/* Requests */

// MembershipMultiGetRequest contains data for making request to get limited memberships count
type MembershipMultiGetRequest struct {
	Offset int64
	Limit  int64
}

/* Results */

// MembershipResult stores project memberships requests processing result
type MembershipResult struct {
	Memberships []MembershipObject `json:"memberships"`
	TotalCount  int64              `json:"total_count"`
	Offset      int64              `json:"offset"`
	Limit       int64              `json:"limit"`
}

/* Internal types */

type membershipSingleResult struct {
	Membership MembershipObject `json:"membership"`
}

// MembershipAllGet gets info for all memberships for project with specified ID
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Memberships#GET
func (r *Context) MembershipAllGet(projectID string) (MembershipResult, StatusCode, error) {

	var (
		membership MembershipResult
		offset     int64
		status     StatusCode
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
func (r *Context) MembershipMultiGet(projectID string, request MembershipMultiGetRequest) (MembershipResult, StatusCode, error) {

	var m MembershipResult

	urlParams := url.Values{}
	urlParams.Add("offset", strconv.FormatInt(request.Offset, 10))
	urlParams.Add("limit", strconv.FormatInt(request.Limit, 10))

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
func (r *Context) MembershipSingleGet(membershipID int64) (MembershipObject, StatusCode, error) {

	var m membershipSingleResult

	ur := url.URL{
		Path: "/memberships/" + strconv.FormatInt(membershipID, 10) + ".json",
	}

	status, err := r.Get(&m, ur, http.StatusOK)

	return m.Membership, status, err
}

// MembershipAdd adds new member to project with specified ID
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Memberships#POST
func (r *Context) MembershipAdd(projectID string, membership MembershipAdd) (MembershipObject, StatusCode, error) {

	var m membershipSingleResult

	ur := url.URL{
		Path: "/projects/" + projectID + "/memberships.json",
	}

	status, err := r.Post(membership, &m, ur, http.StatusCreated)

	return m.Membership, status, err
}

// MembershipUpdate updates project membership with specified ID
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Memberships#PUT
func (r *Context) MembershipUpdate(membershipID int64, membership MembershipUpdate) (StatusCode, error) {

	ur := url.URL{
		Path: "/memberships/" + strconv.FormatInt(membershipID, 10) + ".json",
	}

	status, err := r.Put(membership, nil, ur, http.StatusNoContent)

	return status, err
}

// MembershipDelete deletes project membership with specified ID
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Memberships#DELETE
func (r *Context) MembershipDelete(membershipID int64) (StatusCode, error) {

	ur := url.URL{
		Path: "/memberships/" + strconv.FormatInt(membershipID, 10) + ".json",
	}

	status, err := r.Del(nil, nil, ur, http.StatusNoContent)

	return status, err
}
