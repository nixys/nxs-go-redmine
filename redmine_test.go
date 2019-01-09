package redmine

import (
	"os"
	"testing"
)

func initTest(r *Redmine, t *testing.T) {

	rdmnHost := os.Getenv("REDMINE_HOST")
	if rdmnHost == "" {
		t.Fatal("Init error: undefined env var `REDMINE_HOST`")
	}

	rdmnApiKey := os.Getenv("REDMINE_API_KEY")
	if rdmnApiKey == "" {
		t.Fatal("Init error: undefined env var `REDMINE_API_KEY`")
	}

	r.SetEndpoint(rdmnHost)
	r.SetApiKey(rdmnApiKey)
	r.SetLimit(100)

	t.Logf("Init: success")
}
