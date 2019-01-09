package redmine

import (
	"os"
	"strconv"
	"testing"
)

var (
	testMembershipRoleID1 int
	testMembershipRoleID2 int
)

func TestMembershipCRUD(t *testing.T) {

	var r Redmine

	/* Get env variables */
	testMembershipRoleID1, _ = strconv.Atoi(os.Getenv("REDMINE_ROLE_ID_1"))
	testMembershipRoleID2, _ = strconv.Atoi(os.Getenv("REDMINE_ROLE_ID_2"))

	if testMembershipRoleID1 == 0 || testMembershipRoleID2 == 0 {
		t.Fatal("Membership test error: env variables `REDMINE_ROLE_ID_1` or `REDMINE_ROLE_ID_2` does not set")
	}

	/* Init Redmine context */
	initTest(&r, t)

	/* Preparing auxiliary data */
	uCreated := testUserCreate(t, r)
	defer testUserDetele(t, r, uCreated.ID)

	pCreated := testProjectCreate(t, r)
	defer testProjectDetele(t, r, pCreated.ID)

	/* Add and delete */
	mCreated := testMembershipAdd(t, r, pCreated.ID, uCreated.ID)
	defer testMembershipDetele(t, r, mCreated.ID)

	/* Get multi */
	testMembershipMultiGet(t, r, mCreated.ID, pCreated.ID)

	/* Update */
	testMembershipUpdate(t, r, mCreated.ID)

	/* Get single (check role `testMembershipRoleID2` exists in specified membership) */
	testMembershipSingleGet(t, r, mCreated.ID)
}

func testMembershipAdd(t *testing.T, r Redmine, projectID, userID int) MembershipObject {

	m, _, err := r.MembershipAdd(projectID, MembershipAddObject{
		UserID:  userID,
		RoleIDs: []int{testMembershipRoleID1},
	})
	if err != nil {
		t.Fatal("Membership add error:", err)
	}

	t.Logf("Membership add: success")

	return m
}

func testMembershipUpdate(t *testing.T, r Redmine, id int) {

	_, err := r.MembershipUpdate(id, MembershipUpdateObject{
		RoleIDs: []int{testMembershipRoleID1, testMembershipRoleID2},
	})
	if err != nil {
		t.Fatal("Membership update error:", err)
	}

	t.Logf("Membership update: success")
}

func testMembershipDetele(t *testing.T, r Redmine, id int) {

	_, err := r.MembershipDelete(id)
	if err != nil {
		t.Fatal("Membership delete error:", err)
	}

	t.Logf("Membership delete: success")
}

func testMembershipMultiGet(t *testing.T, r Redmine, id, projectID int) {

	m, _, err := r.MembershipMultiGet(projectID)
	if err != nil {
		t.Fatal("Memberships get error:", err)
	}

	for _, e := range m {
		if e.ID == id {
			for _, role := range e.Roles {
				if role.ID == testMembershipRoleID1 {
					t.Logf("Memberships get: success")
					return
				}
			}

			t.Fatal("Memberships get error: can't find role in added membership")
		}
	}

	t.Fatal("Memberships get error: can't find added membership")
}

func testMembershipSingleGet(t *testing.T, r Redmine, id int) {

	m, _, err := r.MembershipSingleGet(id)
	if err != nil {
		t.Fatal("Membership get error:", err)
	}

	for _, role := range m.Roles {
		if role.ID == testMembershipRoleID2 {
			t.Logf("Membership get: success")
			return
		}
	}

	t.Fatal("Membership get error: can't find role in added membership")
}
