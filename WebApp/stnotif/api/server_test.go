package api

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAllowedHosts(t *testing.T) {
	assert := assert.New(t)
	assert.NotNil(s)

	s.config.AllowedHosts = s.config.AllowedHosts[0:0]
	s.config.AllowedHosts = append(s.config.AllowedHosts, "foo.com")

	// /noop is unknown route
	req, err := http.NewRequest("GET", "/noop", nil)
	assert.Nil(err)
	rr := httptest.NewRecorder()
	s.wrapRequest(s.router).ServeHTTP(rr, req)
	assert.Equal(http.StatusForbidden, rr.Code)

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
