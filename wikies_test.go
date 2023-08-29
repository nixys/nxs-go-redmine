package redmine

import (
	"os"
	"strconv"
	"testing"
)

var (
	testWikiTitle   = "TestTitle"
	testWikiText    = "TestText"
	testWikiComment = "TestComment"

	testWikiTextUpdated    = "TestTextUpdated"
	testWikiCommentUpdated = "TestCommentUpdated"
)

func TestWikiesCRUD(t *testing.T) {

	var r Context

	// Get env variables
	testIssueTrackerID, _ := strconv.Atoi(os.Getenv("REDMINE_TRACKER_ID"))

	if testIssueTrackerID == 0 {
		t.Fatal("Issue test error: env variable `REDMINE_TRACKER_ID` does not set")
	}

	// Init Redmine context
	initTest(&r, t)

	// Preparing auxiliary data
	pCreated := testProjectCreate(t, r, []int{testIssueTrackerID})
	defer testProjectDetele(t, r, pCreated.Identifier)

	// Add and delete
	testWikiCreate(t, r, pCreated.Identifier, testWikiTitle)
	defer testWikiDetele(t, r, pCreated.Identifier, testWikiTitle)

	// Get all
	testwikiAllGet(t, r, pCreated.Identifier)

	// Single get
	testWikiSingleGet(t, r, pCreated.Identifier, testWikiTitle)

	// Update
	testWikiUpdate(t, r, pCreated.Identifier, testWikiTitle)

	// Single version get
	testWikiSingleVersionGet(t, r, pCreated.Identifier, testWikiTitle, 2)
}

func testWikiCreate(t *testing.T, r Context, projectID, wikiTitle string) WikiObject {

	u, s, err := r.AttachmentUpload(testAttachmentFile)
	if err != nil {
		t.Fatal("Wiki create error:", err, s)
	}

	w, s, err := r.WikiCreate(
		projectID,
		wikiTitle,
		WikiCreateObject{
			Text:     testWikiText,
			Comments: testWikiComment,
			Uploads:  []AttachmentUploadObject{u},
		})
	if err != nil {
		t.Fatal("Wiki create error:", err, s)
	}

	if w.Title != wikiTitle {
		t.Fatal("Wiki get error: incorrect title")
	}

	t.Logf("Wiki create: success")

	return w
}

func testwikiAllGet(t *testing.T, r Context, projectID string) {

	w, s, err := r.WikiAllGet(projectID)
	if err != nil {
		t.Fatal("Wikies all get error:", err, s)
	}

	if len(w) == 0 {
		t.Fatal("Wikies all get error: can't find any wikies satisfies specified filters")
	}

	t.Logf("Wikies all get: success")
}

func testWikiSingleGet(t *testing.T, r Context, projectID, wikiTitle string) {

	w, s, err := r.WikiSingleGet(
		projectID,
		wikiTitle,
		WikiSingleGetRequest{
			Includes: []string{"attachments"},
		})
	if err != nil {
		t.Fatal("Wiki get error:", err, s)
	}

	if w.Title != wikiTitle {
		t.Fatal("Wiki get error: incorrect title")
	}

	if len(*w.Attachments) == 0 {
		t.Fatal("Wiki get error: incorrect attachments count")
	}

	if (*w.Attachments)[0].FileName != testAttachmentFile {
		t.Fatal("Wiki get error: incorrect attachments name")
	}

	t.Logf("Wiki get: success")
}

func testWikiSingleVersionGet(t *testing.T, r Context, projectID, wikiTitle string, version int) {

	w, s, err := r.WikiSingleVersionGet(
		projectID,
		wikiTitle,
		version,
		WikiSingleGetRequest{
			Includes: []string{"attachments"},
		})
	if err != nil {
		t.Fatal("Wiki version get error:", err, s)
	}

	if w.Title != wikiTitle {
		t.Fatal("Wiki version get error: incorrect title")
	}

	if w.Text != testWikiTextUpdated {
		t.Fatal("Wiki version get error: incorrect text")
	}

	t.Logf("Wiki version get: success")
}

func testWikiUpdate(t *testing.T, r Context, projectID, wikiTitle string) {

	s, err := r.WikiUpdate(
		projectID,
		wikiTitle,
		WikiUpdateObject{
			Text:     testWikiTextUpdated,
			Comments: testWikiCommentUpdated,
		})
	if err != nil {
		t.Fatal("Wiki update error:", err, s)
	}

	t.Logf("Wiki update: success")
}

func testWikiDetele(t *testing.T, r Context, projectID, wikiTitle string) {

	_, err := r.WikiDelete(projectID, wikiTitle)
	if err != nil {
		t.Fatal("Wiki delete error:", err)
	}

	t.Logf("Wiki delete: success")
}
