package api

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"net/http"

	"code.dsg.com/smartthings_notif/stnotif/dao"
	"github.com/gorilla/mux"

	log "github.com/sirupsen/logrus"
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
			var j []byte
			if len(events) > 0 {
				j, err = json.Marshal(events)
			} else {
				j = []byte("{}")
			}
			if err == nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write(j)
			}
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
		status := http.StatusInternalServerError
		t, err := template.ParseFiles("views/state.html")
		if err == nil {
			data := struct {
				Items []dao.NotifRec
			}{
				Items: events,
			}
			err = t.Execute(w, data)
			if err == nil {
				status = http.StatusOK
			} else {
				log.WithError(err).Error()
			}
		} else {
			log.WithError(err).Error()
		}
		w.WriteHeader(status)
	} else {
		log.WithError(err).Error("cannot get events")
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// create the routes we will support
func (s *server) initRoutes() {
	/* With named routes, other code may lookup the Path
	   url, err := r.Get(route_name).URL(param_name, param_value, ...)
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
	s.router.PathPrefix("/").Handler(http.FileServer(http.Dir("./assets/")))
}
