# nxs-go-redmine

This Go package provides access to Redmine API.

Compatible with Redmine 4.2+

Follows Redmine resources are fully implemented at this moment:
- [Issues](https://www.redmine.org/projects/redmine/wiki/Rest_Issues)
- [Projects](https://www.redmine.org/projects/redmine/wiki/Rest_Projects)
- [Project Memberships](https://www.redmine.org/projects/redmine/wiki/Rest_Memberships)
- [Users](https://www.redmine.org/projects/redmine/wiki/Rest_Users)
- [Attachments](https://www.redmine.org/projects/redmine/wiki/Rest_Attachments)
- [Issue Statuses](https://www.redmine.org/projects/redmine/wiki/Rest_IssueStatuses)
- [Trackers](https://www.redmine.org/projects/redmine/wiki/Rest_Trackers)
- [Enumerations](https://www.redmine.org/projects/redmine/wiki/Rest_Enumerations)
- [Groups](https://www.redmine.org/projects/redmine/wiki/Rest_Groups)
- [Custom Fields](https://www.redmine.org/projects/redmine/wiki/Rest_CustomFields)

## Install

```
go get github.com/nixys/nxs-go-redmine/v4
```

## Example of usage

*You may find more examples in unit-tests in this repository*

**Get all projects from the Redmine server:**

```go
package main

import (
	"fmt"
	"os"

	"github.com/nixys/nxs-go-redmine/v4"
)

func main() {

	var r redmine.Context

	// Get variables from environment for connect to Redmine server 
	rdmnHost := os.Getenv("REDMINE_HOST")
	rdmnAPIKey := os.Getenv("REDMINE_API_KEY")
	if rdmnHost == "" || rdmnAPIKey == "" {
		fmt.Println("Init error: make sure environment variables `REDMINE_HOST` and `REDMINE_API_KEY` are defined")
		os.Exit(1)
	}

	// Init Redmine ctx 
	r.SetEndpoint(rdmnHost)
	r.SetAPIKey(rdmnAPIKey)

	fmt.Println("Init: success")

	// Get all projects 
	p, _, err := r.ProjectAllGet(redmine.ProjectAllGetRequest{
		Includes: []string{"trackers", "issue_categories", "enabled_modules"},
		Filters: redmine.ProjectGetRequestFilters{
			Status: redmine.ProjectStatusActive,
		},
	})
	if err != nil {
		fmt.Println("Projects get error:", err)
		os.Exit(1)
	}

	fmt.Println("Projects:")
	for _, e := range p.Projects {
		fmt.Println("-", e.Name)
	}
}
```

Run:

```
REDMINE_HOST="https://redmine.yourdomain.com" REDMINE_API_KEY="YOUR_API_KEY" go run main.go
```
