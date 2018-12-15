package api

import (
	"bytes"
	"code.dsg.com/smartthings_notif/stnotif/dao"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	GsVersion = "Version 01.03.00"
)

// GET https://script.google.com/macros/s/{key}/GG/exec
// A GET returns 200 text/plain of "Version 01.03.00"
func (s *server) googleSheetsEndpointGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(GsVersion))
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

type postbackData struct {
	Success           bool
	EventsArchived    bool   // always false
	LogIsFull         bool   // always false
	GsVersion         string // GsVersion
	Finished          int64  // time in millisec
	EventsLogged      int    // count
	TotalEventsLogged int    // count
	FreeSpace         string // "unlimited"
	Error             string `json:"error",omitempty` // always empty
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
	pb := postbackData{
		Success:        false,
		EventsArchived: false,
		LogIsFull:      false,
		GsVersion:      GsVersion,
		FreeSpace:      "unlimited",
		Error:          "",
	}
	rd := postRequestData{}
	body, err := ioutil.ReadAll(r.Body)
	if err == nil {
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
				pb.EventsLogged++
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
				} else {
					pb.TotalEventsLogged++
				}
			}
			w.WriteHeader(http.StatusFound)
			pb.Success = true
		} else {
			log.WithField("body", body).Debug("can't parse body")
		}
	}
	if err != nil {
		pb.Error = "service error"
		log.WithError(err).Error("cannot add event")
		w.WriteHeader(http.StatusInternalServerError)
	}

	// We have to call back into the SmartApp to tell it how we did or it may continue to send
	// us the same data or it may just log it.
	pb.Finished = time.Now().UnixNano() / int64(time.Millisecond)
	if len(rd.PostBackUrl) > 0 {
		pbdata, err := json.Marshal(&pb)
		if err == nil {
			log.WithField("pbdata", pb).WithField("postback-url", rd.PostBackUrl).
				Info("calling gs postback")
			rdr := bytes.NewReader(pbdata)
			resp, err := http.Post(rd.PostBackUrl, "application/json", rdr)
			if err == nil {
				body, err := ioutil.ReadAll(resp.Body)
				if err == nil {
					log.WithField("response", string(body)).WithField("post", pbdata).Debug("gs postback success")
				} else {
					log.WithError(err).Error("gs postback read body")
				}
			} else {
				log.WithError(err).Error("gs postback post failed")
			}
		} else {
			log.WithError(err).Error("gs postback request json unmarshal")
		}
	}
}
