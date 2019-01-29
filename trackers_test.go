package redmine

import (
	"fmt"
	"testing"
)

func TestTrackersCRUD(t *testing.T) {

	var r Context

	// Init Redmine context
	initTest(&r, t)

	// Get
	testTrackerAllGet(t, r)
}

func testTrackerAllGet(t *testing.T, r Context) {

	tr, _, err := r.TrackerAllGet()
	if err != nil {
		t.Fatal("Trackers get error:", err)
	}

	if len(tr) == 0 {
		t.Fatal("Trackers get error: can't find any trackers")
	}

	fmt.Println(tr)

	t.Logf("Trackers get: success")
}
