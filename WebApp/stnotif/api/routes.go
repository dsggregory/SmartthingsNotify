package api

import (
	"code.dsg.com/smartthings_notif/stnotif/dao"
	"encoding/json"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"path"
	"time"
)

// Store an event notification record.
// POST /events
// Body is JSON array of dao.NotifRec
func (s *server) addEvent(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err == nil {
		n := dao.NotifRec{}
		err = json.Unmarshal(body, &n)
		if err == nil {
			err = s.db.AddEvent(n)
			if err == nil {
				w.WriteHeader(http.StatusCreated)
			}
		}
	}
	if err != nil {
		log.WithError(err).Error("cannot add event")
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// Get events since some time.
// GET /events?since={mm/dd/yy+HH:MM:SS}
func (s *server) getEvents(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	t, err := dao.SinceFormatToTime(vars["since"])
	if err == nil {
		var events []dao.NotifRec
		events, err = s.db.GetEvents(t)
		if err == nil {
			s.respondWithEvents(events, path.Join(s.appDir, "views/events.html"), w, r)
		}
	}
	if err != nil {
		log.WithError(err).WithField("since", vars["since"]).Error("cannot get events")
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// Get current state of all known devices.
// GET /events/state
func (s *server) getEventsState(w http.ResponseWriter, r *http.Request) {
	var events []dao.NotifRec
	events, err := s.db.GetLastByDevice()
	if err == nil {
		s.respondWithEvents(events, path.Join(s.appDir, "views/state.html"), w, r)
	} else {
		log.WithError(err).Error("cannot get events")
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *server) getDeviceEvents(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var t time.Time
	var err error
	since := vars["since"] // optional
	if len(since) == 0 {
		since = "01/01/1970 00:00:00"
	}
	t, err = dao.SinceFormatToTime(since)
	if err == nil {
		var events []dao.NotifRec
		events, err = s.db.GetDeviceEvents(vars["device"], &t)
		if err == nil {
			s.respondWithEvents(events, "views/events.html", w, r)
		}
	}
	if err != nil {
		log.WithError(err).WithField("since", vars["since"]).Error("cannot get events")
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// create the routes we will support
func (s *server) initRoutes() {
	/* Notes on Gorilla MUX:
		o With named routes, other code may lookup the Path
	   		url, err := r.Get(route_name).URL(param_name, param_value, ...)
		o All params are required to match - this includes query (?.*) params that are defined for the route
	*/
	s.router.HandleFunc("/event", s.addEvent).
		Methods("POST").
		Name("addEvent")
	s.router.HandleFunc("/events", s.getEvents).
		Methods("GET").
		Queries("since", "{since}").
		Name("getEvents")
	s.router.HandleFunc("/events/state", s.getEventsState).
		Methods("GET").
		Name("getEventsState")
	s.router.HandleFunc("/", s.getEventsState).
		Methods("GET")
	s.router.HandleFunc("/events/device/{device}", s.getDeviceEvents).
		Methods("GET").
		Queries("since", "{since}").
		Name("getDeviceEvents")

	// Google Sheets faux routes for use by the SmartApp
	s.router.HandleFunc("/gs", s.googleSheetsEndpointGet).Methods("GET")
	s.router.HandleFunc("/gs", s.googleSheetsEndpointPost).Methods("POST")

	// default for static files (CSS, images, et.al.)
	s.router.PathPrefix("/").Handler(http.FileServer(http.Dir("./assets/")))
}
