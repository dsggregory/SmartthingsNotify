package api

import (
	"code.dsg.com/smartthings_notif/stnotif"
	"code.dsg.com/smartthings_notif/stnotif/dao"
	"encoding/json"
	A "github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
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
	assert.Equal(len(recs), 1256, "we have this many recs since in fixtures")

	f.AddFixture() // add one for the next test
	rr = newEvReq(assert, "/events?since=1h")
	err = json.Unmarshal(rr.Body.Bytes(), &recs)
	assert.Nil(err)
	assert.Equal(len(recs), 1, "only the one we just added")
	f, _ = stnotif.NewFixtures() // reload for other tests
}
