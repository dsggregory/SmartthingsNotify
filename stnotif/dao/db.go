package dao

import (
	"database/sql"
	"time"

	"code.dsg.com/smartthings_notif/stnotif/conf"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

const (
	goMysqlTimeFormat string = "2006-01-02 15:04:05" // "12/08/2018 7:08:16"
	// SinceDateFormat is the expected time format for GetEvents
	SinceDateFormat string = "01/02/2006 15:04:05"
)

// NotifRec is a notifications record
type NotifRec struct {
	ID          int
	Device      string
	EvTime      int64
	Event       string
	Value       string
	Description string
}

// DbHandle is a handle to the database
type DbHandle struct {
	conn    *sql.DB
	addStmt *sql.Stmt
	getStmt *sql.Stmt
}

// UnixToMysqlTime converts time_t to YYYY-MM-DD hh:mm:ss
func UnixToMysqlTime(ti int64) string {
	return time.Unix(ti, 0).UTC().Format(goMysqlTimeFormat)
}

// MysqlTimeToUnix converts YYYY-MM-DD hh:mm:ss to time_t
func MysqlTimeToUnix(ts string) int64 {
	t, _ := time.Parse(goMysqlTimeFormat, ts)
	return t.Unix()
}

// AddEvent inserts an event into the table
func (d *DbHandle) AddEvent(n NotifRec) error {
	_, err := d.addStmt.Exec(n.Device, n.EvTime, n.Event, n.Value, n.Description)
	return err
}

// GetEvents returns an array of events since some time
func (d *DbHandle) GetEvents(since time.Time) ([]NotifRec, error) {
	var recs []NotifRec
	rows, err := d.getStmt.Query(UnixToMysqlTime(since.Unix()))
	defer rows.Close()
	if err == nil {
		for rows.Next() {
			var n NotifRec
			err = rows.Scan(
				&n.ID,
				&n.Device,
				&n.EvTime,
				&n.Event,
				&n.Value,
				&n.Description)
			if err == nil {
				recs = append(recs, n)
			}
		}
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return recs, nil
}

// NewDbHandler creates an instance of the dao
func NewDbHandler(conf *conf.Conf) (*DbHandle, error) {
	dsn := conf.DbServer.User
	if len(conf.DbServer.Password) > 0 {
		dsn += ":" + conf.DbServer.Password
	}
	if len(conf.DbServer.Host) > 0 {
		//tcp(127.0.0.1:3306)
		dsn += "@tcp(" + conf.DbServer.Host + ":" + conf.DbServer.Port + ")"
	}
	dsn += "/" + conf.DbServer.Database + "?charset=utf8"
	log.WithField("dsn", dsn).Debug("opening database")
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	d := DbHandle{conn: conn}

	d.getStmt, err = conn.Prepare("SELECT id, device_name, time, event, value, description FROM smartthings.notifications WHERE time>=?")
	if err != nil {
		log.WithError(err).Fatal("can't prepare GET statement")
	}

	d.addStmt, err = conn.Prepare("INSERT INTO smartthings.notifications SET id=0, device_name=?, time=?, event=?, value=?, description=?")
	if err != nil {
		log.WithError(err).Fatal("can't prepare ADD statement")
	}

	return &d, nil
}
