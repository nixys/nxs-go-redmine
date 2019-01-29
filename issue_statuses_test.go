package redmine

import (
	"testing"
)

func TestIssueStatusesCRUD(t *testing.T) {

	var r Context

	// Init Redmine context
	initTest(&r, t)

	// Get
	testIssueStatusAllGet(t, r)
}

func testIssueStatusAllGet(t *testing.T, r Context) {

	is, _, err := r.IssueStatusAllGet()
	if err != nil {
		t.Fatal("Issue statuses get error:", err)
	}

	if len(is) > 0 {
		t.Logf("Issue statuses get: success")
		return
	}

	t.Fatal("Issue statuses get error: can't find any issue statuses")
}
