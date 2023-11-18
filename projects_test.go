package redmine

import (
	"os"
	"strconv"
	"testing"
)

var (
	testProjectName       = "test-project"
	testProjectName2      = "test-project2"
	testProjectIdentifier = "test_project"
)

func TestProjectsCRUDIdentifier(t *testing.T) {

	var r Context

	// Get env variables
	testProjectTrackerID, err := strconv.ParseInt(os.Getenv("REDMINE_TRACKER_ID"), 10, 64)
	if err != nil {
		t.Fatal("Project test error: env variable `REDMINE_TRACKER_ID` is incorrect")
	}

	if testProjectTrackerID == 0 {
		t.Fatal("Project test error: env variables `REDMINE_TRACKER_ID` does not set")
	}

	// Init Redmine context
	initTest(&r, t)

	// Create and delete
	pCreated := testProjectCreate(t, r, []int64{testProjectTrackerID})
	defer testProjectDetele(t, r, pCreated.Identifier)

	// Archive and unarcheve
	testProjectArchive(t, r, pCreated.Identifier)
	testProjectUnarchive(t, r, pCreated.Identifier)

	// Get
	testProjectAllGet(t, r)
	testProjectSingleGet(t, r, pCreated.Identifier)

	// Update
	testProjectUpdate(t, r, pCreated.Identifier)
}

func TestProjectsCRUDID(t *testing.T) {

	var r Context

	// Get env variables
	testProjectTrackerID, err := strconv.ParseInt(os.Getenv("REDMINE_TRACKER_ID"), 10, 64)
	if err != nil {
		t.Fatal("Project test error: env variable `REDMINE_TRACKER_ID` is incorrect")
	}

	if testProjectTrackerID == 0 {
		t.Fatal("Project test error: env variables `REDMINE_TRACKER_ID` does not set")
	}

	// Init Redmine context
	initTest(&r, t)

	// Create and delete
	pCreated := testProjectCreate(t, r, []int64{testProjectTrackerID})
	defer testProjectDetele(t, r, strconv.FormatInt(pCreated.ID, 10))

	// Get
	testProjectAllGet(t, r)
	testProjectSingleGet(t, r, strconv.FormatInt(pCreated.ID, 10))

	// Update
	testProjectUpdate(t, r, strconv.FormatInt(pCreated.ID, 10))
}

func testProjectCreate(t *testing.T, r Context, trackerIDs []int64) ProjectObject {

	p, _, err := r.ProjectCreate(
		ProjectCreate{
			Project: ProjectCreateObject{
				Name:       testProjectName,
				Identifier: testProjectIdentifier,
				TrackerIDs: &trackerIDs,
			},
		},
	)
	if err != nil {
		t.Fatal("Project create error:", err)
	}

	t.Logf("Project create: success")

	return p
}

func testProjectUpdate(t *testing.T, r Context, id string) {

	_, err := r.ProjectUpdate(
		id,
		ProjectUpdate{
			Project: ProjectUpdateObject{
				Name: &testProjectName2,
			},
		},
	)
	if err != nil {
		t.Fatal("Project update error:", err)
	}

	t.Logf("Project update: success")
}

func testProjectDetele(t *testing.T, r Context, id string) {

	_, err := r.ProjectDelete(id)
	if err != nil {
		t.Fatal("Project delete error:", err)
	}

	t.Logf("Project delete: success")
}

func testProjectAllGet(t *testing.T, r Context) {

	p, _, err := r.ProjectAllGet(
		ProjectAllGetRequest{
			Includes: []ProjectInclude{
				ProjectIncludeTrackers,
				ProjectIncludeIssueCategories,
				ProjectIncludeEnabledModules,
				ProjectIncludeTimeEntryActivities,
				ProjectIncludeIssueCustomFields,
			},
			Filters: ProjectGetRequestFilters{
				Status: ProjectStatusActive,
			},
		},
	)
	if err != nil {
		t.Fatal("Projects get error:", err)
	}

	for _, e := range p.Projects {
		if e.Name == testProjectName {
			t.Logf("Projects all get: success")
			return
		}
	}

	t.Fatal("Projects get error: can't find created project")
}

func testProjectSingleGet(t *testing.T, r Context, id string) {

	_, _, err := r.ProjectSingleGet(
		id,
		ProjectSingleGetRequest{
			Includes: []ProjectInclude{
				ProjectIncludeTrackers,
				ProjectIncludeIssueCategories,
				ProjectIncludeEnabledModules,
				ProjectIncludeTimeEntryActivities,
				ProjectIncludeIssueCustomFields,
			},
		},
	)
	if err != nil {
		t.Fatal("Project get error:", err)
	}

	t.Logf("Project get: success")
}

func testProjectArchive(t *testing.T, r Context, id string) {

	_, err := r.ProjectArchive(id)
	if err != nil {
		t.Fatal("Project archive error:", err)
	}

	t.Logf("Project archive: success")
}

func testProjectUnarchive(t *testing.T, r Context, id string) {

	_, err := r.ProjectUnarchive(id)
	if err != nil {
		t.Fatal("Project unarchive error:", err)
	}

	t.Logf("Project unarchive: success")
}
