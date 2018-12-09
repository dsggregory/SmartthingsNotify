package api

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"

	log "github.com/sirupsen/logrus"
)

func (s *server) fooHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!\n", r.URL.Path[1:])
	log.Infof("%+v\n", r)
}

func (s *server) addEvent(w http.ResponseWriter, r *http.Request) {
	log.Info("in addEvent")
	body, err := ioutil.ReadAll(r.Body)
	if err == nil {
		log.WithField("body", string(body)).Info()
	} else {
		log.WithError(err).Error("cannot read body")
		w.WriteHeader(500)
	}
}

func (s *server) getEvents(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.WithField("since", vars["since"]).Infof("in getEvents")
}

// create the routes we will support
func (s *server) initRoutes() {
	/* With named routes, other code may lookup the Path
	   url, err := r.Get(route_name).URL(param_name, param_value, ...)
	*/
	s.router.HandleFunc("/foo", s.fooHandler).
		Methods("GET").
		Name("foo")

	s.router.HandleFunc("/event", s.addEvent).
		Methods("POST").
		Name("addEvent")
	s.router.HandleFunc("/events", s.getEvents).
		Methods("GET").
		Queries("since", "{since}").
		Name("getEvents")

	s.dumpRoutes()
}
