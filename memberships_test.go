package redmine

import (
	"os"
	"strconv"
	"testing"
)

func TestMembershipCRUD(t *testing.T) {

	var r Context

	// Get env variables
	testMembershipRoleID1, _ := strconv.Atoi(os.Getenv("REDMINE_ROLE_ID_1"))
	testMembershipRoleID2, _ := strconv.Atoi(os.Getenv("REDMINE_ROLE_ID_2"))

	if testMembershipRoleID1 == 0 || testMembershipRoleID2 == 0 {
		t.Fatal("Membership test error: env variables `REDMINE_ROLE_ID_1` or `REDMINE_ROLE_ID_2` does not set")
	}

	// Init Redmine context
	initTest(&r, t)

	// Preparing auxiliary data
	uCreated := testUserCreate(t, r)
	defer testUserDetele(t, r, uCreated.ID)

	pCreated := testProjectCreate(t, r, []int{})
	defer testProjectDetele(t, r, pCreated.Identifier)

	// Add and delete
	mCreated := testMembershipAdd(t, r, pCreated.Identifier, uCreated.ID, testMembershipRoleID1)
	defer testMembershipDetele(t, r, mCreated.ID)

	// Get all
	testMembershipAllGet(t, r, mCreated.ID, pCreated.Identifier, testMembershipRoleID1)

	// Update
	testMembershipUpdate(t, r, mCreated.ID, testMembershipRoleID1, testMembershipRoleID2)

	// Get single (check role `testMembershipRoleID2` exists in specified membership)
	testMembershipSingleGet(t, r, mCreated.ID, testMembershipRoleID2)
}

func testMembershipAdd(t *testing.T, r Context, projectID string, userID, roleID int) MembershipObject {

	m, _, err := r.MembershipAdd(projectID, MembershipAddObject{
		UserID:  userID,
		RoleIDs: []int{roleID},
	})
	if err != nil {
		t.Fatal("Membership add error:", err)
	}

	t.Logf("Membership add: success")

	return m
}

func testMembershipUpdate(t *testing.T, r Context, id, roleID1, roleID2 int) {

	_, err := r.MembershipUpdate(id, MembershipUpdateObject{
		RoleIDs: []int{roleID1, roleID2},
	})
	if err != nil {
		t.Fatal("Membership update error:", err)
	}

	t.Logf("Membership update: success")
}

func testMembershipDetele(t *testing.T, r Context, id int) {

	_, err := r.MembershipDelete(id)
	if err != nil {
		t.Fatal("Membership delete error:", err)
	}

	t.Logf("Membership delete: success")
}

func testMembershipAllGet(t *testing.T, r Context, id int, projectID string, roleID int) {

	m, _, err := r.MembershipAllGet(projectID)
	if err != nil {
		t.Fatal("Memberships get error:", err)
	}

	for _, e := range m.Memberships {
		if e.ID == id {
			for _, role := range e.Roles {
				if role.ID == roleID {
					t.Logf("Memberships get: success")
					return
				}
			}

			t.Fatal("Memberships get error: can't find role in added membership")
		}
	}

	t.Fatal("Memberships get error: can't find added membership")
}

func testMembershipSingleGet(t *testing.T, r Context, id, roleID int) {

	m, _, err := r.MembershipSingleGet(id)
	if err != nil {
		t.Fatal("Membership get error:", err)
	}

	for _, role := range m.Roles {
		if role.ID == roleID {
			t.Logf("Membership get: success")
			return
		}
	}

	t.Fatal("Membership get error: can't find role in added membership")
}
