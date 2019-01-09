package redmine

import (
	"testing"
)

const (
	testProjectName       = "test-project"
	testProjectName2      = "test-project2"
	testProjectIdentifier = "test_project"
)

func TestProjectsCRUD(t *testing.T) {

	var r Redmine

	/* Init Redmine context */
	initTest(&r, t)

	/* Create and delete */
	pCreated := testProjectCreate(t, r)
	defer testProjectDetele(t, r, pCreated.ID)

	/* Get */
	testProjectMultiGet(t, r)
	testProjectSingleGet(t, r, pCreated.ID)

	/* Update */
	testProjectUpdate(t, r, pCreated.ID)
}

func testProjectCreate(t *testing.T, r Redmine) ProjectObject {

	p, _, err := r.ProjectCreate(ProjectCreateObject{
		Name:           testProjectName,
		Identifier:     testProjectIdentifier,
		IsPublic:       false,
		InheritMembers: false,
	})
	if err != nil {
		t.Fatal("Project create error:", err)
	}

	t.Logf("Project create: success")

	return p
}

func testProjectUpdate(t *testing.T, r Redmine, id int) {

	_, err := r.ProjectUpdate(id, ProjectUpdateObject{
		Name: testProjectName2,
	})
	if err != nil {
		t.Fatal("Project update error:", err)
	}

	t.Logf("Project update: success")
}

func testProjectDetele(t *testing.T, r Redmine, id int) {

	_, err := r.ProjectDelete(id)
	if err != nil {
		t.Fatal("Project delete error:", err)
	}

	t.Logf("Project delete: success")
}

func testProjectMultiGet(t *testing.T, r Redmine) {

	p, _, err := r.ProjectMultiGet([]string{"trackers", "issue_categories", "enabled_modules"})
	if err != nil {
		t.Fatal("Projects get error:", err)
	}

	for _, e := range p {
		if e.Name == testProjectName {
			t.Logf("Projects get: success")
			return
		}
	}

	t.Fatal("Projects get error: can't find created project")
}

func testProjectSingleGet(t *testing.T, r Redmine, id int) {

	_, _, err := r.ProjectSingleGet(id, []string{"trackers", "issue_categories", "enabled_modules"})
	if err != nil {
		t.Fatal("Project get error:", err)
	}

	t.Logf("Project get: success")
}