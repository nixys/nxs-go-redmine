package redmine

import (
	"os"
	"strconv"
	"testing"
)

var (
	testIssueSubject     = "Test issue subject"
	testIssueDescription = "Test issue description"

	testIssueSubject2     = "Test issue subject2"
	testIssueDescription2 = "Test issue description2"
)

func TestIssuesCRUD(t *testing.T) {

	var r Context

	// Get env variables
	testIssueTrackerID, _ := strconv.Atoi(os.Getenv("REDMINE_TRACKER_ID"))
	testMembershipRoleID, _ := strconv.Atoi(os.Getenv("REDMINE_ROLE_ID_1"))

	if testIssueTrackerID == 0 ||
		testMembershipRoleID == 0 {
		t.Fatal("Issue test error: env variables `REDMINE_TRACKER_ID` or `REDMINE_ROLE_ID_1` does not set")
	}

	// Init Redmine context
	initTest(&r, t)

	// Preparing auxiliary data
	pCreated := testProjectCreate(t, r, []int{testIssueTrackerID})
	defer testProjectDetele(t, r, pCreated.ID)

	uCreated := testUserCreate(t, r)
	defer testUserDetele(t, r, uCreated.ID)

	testMembershipAdd(t, r, pCreated.ID, uCreated.ID, testMembershipRoleID)

	// Add and delete
	iCreated := testIssueCreate(t, r, pCreated.ID, uCreated.ID, nil)
	defer testIssueDetele(t, r, iCreated.ID)

	// Signle get
	testIssueSingleGet(t, r, iCreated.ID, uCreated.ID)

	// Update
	testIssueUpdate(t, r, iCreated.ID)

	// Get multi
	testIssueMultiGet(t, r, iCreated.ID)

	// Get all
	testIssueAllGet(t, r, iCreated.ID)

	// Watchers delete and add
	testIssueWatcherDelete(t, r, iCreated.ID, uCreated.ID)
	testIssueWatcherAdd(t, r, iCreated.ID, uCreated.ID)
}

func testIssueCreate(t *testing.T, r Context, projectID, userID int, upload *AttachmentUploadObject) IssueObject {

	var (
		u []AttachmentUploadObject
		w []int
	)

	if upload != nil {
		u = append(u, *upload)
	}

	if userID > 0 {
		w = append(w, userID)
	}

	i, s, err := r.IssueCreate(IssueCreateObject{
		ProjectID:      projectID,
		Subject:        testIssueSubject,
		WatcherUserIDs: w,
		Description:    testIssueDescription,
		Uploads:        u,
	})
	if err != nil {
		t.Fatal("Issue create error:", err, s)
	}

	t.Logf("Issue create: success")

	return i
}

func testIssueUpdate(t *testing.T, r Context, id int) {

	s, err := r.IssueUpdate(id, IssueUpdateObject{
		Subject:     testIssueSubject2,
		Description: testIssueDescription2,
		IsPrivate:   true,
	})
	if err != nil {
		t.Fatal("Issue update error:", err, s)
	}

	t.Logf("Issue update: success")
}

func testIssueDetele(t *testing.T, r Context, id int) {

	_, err := r.IssueDelete(id)
	if err != nil {
		t.Fatal("Issue delete error:", err)
	}

	t.Logf("Issue delete: success")
}

func testIssueAllGet(t *testing.T, r Context, id int) {

	i, s, err := r.IssuesAllGet(IssueAllGetRequest{
		Includes: []string{"relations", "attachments"},
		Filters: IssueGetRequestFilters{
			Fields: map[string][]string{
				"issue_id":    {strconv.Itoa(id)},
				"subject":     {testIssueSubject2},
				"description": {testIssueDescription2},
			},
		},
	})
	if err != nil {
		t.Fatal("Issues all get error:", err, s)
	}

	if len(i.Issues) == 0 {
		t.Fatal("Issues all get error: can't find any issues satisfies specified filters")
	}

	t.Logf("Issues all get: success")
}

func testIssueMultiGet(t *testing.T, r Context, id int) {

	i, s, err := r.IssuesMultiGet(IssueMultiGetRequest{
		Includes: []string{"relations", "attachments"},
		Filters: IssueGetRequestFilters{
			Fields: map[string][]string{
				"issue_id":    {strconv.Itoa(id)},
				"subject":     {testIssueSubject2},
				"description": {testIssueDescription2},
			},
		},
		Limit:  100,
		Offset: 0,
	})
	if err != nil {
		t.Fatal("Issues multi get error:", err, s)
	}

	if len(i.Issues) == 0 {
		t.Fatal("Issues multi get error: can't find any issues satisfies specified filters")
	}

	t.Logf("Issues multi get: success")
}

func testIssueSingleGet(t *testing.T, r Context, id, userID int) {

	i, s, err := r.IssueSingleGet(id, []string{"children", "attachments", "relations", "changesets", "journals", "watchers"})
	if err != nil {
		t.Fatal("Issue get error:", err, s)
	}

	if i.Subject != testIssueSubject {
		t.Fatal("Issue get error: incorrect subject")
	}

	if i.Description != testIssueDescription {
		t.Fatal("Issue get error: incorrect description")
	}

	if len(i.Watchers) != 1 {
		t.Fatal("Issue get error: incorrect watchers count")
	}

	if i.Watchers[0].ID != userID {
		t.Fatal("Issue get error: incorrect issue watchers")
	}

	t.Logf("Issue get: success")
}

func testIssueWatcherAdd(t *testing.T, r Context, id int, userID int) {

	s, err := r.IssueWatcherAdd(id, userID)
	if err != nil {
		t.Fatal("Issue add watcher error:", err, s)
	}

	t.Logf("Issue add watcher: success")
}

func testIssueWatcherDelete(t *testing.T, r Context, id int, userID int) {

	s, err := r.IssueWatcherDelete(id, userID)
	if err != nil {
		t.Fatal("Issue delete watcher error:", err, s)
	}

	t.Logf("Issue delete watcher: success")
}
