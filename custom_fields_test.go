package redmine

import (
	"testing"
)

func TestCustomFieldsCRUD(t *testing.T) {

	var r Context

	// Init Redmine context
	initTest(&r, t)

	// Get
	testCustomFieldAllGet(t, r)
}

func testCustomFieldAllGet(t *testing.T, r Context) {

	cf, _, err := r.CustomFieldAllGet()
	if err != nil {
		t.Fatal("Custom fields get error:", err)
	}

	if len(cf) > 0 {
		t.Logf("Custom fields get: success")
		return
	}

	t.Fatal("Custom fields get error: can't find any custom fields")
}
