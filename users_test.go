package redmine

import (
	"testing"
)

var (
	testUserLogin      = "test-user-login"
	testUserFirstName  = "First"
	testUserLastName   = "Last"
	testUserFirstName2 = "First2"
	testUserLastName2  = "Last2"
	testUserMail       = "test@domain.com"
)

func TestUserCRUD(t *testing.T) {

	var r Context

	// Init Redmine context
	initTest(&r, t)

	// Create and delete
	uCreated := testUserCreate(t, r)
	defer testUserDetele(t, r, uCreated.ID)

	// Get
	testUserAllGet(t, r)
	testUserSingleGet(t, r, uCreated.ID)

	// Update
	testUserUpdate(t, r, uCreated.ID)

	// Current
	testUserCurrentGet(t, r)
}

func testUserCreate(t *testing.T, r Context) UserObject {

	u, _, err := r.UserCreate(
		UserCreate{
			User: UserCreateObject{
				Login:            testUserLogin,
				FirstName:        testUserFirstName,
				LastName:         testUserLastName,
				Mail:             testUserMail,
				MailNotification: StringPtr(UserNotificationOnlyAssigned.String()),
				MustChangePasswd: BoolPtr(true),
				GeneratePassword: BoolPtr(true),
			},
			SendInformation: BoolPtr(true),
		},
	)
	if err != nil {
		t.Fatal("User create error:", err)
	}

	t.Logf("User create: success")

	return u
}

func testUserUpdate(t *testing.T, r Context, id int64) {

	_, err := r.UserUpdate(
		id,
		UserUpdate{
			User: UserUpdateObject{
				FirstName:        &testUserFirstName2,
				LastName:         &testUserLastName2,
				MailNotification: StringPtr(UserNotificationOnlyNone.String()),
			},
			SendInformation: BoolPtr(true),
		},
	)
	if err != nil {
		t.Fatal("User update error:", err)
	}

	t.Logf("User update: success")
}

func testUserDetele(t *testing.T, r Context, id int64) {

	_, err := r.UserDelete(id)
	if err != nil {
		t.Fatal("User delete error:", err)
	}

	t.Logf("User delete: success")
}

func testUserAllGet(t *testing.T, r Context) {

	u, _, err := r.UserAllGet(UserAllGetRequest{
		Filters: UserGetRequestFiltersInit().
			StatusSet(UserStatusActive),
	})
	if err != nil {
		t.Fatal("Users get error:", err)
	}

	for _, e := range u.Users {
		if e.Login == testUserLogin {
			t.Logf("Users get: success")
			return
		}
	}

	t.Fatal("Users get error: can't find created user")
}

func testUserSingleGet(t *testing.T, r Context, id int64) {

	_, _, err := r.UserSingleGet(id, UserSingleGetRequest{
		Includes: []UserInclude{
			UserIncludeGroups,
			UserIncludeMemberships,
		},
	})
	if err != nil {
		t.Fatal("User get error:", err)
	}

	t.Logf("User get: success")
}

func testUserCurrentGet(t *testing.T, r Context) {

	_, _, err := r.UserCurrentGet(UserCurrentGetRequest{
		Includes: []UserInclude{
			UserIncludeGroups,
			UserIncludeMemberships,
		},
	})
	if err != nil {
		t.Fatal("Current user get error:", err)
	}

	t.Logf("Current user get: success")
}
