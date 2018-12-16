package api

import (
	"code.dsg.com/smartthings_notif/stnotif/dao"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"html/template"
	"net/http"
	"path"
)

func requestAccepts(r *http.Request, mime string) bool {
	accepts := r.Header["Accept"]
	for _, n := range accepts {
		if n == mime {
			return true
		}
	}
	return false
}

// GetDeviceEventsHref returns the route URL to get a devices events
func (s *server) GetDeviceEventsHref(device string) string {
	href, err := s.router.Get("getDeviceEvents").URL("device", device, "since", "")
	if err == nil {
		return href.String()
	} else {
		return ""
	}
}

// Respond with an array of events in JSON or by using the view specified by templatePath for HTML.
func (s *server) respondWithEvents(events []dao.NotifRec, templatePath string, w http.ResponseWriter, r *http.Request) {
	if requestAccepts(r, "application/json") {
		var j []byte
		var err error
		if len(events) > 0 {
			j, err = json.Marshal(events)
		} else {
			j = []byte("{}")
		}
		if err == nil {
			w.Header().Set("Content-Type", "application/json")
			w.Write(j)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		// see http://goinbigdata.com/example-of-using-templates-in-golang/ for a good doc on templates
		fmap := template.FuncMap{
			"getDeviceEventsHref": s.GetDeviceEventsHref,
		}
		t, err := template.New(path.Base(templatePath)).Funcs(fmap).ParseFiles(templatePath)
		if err == nil {
			eventsUrl, _ := s.router.Get("getEvents").URL()
			data := struct {
				Events        []dao.NotifRec
				Server        *server
				GetEventsHref string
			}{
				Events:        events,
				Server:        s,
				GetEventsHref: eventsUrl.String(),
			}
			err = t.Execute(w, data)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.WithError(err).Error()
			}
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			log.WithError(err).Error()
		}
	}
}
