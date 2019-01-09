package redmine

import (
	"testing"
)

const (
	testUserLogin      = "test-user-login"
	testUserFirstName  = "First"
	testUserLastName   = "Last"
	testUserFirstName2 = "First2"
	testUserLastName2  = "Last2"
	testUserMail       = "test@domain.com"
)

func TestUserCRUD(t *testing.T) {

	var r Redmine

	/* Init Redmine context */
	initTest(&r, t)

	/* Create and delete */
	uCreated := testUserCreate(t, r)
	defer testUserDetele(t, r, uCreated.ID)

	/* Get */
	testUserMultiGet(t, r)
	testUserSingleGet(t, r, uCreated.ID)

	/* Update */
	testUserUpdate(t, r, uCreated.ID)
}

func testUserCreate(t *testing.T, r Redmine) UserObject {

	u, _, err := r.UserCreate(UserCreateObject{
		Login:            testUserLogin,
		FirstName:        testUserFirstName,
		LastName:         testUserLastName,
		Mail:             testUserMail,
		MailNotification: UserNotification[UserNotificationOnlyAssigned],
		MustChangePasswd: true,
		GeneratePassword: true,
	})
	if err != nil {
		t.Fatal("User create error:", err)
	}

	t.Logf("User create: success")

	return u
}

func testUserUpdate(t *testing.T, r Redmine, id int) {

	_, err := r.UserUpdate(id, UserUpdateObject{
		FirstName:        testUserFirstName2,
		LastName:         testUserLastName2,
		MailNotification: UserNotification[UserNotificationOnlyNone],
	})
	if err != nil {
		t.Fatal("User update error:", err)
	}

	t.Logf("User update: success")
}

func testUserDetele(t *testing.T, r Redmine, id int) {

	_, err := r.UserDelete(id)
	if err != nil {
		t.Fatal("User delete error:", err)
	}

	t.Logf("User delete: success")
}

func testUserMultiGet(t *testing.T, r Redmine) {

	u, _, err := r.UserMultiGet(UserStatusActive, "", 0)
	if err != nil {
		t.Fatal("Users get error:", err)
	}

	for _, e := range u {
		if e.Login == testUserLogin {
			t.Logf("Users get: success")
			return
		}
	}

	t.Fatal("Users get error: can't find created user")
}

func testUserSingleGet(t *testing.T, r Redmine, id int) {

	_, _, err := r.UserSingleGet(id, []string{"groups", "memberships"})
	if err != nil {
		t.Fatal("User get error:", err)
	}

	t.Logf("User get: success")
}
