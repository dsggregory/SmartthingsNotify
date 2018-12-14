package api

import (
	"code.dsg.com/smartthings_notif/stnotif/dao"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEventsState(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(s)

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/events/state", nil)
	assert.Nil(err)

	// we want JSON
	req.Header["Accept"] = append(req.Header["Accept"], "application/json")

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()

	// a testing example for when the route doesn't require the gorilla mux
	handler := http.HandlerFunc(s.getEventsState)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	assert.Equal(http.StatusOK, rr.Code)

	// Check the response body is what we expect.
	var recs []dao.NotifRec
	err = json.Unmarshal(rr.Body.Bytes(), &recs)
	assert.Nil(err)
	assert.Equal(len(recs), 17, "we have this many distinct devices in fixtures")
}

func TestGetEvents(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(s)

	// Create a request to pass to our handler
	req, err := http.NewRequest("GET", "/events?since=12/08/2018+00:00:00", nil)
	assert.Nil(err)

	// we want JSON
	req.Header["Accept"] = append(req.Header["Accept"], "application/json")

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()

	// We need to use the gorilla mux to process query params of the request
	s.router.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	assert.Equal(http.StatusOK, rr.Code)

	// Check the response body is what we expect.
	var recs []dao.NotifRec
	err = json.Unmarshal(rr.Body.Bytes(), &recs)
	assert.Nil(err)
	assert.Equal(len(recs), 1256, "we have this many recs since in fixtures")
}
