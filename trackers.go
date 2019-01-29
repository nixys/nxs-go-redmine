package redmine

import "net/url"

/* Get */

// TrackerObject struct used for trackers get operations
type TrackerObject struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	DefaultStatus IDName `json:"default_status"` // Since 3.0
}

/* Internal types */

type trackerAllResult struct {
	Trackers []TrackerObject `json:"trackers"`
}

// TrackerAllGet gets info for all trackers
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Trackers#GET
func (r *Context) TrackerAllGet() ([]TrackerObject, int, error) {

	var t trackerAllResult

	ur := url.URL{
		Path: "/trackers.json",
	}

	status, err := r.get(&t, ur, 200)

	return t.Trackers, status, err
}
