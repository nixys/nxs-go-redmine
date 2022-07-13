package redmine

import (
	"os"
	"strconv"
	"testing"
)

const (
	testAttachmentFile         = "attachments_test.go"
	testAttachmentFileDownload = "/tmp/" + testAttachmentFile
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
	defer testProjectDetele(t, r, pCreated.Identifier)

	// Get multi

	// Upload
	aCreated := testAttachmentUpload(t, r, pCreated.ID, 0)

	// Get single
	testAttachmentGetSingle(t, r, aCreated)

	// Download
	testAttachmentDownload(t, r, aCreated)
}

func testAttachmentUpload(t *testing.T, r Context, projectID, userID int) int {

	u, s, err := r.AttachmentUpload(testAttachmentFile)
	if err != nil {
		t.Fatal("Upload attachment error:", err, s)
	}

	// Created issue will be deleted with the project
	i := testIssueCreate(t, r, projectID, userID, &u)

	// Request single issue to get Attachment ID
	j, s, err := r.IssueSingleGet(i.ID, IssueSingleGetRequest{
		Includes: []string{"attachments"},
	})
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

	a, s, err := r.AttachmentSingleGet(id)
	if err != nil {
		t.Fatal("Attachment get error:", err, s)
	}

	if a.FileName != testAttachmentFile {
		t.Fatal("Attachment get error: wrong attachment file name")
	}

	t.Logf("Attachment get: success")
}

func testAttachmentDownload(t *testing.T, r Context, id int) {

	a, s, err := r.AttachmentDownload(id, testAttachmentFileDownload)
	if err != nil {
		t.Fatal("Attachment download error:", err, s)
	}
	defer os.Remove(testAttachmentFileDownload)

	if a.FileName != testAttachmentFile {
		t.Fatal("Attachment download error: wrong attachment file name")
	}

	t.Logf("Attachment get: success")
}
