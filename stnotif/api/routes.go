package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"code.dsg.com/smartthings_notif/stnotif/dao"
	"github.com/gorilla/mux"

	log "github.com/sirupsen/logrus"
)

func (s *server) addEvent(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err == nil {
		n := dao.NotifRec{}
		err = json.Unmarshal(body, &n)
		if err == nil {
			err = s.db.AddEvent(n)
			if err == nil {
				w.WriteHeader(201)
			}
		}
	}
	if err != nil {
		log.WithError(err).Error("cannot add event")
		w.WriteHeader(500)
	}
}

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
				w.WriteHeader(200)
				w.Write(j)
			}
		}
	}
	if err != nil {
		log.WithError(err).WithField("since", vars["since"]).Error("cannot get events")
		w.WriteHeader(500)
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
		Queries("since", "{since}"). // time.UnixFormat()
		Name("getEvents")
}
