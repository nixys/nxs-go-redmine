package redmine

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/mitchellh/mapstructure"
)

const (
	limitDefault = 100
)

// Context struct used for store settings to communicate with Redmine API
type Context struct {
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

// SetAPIKey is used to set Redmine API key
func (r *Context) SetAPIKey(apiKey string) {
	r.apiKey = apiKey
}

// SetEndpoint is used to set Redmine endpoint
func (r *Context) SetEndpoint(endpoint string) {
	r.endpoint = endpoint
}

// SetLimit is used to set elements limit on page
func (r *Context) SetLimit(limit int) {
	r.limit = limit
}

func (r *Context) get(out interface{}, uri string, statusExpected int) (int, error) {

	var er errorsResult

	u := r.endpoint + uri

	// Create request
	req, err := http.NewRequest("GET", u, nil)
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

	dJ := json.NewDecoder(res.Body)

	if res.StatusCode != statusExpected {
		if err := dJ.Decode(&er); err != nil {
			return res.StatusCode, err
		}
		err = errors.New(strings.Join(er.Errors, "\n"))
	} else {
		if out != nil {

			rawConf := make(map[string]interface{})

			if err := dJ.Decode(&rawConf); err != nil {
				return res.StatusCode, fmt.Errorf("json decode error: %v", err)
			}

			dM, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
				WeaklyTypedInput: true,
				Result:           out,
				TagName:          "json",
			})
			if err != nil {
				return res.StatusCode, fmt.Errorf("mapstructure create decoder error: %v", err)
			}

			if err := dM.Decode(rawConf); err != nil {
				return res.StatusCode, fmt.Errorf("mapstructure decode error: %v", err)
			}
		}
	}

	return res.StatusCode, err
}

func (r *Context) post(in interface{}, out interface{}, uri string, statusExpected int) (int, error) {

	return r.alter("POST", in, out, uri, statusExpected)
}

func (r *Context) put(in interface{}, out interface{}, uri string, statusExpected int) (int, error) {

	return r.alter("PUT", in, out, uri, statusExpected)
}

func (r *Context) del(in interface{}, out interface{}, uri string, statusExpected int) (int, error) {

	return r.alter("DELETE", in, out, uri, statusExpected)
}

func (r *Context) alter(method string, in interface{}, out interface{}, uri string, statusExpected int) (int, error) {

	var er errorsResult

	u := r.endpoint + uri

	s, err := json.Marshal(in)
	if err != nil {
		return 0, err
	}
	req, err := http.NewRequest(method, u, strings.NewReader(string(s)))
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

	dJ := json.NewDecoder(res.Body)

	if res.StatusCode != statusExpected {
		if err := dJ.Decode(&er); err != nil {
			return res.StatusCode, err
		}
		err = errors.New(strings.Join(er.Errors, "\n"))
	} else {
		if out != nil {

			rawConf := make(map[string]interface{})

			if err := dJ.Decode(&rawConf); err != nil {
				return res.StatusCode, fmt.Errorf("json decode error: %v", err)
			}

			dM, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
				WeaklyTypedInput: true,
				Result:           out,
				TagName:          "json",
			})
			if err != nil {
				return res.StatusCode, fmt.Errorf("mapstructure create decoder error: %v", err)
			}

			if err := dM.Decode(rawConf); err != nil {
				return res.StatusCode, fmt.Errorf("mapstructure decode error: %v", err)
			}
		}
	}

	return res.StatusCode, err
}

func (r *Context) uploadFile(filPath string, out interface{}, uri string, statusExpected int) (int, error) {

	var er errorsResult

	u := r.endpoint + uri

	c, err := ioutil.ReadFile(filPath)
	if err != nil {
		return 0, err
	}

	// Create request
	req, err := http.NewRequest("POST", u, bytes.NewBuffer(c))
	if err != nil {
		return 0, err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Add("X-Redmine-API-Key", r.apiKey)

	// Make request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	dJ := json.NewDecoder(res.Body)

	if res.StatusCode != statusExpected {
		if err := dJ.Decode(&er); err != nil {
			return res.StatusCode, err
		}
		err = errors.New(strings.Join(er.Errors, "\n"))
	} else {
		if out != nil {
			if err := dJ.Decode(&out); err != nil {
				return res.StatusCode, fmt.Errorf("json decode error: %v", err)
			}
		}
	}

	return res.StatusCode, err
}
