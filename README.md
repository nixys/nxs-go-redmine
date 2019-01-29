# nxs-go-redmine

This Go package provides access to Redmine API.

Follows Redmine resources are fully implemented at this moment:
- [Issues](http://www.redmine.org/projects/redmine/wiki/Rest_Issues)
- [Projects](http://www.redmine.org/projects/redmine/wiki/Rest_Projects)
- [Project Memberships](http://www.redmine.org/projects/redmine/wiki/Rest_Memberships)
- [Users](http://www.redmine.org/projects/redmine/wiki/Rest_Users)
- [Attachments](http://www.redmine.org/projects/redmine/wiki/Rest_Attachments)
- [Issue Statuses](http://www.redmine.org/projects/redmine/wiki/Rest_IssueStatuses)
- [Trackers](http://www.redmine.org/projects/redmine/wiki/Rest_Trackers)
- [Groups](http://www.redmine.org/projects/redmine/wiki/Rest_Groups)
- [Custom Fields](http://www.redmine.org/projects/redmine/wiki/Rest_CustomFields)

## Install

```
go get github.com/nixys/nxs-go-redmine
```

## Example of usage

*You may find more examples in unit-tests in this repository*

**Get all projects from the Redmine server:**

```go
package main

import (
	"fmt"
	"os"

	"github.com/nixys/nxs-go-redmine"
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
	p, _, err := r.ProjectAllGet([]string{"trackers", "issue_categories", "enabled_modules"})
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
