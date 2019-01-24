package redmine

import (
	"testing"
)

func TestTrackerCRUD(t *testing.T) {

	var r Context

	// Init Redmine context
	initTest(&r, t)

	// Get
	testTrackerMultiGet(t, r)
}

func testTrackerMultiGet(t *testing.T, r Context) {

	tr, _, err := r.TrackerMultiGet()
	if err != nil {
		t.Fatal("Trackers get error:", err)
	}

	if len(tr) > 0 {
		t.Logf("Trackers get: success")
		return
	}

	t.Fatal("Trackers get error: can't find any trackers")
}
