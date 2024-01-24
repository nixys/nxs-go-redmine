package redmine

import (
	"os"
	"strconv"
	"testing"
)

var (
	testTimeEntriesHours1  float64 = 1.23
	testTimeEntriesHours2  float64 = 1.24
	testTimeEntriesComment string  = "Comment"
)

func TestTimeEntriesCRUD(t *testing.T) {

	var r Context

	// Init Redmine context
	initTest(&r, t)

	// Get env variables
	testIssueTrackerID, err := strconv.ParseInt(os.Getenv("REDMINE_TRACKER_ID"), 10, 64)
	if err != nil {
		t.Fatal("Time entry test error: env variable `REDMINE_TRACKER_ID` is incorrect")
	}
	if testIssueTrackerID == 0 {
		t.Fatal("Time entry test error: env variable `REDMINE_TRACKER_ID` does not set")
	}

	testTimeEntriesActivityID, err := strconv.ParseInt(os.Getenv("REDMINE_ACTIVITY_ID"), 10, 64)
	if err != nil {
		t.Fatal("Time entry test error: env variable `REDMINE_ACTIVITY_ID` is incorrect")
	}
	if testTimeEntriesActivityID == 0 {
		t.Fatal("Time entry test error: env variable `REDMINE_ACTIVITY_ID` does not set")
	}

	// Preparing auxiliary data
	pCreated := testProjectCreate(t, r, []int64{testIssueTrackerID})
	defer testProjectDetele(t, r, pCreated.Identifier)
	iCreated := testIssueCreate(t, r, pCreated.ID, 0, nil)

	// Create and delete
	teCreated := testTimeEntryCreate(t, r, iCreated.ID, testTimeEntriesActivityID)
	defer testTimeEntryDetele(t, r, teCreated.ID)

	// Get
	testTimeEntriesAllGet(t, r, pCreated.Identifier)
	testTimeEntrySingleGet(t, r, teCreated.ID)

	// Update
	testTimeEntryUpdate(t, r, teCreated.ID)
}

func testTimeEntryCreate(t *testing.T, r Context, issueID, activityID int64) TimeEntryObject {

	te, _, err := r.TimeEntryCreate(
		TimeEntryCreate{
			TimeEntry: TimeEntryCreateObject{
				IssueID:    &issueID,
				ActivityID: activityID,
				Hours:      testTimeEntriesHours1,
				Comments:   testTimeEntriesComment,
			},
		},
	)
	if err != nil {
		t.Fatal("Time entry create error:", err)
	}

	t.Logf("Time entry create: success")

	return te
}

func testTimeEntryUpdate(t *testing.T, r Context, id int64) {

	_, err := r.TimeEntryUpdate(
		id,
		TimeEntryUpdate{
			TimeEntry: TimeEntryUpdateObject{
				Hours: &testTimeEntriesHours2,
			},
		},
	)
	if err != nil {
		t.Fatal("Time entry update error:", err)
	}

	s, _, err := r.TimeEntrySingleGet(id, TimeEntrySingleGetRequest{})
	if err != nil {
		t.Fatal("Time entry update error:", err)
	}

	if s.Hours != testTimeEntriesHours2 {
		t.Fatal("Time entry update error: incorrect hours")
	}

	if s.Comments != testTimeEntriesComment {
		t.Fatal("Time entry update error: incorrect comment")
	}

	t.Logf("Time entry update: success")
}

func testTimeEntryDetele(t *testing.T, r Context, id int64) {

	_, err := r.TimeEntryDelete(id)
	if err != nil {
		t.Fatal("Time entry delete error:", err)
	}

	t.Logf("Time entry delete: success")
}

func testTimeEntriesAllGet(t *testing.T, r Context, projectID string) {

	te, _, err := r.TimeEntryAllGet(
		TimeEntryAllGetRequest{
			Filters: TimeEntryGetRequestFiltersInit().
				ProjectSet(projectID),
		},
	)
	if err != nil {
		t.Fatal("Time entries get all error:", err)
	}

	if len(te.TimeEntries) == 0 {
		t.Fatal("Time entries get all error: can't find any time entries")
	}

	t.Logf("Time entries get all: success")
}

func testTimeEntrySingleGet(t *testing.T, r Context, id int64) {

	s, _, err := r.TimeEntrySingleGet(id, TimeEntrySingleGetRequest{})
	if err != nil {
		t.Fatal("Time entry single get error:", err)
	}

	if s.Hours != testTimeEntriesHours1 {
		t.Fatal("Time entry single get error: incorrect hours")
	}

	if s.Comments != testTimeEntriesComment {
		t.Fatal("Time entry single get error: incorrect comment")
	}

	t.Logf("Time entry single get: success")
}
