package redmine

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

// Redmine struct used for store settings to communicate with Redmine API
type Redmine struct {
	endpoint string
	apiKey   string
	limit    int
}

// IDName used as embedded struct for other structs within package
type IDName struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type errorsResult struct {
	Errors []string `json:"errors"`
}

// SetAPIKey is used to set Redmine API key in ctx
func (r *Redmine) SetAPIKey(apiKey string) {
	r.apiKey = apiKey
}

// SetEndpoint is used to set Redmine endpoint in ctx
func (r *Redmine) SetEndpoint(endpoint string) {
	r.endpoint = endpoint
}

// SetLimit is used to set elements limit on page
func (r *Redmine) SetLimit(limit int) {
	r.limit = limit
}

func (r *Redmine) get(out interface{}, uri string, statusExpected int) (int, error) {

	var er errorsResult

	url := r.endpoint + uri

	// Create request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}

	// Set headers
	req.Header.Add("X-Redmine-API-Key", r.apiKey)

	// Make request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}

	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	if res.StatusCode != statusExpected {
		err = decoder.Decode(&er)
		if err == nil {
			err = errors.New(strings.Join(er.Errors, "\n"))
		}
	} else {
		if out != nil {
			err = decoder.Decode(out)
		}
	}

	if err != nil {
		return res.StatusCode, err
	}

	return res.StatusCode, nil
}

func (r *Redmine) post(in interface{}, out interface{}, uri string, statusExpected int) (int, error) {

	return r.alter("POST", in, out, uri, statusExpected)
}

func (r *Redmine) put(in interface{}, out interface{}, uri string, statusExpected int) (int, error) {

	return r.alter("PUT", in, out, uri, statusExpected)
}

func (r *Redmine) del(in interface{}, out interface{}, uri string, statusExpected int) (int, error) {

	return r.alter("DELETE", in, out, uri, statusExpected)
}

func (r *Redmine) alter(method string, in interface{}, out interface{}, uri string, statusExpected int) (int, error) {

	var er errorsResult

	url := r.endpoint + uri

	s, err := json.Marshal(in)
	if err != nil {
		return 0, err
	}
	req, err := http.NewRequest(method, url, strings.NewReader(string(s)))
	if err != nil {
		return 0, err
	}

	// Set headers
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Redmine-API-Key", r.apiKey)

	// Make request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}

	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	if res.StatusCode != statusExpected {
		err = decoder.Decode(&er)
		if err == nil {
			err = errors.New(strings.Join(er.Errors, "\n"))
		}
	} else {
		if out != nil {
			err = decoder.Decode(out)
		}
	}

	if err != nil {
		return res.StatusCode, err
	}

	return res.StatusCode, nil
}