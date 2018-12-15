package api

import (
	"code.dsg.com/smartthings_notif/stnotif/dao"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

// GET https://script.google.com/macros/s/{key}/GG/exec
// A GET returns 200 text/plain of "Version 01.03.00"
func (s *server) googleSheetsEndpointGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("Version 01.03.00"))
}

type gsPostEvent struct {
	Time   string
	Device string
	Name   string
	Value  string
	Desc   string
}
type postRequestData struct {
	PostBackUrl        string
	ArchiveOptions     string
	LogDesc            string
	LogReporting       string
	DeleteExtraColumns string
	Events             []gsPostEvent
}

// POST https://script.google.com/macros/s/{key}/GG/exec
// A POST sends data to the spreadsheet.
// application/json array of {time, device, name, value, descr}
/* Request is
[
		postBackUrl: "${state.endpoint}update-logging-status",
		archiveOptions: getArchiveOptions(),
		logDesc: (settings?.logDesc != false),
		logReporting: (settings?.logReporting == true),
		deleteExtraColumns: (settings?.deleteExtraColumns == true),
		events: events
	]
Responds with 302 and no body
*/
func (s *server) googleSheetsEndpointPost(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err == nil {
		rd := postRequestData{}
		err = json.Unmarshal(body, &rd)
		if err == nil {
			log.WithFields(log.Fields{
				"n_events":           len(rd.Events),
				"postBackUrl":        rd.PostBackUrl,
				"archiveOptions":     rd.ArchiveOptions,
				"logDesc":            rd.LogDesc,
				"logReporting":       rd.LogReporting,
				"deleteExtraColumns": rd.DeleteExtraColumns,
			}).Debug("got smartthings events")
			for i := range rd.Events {
				ir := rd.Events[i]
				log.WithField("event", ir).Debugf("event[%d]", i)
				t, _ := dao.SinceFormatToTime(ir.Time)
				n := dao.NotifRec{
					0,
					ir.Device,
					t.Unix(),
					ir.Name,
					ir.Value,
					ir.Desc}
				aerr := s.db.AddEvent(n)
				if aerr != nil {
					log.WithError(aerr).WithField("input-event", ir).Error("AddEvent failed")
				}
			}
			w.WriteHeader(http.StatusFound)
		}
	}
	if err != nil {
		log.WithError(err).Error("cannot add event")
		w.WriteHeader(http.StatusInternalServerError)
	}
}
