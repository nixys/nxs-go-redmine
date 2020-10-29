package redmine

import (
	"testing"
)

func TestEnumerationCRUD(t *testing.T) {

	var r Context

	// Init Redmine context
	initTest(&r, t)

	// Get enumeration priorities
	testEnumerationPrioritiesAllGet(t, r)

	// Get enumeration time entry activities
	testEnumerationTimeEntryActivitiesAllGet(t, r)

	// Get enumeration document categories
	testEnumerationDocumentCategoriesAllGet(t, r)
}

func testEnumerationPrioritiesAllGet(t *testing.T, r Context) {

	e, _, err := r.EnumerationPrioritiesAllGet()
	if err != nil {
		t.Fatal("Enumeration priorities get error:", err)
	}

	if len(e) == 0 {
		t.Fatal("Enumeration priorities get error: can't find any priorities")
	}

	t.Logf("Enumeration priorities: success")
}

func testEnumerationTimeEntryActivitiesAllGet(t *testing.T, r Context) {

	e, _, err := r.EnumerationTimeEntryActivitiesAllGet()
	if err != nil {
		t.Fatal("Enumeration time entry activities get error:", err)
	}

	if len(e) == 0 {
		t.Fatal("Enumeration time entry activities get error: can't find any time entry activities")
	}

	t.Logf("Enumeration time entry activities: success")
}

func testEnumerationDocumentCategoriesAllGet(t *testing.T, r Context) {

	e, _, err := r.EnumerationDocumentCategoriesAllGet()
	if err != nil {
		t.Fatal("Enumeration document categories get error:", err)
	}

	if len(e) == 0 {
		t.Fatal("Enumeration document categories get error: can't find any document categories")
	}

	t.Logf("Enumeration document categories: success")
}
