# nxs-go-redmine

## Introduction

Go client library for [Redmine](https://www.redmine.org)

### Features

- Compatible with Redmine 4.2+
- Implemented following Redmine resources:
  - [Issues](https://www.redmine.org/projects/redmine/wiki/Rest_Issues)
  - [Projects](https://www.redmine.org/projects/redmine/wiki/Rest_Projects)
  - [Project Memberships](https://www.redmine.org/projects/redmine/wiki/Rest_Memberships)
  - [Users](https://www.redmine.org/projects/redmine/wiki/Rest_Users)
  - [Wiki Pages](https://www.redmine.org/projects/redmine/wiki/Rest_WikiPages)
  - [Attachments](https://www.redmine.org/projects/redmine/wiki/Rest_Attachments)
  - [Issue Statuses](https://www.redmine.org/projects/redmine/wiki/Rest_IssueStatuses)
  - [Trackers](https://www.redmine.org/projects/redmine/wiki/Rest_Trackers)
  - [Enumerations](https://www.redmine.org/projects/redmine/wiki/Rest_Enumerations)
  - [Groups](https://www.redmine.org/projects/redmine/wiki/Rest_Groups)
  - [Custom Fields](https://www.redmine.org/projects/redmine/wiki/Rest_CustomFields)

### Who can use the tool

Developer teams or sysadmins who need to automate a business processes that works around Redmine.

## Quickstart

### Import

```
go get github.com/nixys/nxs-go-redmine/v4
```

### Initialize

To initialize this library you need to do:
- Declare variable `redmine.Context`
- Set a Redmine endpoint via method `(r *Context) SetEndpoint(endpoint string)`
- Set a Redmine API key via method `(r *Context) SetAPIKey(apiKey string)`

After thar you be able to use all available methods to interact with Redmine API.

## Example

In the example below will be printed a names for all active projects from Redmine

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

**For more examples see apps based on this library:**
- [nxs-support-bot](https://github.com/nixys/nxs-support-bot)

## Feedback

For support and feedback please contact me:
- telegram: [@borisershov](https://t.me/borisershov)
- e-mail: b.ershov@nixys.ru

## License

nxs-go-redmine is released under the [MIT License](LICENSE).
