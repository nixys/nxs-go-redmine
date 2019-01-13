package redmine

import (
	"os"
	"testing"
)

func initTest(r *Context, t *testing.T) {

	rdmnHost := os.Getenv("REDMINE_HOST")
	if rdmnHost == "" {
		t.Fatal("Init error: undefined env var `REDMINE_HOST`")
	}

	rdmnAPIKey := os.Getenv("REDMINE_API_KEY")
	if rdmnAPIKey == "" {
		t.Fatal("Init error: undefined env var `REDMINE_API_KEY`")
	}

	r.SetEndpoint(rdmnHost)
	r.SetAPIKey(rdmnAPIKey)
	r.SetLimit(100)

	t.Logf("Init: success")
}
