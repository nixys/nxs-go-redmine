package redmine

/* Get */

// TrackerObject struct used for trackers get operations
type TrackerObject struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

/* Internal types */

type trackerMultiResult struct {
	Trackers []TrackerObject `json:"trackers"`
}

// TrackerMultiGet gets multiple trackers
//
// see: http://www.redmine.org/projects/redmine/wiki/Rest_Trackers#Rest-Trackers
func (r *Context) TrackerMultiGet() ([]TrackerObject, int, error) {

	var t trackerMultiResult

	uri := "/trackers.json"

	status, err := r.get(&t, uri, 200)

	return t.Trackers, status, err
}
