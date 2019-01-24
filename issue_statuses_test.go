package redmine

import (
	"fmt"
	"testing"
)

func TestIssueStatusesCRUD(t *testing.T) {

	var r Context

	// Init Redmine context
	initTest(&r, t)

	// Get
	testIssueStatusMultiGet(t, r)
}

func testIssueStatusMultiGet(t *testing.T, r Context) {

	is, _, err := r.IssueStatusMultiGet()
	if err != nil {
		t.Fatal("Issue statuses get error:", err)
	}

	fmt.Println(is)

	if len(is) > 0 {
		t.Logf("Issue statuses get: success")
		return
	}

	t.Fatal("Issue statuses get error: can't find any issue statuses")
}
