# nxs-go-redmine

This Go package provides access to Redmine API.

At this time not the all Redmine API methods are implemented, but work is in progress.

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

	var r redmine.Redmine

	/* Get variables from environment for connect to Redmine server */
	rdmnHost := os.Getenv("REDMINE_HOST")
	rdmnApiKey := os.Getenv("REDMINE_API_KEY")
	if rdmnHost == "" || rdmnApiKey == "" {
		fmt.Println("Init error: make sure environment variables `REDMINE_HOST` and `REDMINE_API_KEY` are defined")
		os.Exit(1)
	}

	/* Init Redmine ctx */
	r.SetEndpoint(rdmnHost)
	r.SetApiKey(rdmnApiKey)
	r.SetLimit(100)

	fmt.Println("Init: success")

	/* Get all projects */
	p, _, err := r.ProjectMultiGet([]string{"trackers", "issue_categories", "enabled_modules"})
	if err != nil {
		fmt.Println("Projects get error:", err)
		os.Exit(1)
	}

	fmt.Println("Projects:")
	for _, e := range p {
		fmt.Println("-", e.Name)
	}
}
```

Run:

```
REDMINE_HOST="https://redmine.yourdomain.com" REDMINE_API_KEY="YOUR_API_KEY" go run main.go
```
