package api

import (
	"code.dsg.com/smartthings_notif/stnotif"
	"code.dsg.com/smartthings_notif/stnotif/dao"
	"encoding/json"
	A "github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func newEvReq(assert *A.Assertions, url string) *httptest.ResponseRecorder {
	// Create a request to pass to our handler
	req, err := http.NewRequest("GET", url, nil)
	assert.Nil(err)

	// we want JSON
	req.Header["Accept"] = append(req.Header["Accept"], "application/json")

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()

	// We need to use the gorilla mux to process query params of the request
	s.router.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	assert.Equal(http.StatusOK, rr.Code)

	return rr
}

func TestEventsState(t *testing.T) {
	assert := A.New(t)
	assert.NotNil(s)

	rr := newEvReq(assert, "/events/state")

	// Check the response body is what we expect.
	var recs []dao.NotifRec
	err := json.Unmarshal(rr.Body.Bytes(), &recs)
	assert.Nil(err)
	assert.Equal(len(recs), 17, "we have this many distinct devices in fixtures")
}

func TestGetEvents(t *testing.T) {
	assert := A.New(t)
	assert.NotNil(s)

	rr := newEvReq(assert, "/events?since=12/08/2018+00:00:00")
	// Check the response body is what we expect.
	var recs []dao.NotifRec
	err := json.Unmarshal(rr.Body.Bytes(), &recs)
	assert.Nil(err)
	assert.Equal(1256, len(recs), "we have this many recs since in fixtures")

	f.AddFixture(nil) // add one for the next test
	rr = newEvReq(assert, "/events?since=1h")
	err = json.Unmarshal(rr.Body.Bytes(), &recs)
	assert.Nil(err)
	assert.Equal(1, len(recs), "only the one we just added")

	rr = newEvReq(assert, "/events/device/fixture?since=1h")
	err = json.Unmarshal(rr.Body.Bytes(), &recs)
	assert.Nil(err)
	assert.Equal(1, len(recs), "only the one we just added")

	// add event from 2 hours ago
	f.AddFixture(&dao.NotifRec{0, "fixture", time.Now().UTC().Unix() - (60 * 60 * 2), "ev", "eval", "fixture 2hr old event"})
	f.AddFixture(&dao.NotifRec{0, "other", time.Now().UTC().Unix() - (60 * 60 * 2), "ev", "eval", "other 2hr old event"})
	rr = newEvReq(assert, "/events/device/fixture?since=1h")
	err = json.Unmarshal(rr.Body.Bytes(), &recs)
	assert.Nil(err)
	assert.Equal(1, len(recs), "not to include the 2hr old event")
	rr = newEvReq(assert, "/events/device/fixture?since=2.1h")
	err = json.Unmarshal(rr.Body.Bytes(), &recs)
	assert.Nil(err)
	assert.Equal(2, len(recs), "to include the 2hr old fixture event")
	rr = newEvReq(assert, "/events?since=2.1h")
	err = json.Unmarshal(rr.Body.Bytes(), &recs)
	assert.Nil(err)
	assert.Equal(3, len(recs), "to include the 2hr old events for all")

	f, _ = stnotif.NewFixtures() // reload for other tests
}
