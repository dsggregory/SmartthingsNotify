package api

import (
	"code.dsg.com/smartthings_notif/stnotif/dao"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func assertPanic(t *testing.T, msg string, f func()) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf(msg)
		}
	}()
	f()
}

func TestAllowedHosts(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(s)

	s.config.AllowedHosts = s.config.AllowedHosts[0:0]
	s.config.AllowedHosts = append(s.config.AllowedHosts, "foo.com")

	// /noop is unknown route
	req, err := http.NewRequest("GET", "/noop", nil)
	assert.Nil(err)
	rr := httptest.NewRecorder()
	rr.Code = 0
	assertPanic(t, "did not silently refuse", func() {
		s.wrapRequest(s.router).ServeHTTP(rr, req)
	})
	assert.Equal(0, rr.Code, "response should not have written or set code")

	s.config.AllowedHosts = append(s.config.AllowedHosts, "\\[::1\\]")  // ipv6
	hosts := []string{"foo.com", "foo.com:port", "[::1]", "[::1]:port"} // remote addrs to test
	for i := range hosts {
		req, err = http.NewRequest("GET", "/noop", nil)
		req.RemoteAddr = hosts[i]
		assert.Nil(err)
		rr = httptest.NewRecorder()
		s.wrapRequest(s.router).ServeHTTP(rr, req)
		assert.Equal(http.StatusNotFound, rr.Code, "should get 404 for "+hosts[i])
	}
}

func TestDbHandle_AddEvent(t *testing.T) {
	assert := assert.New(t)

	now := time.Now().UTC()

	f.Config.HubTzLocation = time.FixedZone("negOne", -100)
	f.AddFixture(nil)
	tm, err := dao.SinceFormatToTime("5m")
	ev, err := f.DbHandle.GetEvents(tm)
	assert.Nil(err)
	assert.Equal(1, len(ev))
	assert.WithinDuration(now, time.Unix(ev[0].EvTime, 0), 2*time.Second)

	// change hub tz and still get the event we added
	f.Config.HubTzLocation = time.UTC
	ev, err = f.DbHandle.GetEvents(tm)
	assert.Nil(err)
	assert.Equal(1, len(ev))
	assert.WithinDuration(now, time.Unix(ev[0].EvTime, 0), 2*time.Second)
}
