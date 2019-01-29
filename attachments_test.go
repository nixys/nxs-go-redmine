package redmine

import (
	"os"
	"strconv"
	"testing"
)

const (
	testAttachmentFileUpload = "attachments_test.go"
)

func TestAttachmentsCRUD(t *testing.T) {

	var r Context

	// Get env variables
	testIssueTrackerID, _ := strconv.Atoi(os.Getenv("REDMINE_TRACKER_ID"))

	if testIssueTrackerID == 0 {
		t.Fatal("Attachments test error: env variable `REDMINE_TRACKER_ID` does not set")
	}

	// Init Redmine context
	initTest(&r, t)

	// Preparing auxiliary data
	pCreated := testProjectCreate(t, r, []int{testIssueTrackerID})
	defer testProjectDetele(t, r, pCreated.ID)

	// Get multi

	// Upload
	aCreated := testAttachmentUpload(t, r, pCreated.ID, 0)

	// Get single
	testAttachmentGetSingle(t, r, aCreated)

}

func testAttachmentUpload(t *testing.T, r Context, projectID, userID int) int {

	u, s, err := r.AttachmentUpload(testAttachmentFileUpload)
	if err != nil {
		t.Fatal("Upload attachment error:", err, s)
	}

	// Created issue will be deleted with the project
	i := testIssueCreate(t, r, projectID, userID, &u)

	// Request single issue to get Attachment ID
	j, s, err := r.IssueSingleGet(i.ID, []string{"attachments"})
	if err != nil {
		t.Fatal("Issue get error:", err, s)
	}

	if len(j.Attachments) != 1 {
		t.Fatal("Upload attachment error: wrong attachments count")
	}

	t.Logf("Upload attachment and create issue: success")

	return j.Attachments[0].ID
}

func testAttachmentGetSingle(t *testing.T, r Context, id int) {

	_, s, err := r.AttachmentSingleGet(id)
	if err != nil {
		t.Fatal("Attachment get error:", err, s)
	}

	t.Logf("Attachment get: success")
}
