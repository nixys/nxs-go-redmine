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

	testIssueNote        = "Test issue note"
	testIssuePrivateNote = "Test issue private note"

	testIssueStartDate = "2022-07-01"
	testIssueDueDate   = "2022-07-02"

	testIssueStartDate2 = "2022-07-03"
	testIssueDueDate2   = "2022-07-04"
)

func TestIssuesCRUD(t *testing.T) {

	var r Context

	// Get env variables
	testIssueTrackerID, err := strconv.ParseInt(os.Getenv("REDMINE_TRACKER_ID"), 10, 64)
	if err != nil {
		t.Fatal("Issue test error: env variable `REDMINE_TRACKER_ID` is incorrect")
	}

	testMembershipRoleID, err := strconv.ParseInt(os.Getenv("REDMINE_ROLE_ID_1"), 10, 64)
	if err != nil {
		t.Fatal("Issue test error: env variable `REDMINE_ROLE_ID_1` is incorrect")
	}

	if testIssueTrackerID == 0 ||
		testMembershipRoleID == 0 {
		t.Fatal("Issue test error: env variables `REDMINE_TRACKER_ID` or `REDMINE_ROLE_ID_1` does not set")
	}

	// Init Redmine context
	initTest(&r, t)

	// Preparing auxiliary data
	pCreated := testProjectCreate(t, r, []int64{testIssueTrackerID})
	defer testProjectDetele(t, r, pCreated.Identifier)

	uCreated := testUserCreate(t, r)
	defer testUserDetele(t, r, uCreated.ID)

	testMembershipAdd(t, r, pCreated.Identifier, uCreated.ID, testMembershipRoleID)

	// Add and delete
	iCreated := testIssueCreate(t, r, pCreated.ID, uCreated.ID, nil)
	defer testIssueDetele(t, r, iCreated.ID)

	// Signle get
	testIssueSingleGet(t, r, iCreated.ID, uCreated.ID)

	// Update
	testIssueNoteAdd(t, r, iCreated.ID, testIssueNote, false)
	testIssueNoteAdd(t, r, iCreated.ID, testIssuePrivateNote, true)
	testIssueUpdate(t, r, iCreated.ID)

	// Get multi
	testIssueMultiGet(t, r, iCreated.ID)

	// Get all
	testIssueAllGet(t, r, iCreated.ID)

	// Watchers delete and add
	testIssueWatcherDelete(t, r, iCreated.ID, uCreated.ID)
	testIssueWatcherAdd(t, r, iCreated.ID, uCreated.ID)
}

func testIssueCreate(t *testing.T, r Context, projectID, userID int64, upload *AttachmentUploadObject) IssueObject {

	var (
		u []AttachmentUploadObject
		w []int64
	)

	if upload != nil {
		u = append(u, *upload)
	}

	if userID > 0 {
		w = append(w, userID)
	}

	i, s, err := r.IssueCreate(
		IssueCreate{
			Issue: IssueCreateObject{
				ProjectID:      projectID,
				Subject:        testIssueSubject,
				WatcherUserIDs: w,
				Description:    testIssueDescription,
				StartDate:      testIssueStartDate,
				DueDate:        testIssueDueDate,
				Uploads:        u,
			},
		},
	)
	if err != nil {
		t.Fatal("Issue create error:", err, s)
	}

	o, s, err := r.IssueSingleGet(i.ID, IssueSingleGetRequest{})
	if err != nil {
		t.Fatal("Issue create error:", err, s)
	}

	if o.StartDate != testIssueStartDate || o.DueDate != testIssueDueDate {
		t.Fatal("Issue create error: incorrect issue start or due date")
	}

	t.Logf("Issue create: success")

	return i
}

func testIssueUpdate(t *testing.T, r Context, id int64) {

	s, err := r.IssueUpdate(
		id,
		IssueUpdate{
			Issue: IssueUpdateObject{
				Subject:     testIssueSubject2,
				Description: testIssueDescription2,
				StartDate:   &testIssueStartDate2,
				DueDate:     &testIssueDueDate2,
				IsPrivate:   true,
			},
		},
	)
	if err != nil {
		t.Fatal("Issue update error:", err, s)
	}

	o, s, err := r.IssueSingleGet(id, IssueSingleGetRequest{})
	if err != nil {
		t.Fatal("Issue update error:", err, s)
	}

	if o.StartDate != testIssueStartDate2 || o.DueDate != testIssueDueDate2 {
		t.Fatal("Issue update error: incorrect issue start or due date")
	}

	t.Logf("Issue update: success")
}

func testIssueNoteAdd(t *testing.T, r Context, id int64, notes string, privateNotes bool) {

	s, err := r.IssueUpdate(
		id,
		IssueUpdate{
			Issue: IssueUpdateObject{
				Notes:        notes,
				PrivateNotes: privateNotes,
			},
		},
	)
	if err != nil {
		t.Fatal("Issue notes add error:", err, s)
	}

	o, s, err := r.IssueSingleGet(id, IssueSingleGetRequest{
		Includes: []string{"journals"},
	})
	if err != nil {
		t.Fatal("Issue notes add error:", err, s)
	}

	if len(o.Journals) == 0 {
		t.Fatal("Issue notes add error: bad journals count")
	}

	j := o.Journals[len(o.Journals)-1]

	if j.Notes != notes || j.PrivateNotes != privateNotes {
		t.Fatal("Issue notes add error: incorrect comment text or notes privacy")
	}

	t.Logf("Issue notes add: success")
}

func testIssueDetele(t *testing.T, r Context, id int64) {

	_, err := r.IssueDelete(id)
	if err != nil {
		t.Fatal("Issue delete error:", err)
	}

	t.Logf("Issue delete: success")
}

func testIssueAllGet(t *testing.T, r Context, id int64) {

	i, s, err := r.IssuesAllGet(IssueAllGetRequest{
		Includes: []string{"relations", "attachments"},
		Filters: IssueGetRequestFilters{
			Fields: map[string][]string{
				"issue_id":    {strconv.FormatInt(id, 10)},
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

func testIssueMultiGet(t *testing.T, r Context, id int64) {

	i, s, err := r.IssuesMultiGet(IssueMultiGetRequest{
		Includes: []string{"relations", "attachments"},
		Filters: IssueGetRequestFilters{
			Fields: map[string][]string{
				"issue_id":    {strconv.FormatInt(id, 10)},
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

func testIssueSingleGet(t *testing.T, r Context, id, userID int64) {

	i, s, err := r.IssueSingleGet(id, IssueSingleGetRequest{
		Includes: []string{"children", "attachments", "relations", "changesets", "journals", "watchers"},
	})
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

func testIssueWatcherAdd(t *testing.T, r Context, id int64, userID int64) {

	s, err := r.IssueWatcherAdd(id, userID)
	if err != nil {
		t.Fatal("Issue add watcher error:", err, s)
	}

	t.Logf("Issue add watcher: success")
}

func testIssueWatcherDelete(t *testing.T, r Context, id int64, userID int64) {

	s, err := r.IssueWatcherDelete(id, userID)
	if err != nil {
		t.Fatal("Issue delete watcher error:", err, s)
	}

	t.Logf("Issue delete watcher: success")
}
