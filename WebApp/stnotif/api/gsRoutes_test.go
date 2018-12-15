package api

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGsGet(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(s)

	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/gs", nil)
	assert.Nil(err)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()

	// a testing example for when the route doesn't require the gorilla mux
	handler := http.HandlerFunc(s.googleSheetsEndpointGet)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	assert.Equal(http.StatusOK, rr.Code)

	// Check the response body is what we expect.
	assert.Equal(rr.Body.String(), GsVersion)

}

func TestGsPost(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(s)

	originalCount := s.db.GetCount()

	// Create a request to pass to our handler.
	events := []gsPostEvent{
		gsPostEvent{"1/1/2018 00:00:00", "Dev1", "activate", "on", "on"},
	}
	pd := postRequestData{
		PostBackUrl:        "", //"https://postback.domain.com",
		ArchiveOptions:     archiveOptions{LogIsFull: false, Type: "type", Interval: 0},
		LogDesc:            true,
		LogReporting:       false,
		DeleteExtraColumns: true,
		Events:             events,
	}
	j, err := json.Marshal(&pd)
	r := bytes.NewReader(j)
	req, err := http.NewRequest("POST", "/gs", r)
	assert.Nil(err)

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()

	// a testing example for when the route doesn't require the gorilla mux
	handler := http.HandlerFunc(s.googleSheetsEndpointPost)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	assert.Equal(http.StatusFound, rr.Code)

	// Check the response body is what we expect.
	assert.Equal(rr.Body.String(), "")

	assert.Equal(originalCount+1, s.db.GetCount())

	setupTest() // recreate database from fixtures since we modified
}
