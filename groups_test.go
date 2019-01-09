package redmine

import (
	"testing"
)

const (
	testGroupName  = "test-group"
	testGroupName2 = "test-group2"
)

func TestGroupsCRUD(t *testing.T) {

	var r Redmine

	/* Init Redmine context */
	initTest(&r, t)

	/* Preparing auxiliary data */
	uCreated := testUserCreate(t, r)
	defer testUserDetele(t, r, uCreated.ID)

	/* Create and delete */
	gCreated := testGroupCreate(t, r)
	defer testGroupDetele(t, r, gCreated.ID)

	/* Get multi */
	testGroupMultiGet(t, r)

	/* Update */
	testGroupUpdate(t, r, gCreated.ID, uCreated.ID)

	/* Get single (check user is member of specified group) */
	testGroupSingleGet(t, r, gCreated.ID, uCreated.ID)

	/* Delete user */
	testGroupDeteleUser(t, r, gCreated.ID, uCreated.ID)

	/* Add user */
	testGroupAddUser(t, r, gCreated.ID, uCreated.ID)
}

func testGroupCreate(t *testing.T, r Redmine) GroupObject {

	g, _, err := r.GroupCreate(GroupCreateObject{
		Name: testGroupName,
	})
	if err != nil {
		t.Fatal("Group create error:", err)
	}

	t.Logf("Group create: success")

	return g
}

func testGroupUpdate(t *testing.T, r Redmine, id, userID int) {

	_, err := r.GroupUpdate(id, GroupUpdateObject{
		Name:    testGroupName2,
		UserIDs: []int{userID},
	})
	if err != nil {
		t.Fatal("Group update error:", err)
	}

	t.Logf("Group update: success")
}

func testGroupAddUser(t *testing.T, r Redmine, id, userID int) {

	_, err := r.GroupAddUser(id, GroupAddUserObject{UserID: userID})
	if err != nil {
		t.Fatal("Group add user error:", err)
	}

	t.Logf("Group add user: success")
}

func testGroupDeteleUser(t *testing.T, r Redmine, id, userID int) {

	_, err := r.GroupDeleteUser(id, userID)
	if err != nil {
		t.Fatal("Group delete user error:", err)
	}

	t.Logf("Group delete user: success")
}

func testGroupDetele(t *testing.T, r Redmine, id int) {

	_, err := r.GroupDelete(id)
	if err != nil {
		t.Fatal("Group delete error:", err)
	}

	t.Logf("Group delete: success")
}

func testGroupMultiGet(t *testing.T, r Redmine) {

	g, _, err := r.GroupMultiGet()
	if err != nil {
		t.Fatal("Groups get error:", err)
	}

	for _, e := range g {
		if e.Name == testGroupName {
			t.Logf("Groups get: success")
			return
		}
	}

	t.Fatal("Groups get error: can't find created group")
}

func testGroupSingleGet(t *testing.T, r Redmine, id, userID int) {

	g, _, err := r.GroupSingleGet(id, []string{"users", "memberships"})
	if err != nil {
		t.Fatal("Group get error:", err)
	}

	/* Check user is a member of specified group (error if not) */

	for _, e := range g.Users {
		if e.ID == userID {
			t.Logf("Group get: success")
			return
		}
	}

	t.Fatal("Group get error: can't find user in specified group")
}
