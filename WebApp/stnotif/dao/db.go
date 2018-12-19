package dao

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"code.dsg.com/smartthings_notif/stnotif/conf"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

const (
	goMysqlTimeFormat string = "2006-01-02 15:04:05" // "12/08/2018 7:08:16"
	// SinceDateFormat is the expected time format for GetEvents
	SinceDateFormat string = "1/2/2006 15:04:05"
)

// NotifRec is a notifications record
type NotifRec struct {
	ID          int
	Device      string
	EvTime      int64 // time_t of UTC
	Event       string
	Value       string
	Description string
}

// EventTime returns the string version of the NotifyRec.EvTime in localtime
func (n *NotifRec) EventTime() string {
	return time.Unix(n.EvTime, 0).Format("Mon, 2 Jan 15:04:05")
}

// DbHandle is a handle to the database
type DbHandle struct {
	conf          *conf.Conf
	dbname        string
	conn          *sql.DB
	addStmt       *sql.Stmt
	getStmt       *sql.Stmt
	getDeviceStmt *sql.Stmt
}

// UnixToMysqlTime converts time_t to YYYY-MM-DD hh:mm:ss in localtime
func UnixToMysqlTime(ti int64) string {
	return time.Unix(ti, 0).Format(goMysqlTimeFormat)
}

// MysqlTimeToUnix converts YYYY-MM-DD hh:mm:ss in localtime to a time_t
func MysqlTimeToUnix(ts string) int64 {
	loc, _ := time.LoadLocation("Local")
	t, _ := time.ParseInLocation(goMysqlTimeFormat, ts, loc)
	return t.Unix()
}

// SinceFormatToTime converts the "since" format to Time
// The since may be an exact time in "mm/dd/yyyy HH:MM:SS" localtime or a duration "6h" (h|m|s)
func SinceFormatToTime(since string) (time.Time, error) {
	if len(since) > 0 {
		d, err := time.ParseDuration(since)
		if err == nil {
			return time.Now().Add(-d), nil
		}
		loc, _ := time.LoadLocation("Local")
		return time.ParseInLocation(SinceDateFormat, since, loc)
	}
	return time.Parse(SinceDateFormat, "1/1/1970 00:00:00")
}

// AddEvent inserts an event into the table
func (d *DbHandle) AddEvent(n NotifRec) error {
	// convert hub time to UTC
	n.EvTime = time.Unix(n.EvTime, 0).Unix()
	log.WithField("rec", n).Debug("insert record")
	_, err := d.addStmt.Exec(n.Device, n.EvTime, n.Event, n.Value, n.Description)
	return err
}

// parses q Query result set into array of notification records
func (d *DbHandle) notificationsFromQuery(rows *sql.Rows) ([]NotifRec, error) {
	var recs []NotifRec
	err := rows.Err()
	if err != nil {
		return nil, err
	}
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
			// convert UTC time to hub time
			n.EvTime = time.Unix(n.EvTime, 0).Unix()
			recs = append(recs, n)
		}
	}
	return recs, nil
}

// GetEvents returns an array of events since some time
func (d *DbHandle) GetEvents(since time.Time) ([]NotifRec, error) {
	tsince := since.Unix()
	log.WithField("since_t", tsince).WithField("since_tm", since).Debug()
	rows, err := d.getStmt.Query(tsince)
	defer rows.Close()
	if err == nil {
		return d.notificationsFromQuery(rows)
	}
	return nil, err
}

// GetDeviceEvents returns an array of events for the device since some time or epoch if nil
func (d *DbHandle) GetDeviceEvents(device string, since *time.Time) ([]NotifRec, error) {
	var tsince int64
	if since == nil {
		tsince = 0
	} else {
		tsince = since.Unix()
	}
	log.WithFields(log.Fields{
		"device":   device,
		"since_t":  tsince,
		"since_tm": since,
	}).Debug()
	rows, err := d.getDeviceStmt.Query(tsince, device)
	defer rows.Close()
	if err == nil {
		return d.notificationsFromQuery(rows)
	}
	return nil, err
}

// GetLastByDevice returns the current state of all known devices
func (d *DbHandle) GetLastByDevice() ([]NotifRec, error) {
	rows, err := d.conn.Query("select * from notifications where id in (select MAX(id) from notifications group by device_name) order by device_name")
	defer rows.Close()
	if err == nil {
		return d.notificationsFromQuery(rows)
	}
	return nil, err
}

// CreateDatabase creates the database if it doesn't exist
func (d *DbHandle) CreateDatabase(dbname string) error {
	q := "CREATE DATABASE IF NOT EXISTS " + dbname
	_, err := d.conn.Exec(q)
	if err != nil {
		return err
	}

	q = "CREATE TABLE IF NOT EXISTS " + dbname + ".notifications"
	q += ` (
id int primary key auto_increment,
device_name varchar(64),
time bigint unsigned,
event varchar(64),
value varchar(64),
description varchar(128),
KEY (time, device_name)
) ENGINE=INNODB`
	_, err = d.conn.Exec(q)
	if err != nil {
		return err
	}

	return nil
}

// GetCount gets count of rows in notifications
func (d *DbHandle) GetCount() int {
	count := 0
	var col string
	row := d.conn.QueryRow("SELECT count(*) FROM notifications")
	err := row.Scan(&col)
	if err == nil {
		count, _ = strconv.Atoi(col)
	}

	return count
}

// Return the database name of the Data Source Name
func dbnameOfDSN(dsn string) (string, string) {
	var dbname string
	i := strings.LastIndex(dsn, "/")
	if i >= 0 {
		dbname = dsn[i+1:] // save the database name
		j := strings.Index(dbname, "?")
		if j >= 0 {
			dbname = dbname[:j]
		}
		dsn = dsn[:i+1] // stomp on the database name in conf. Requires trailing '/'.
	}

	return dbname, dsn
}

// NewDbHandlerTest creates a test DB after dropping one that exists
func NewDbHandlerTest(conf *conf.Conf) (*DbHandle, error) {
	dbname, dsn := dbnameOfDSN(conf.DbDSN())

	// open without a default database name
	d := DbHandle{conf: conf}
	conn, err := sql.Open(conf.DbDriver(), dsn)
	if err != nil {
		return nil, err
	}
	d.conn = conn

	// drop the existing test database
	_, err = d.conn.Exec("DROP DATABASE IF EXISTS " + dbname)
	if err != nil {
		panic(err)
	}

	// create the test database
	err = d.CreateDatabase(dbname)
	if err != nil {
		return nil, err
	}

	// close and reopen with the database name
	d.conn.Close()

	return NewDbHandler(conf)
}

// NewDbHandler creates an instance of the dao
func NewDbHandler(conf *conf.Conf) (*DbHandle, error) {
	d := DbHandle{conf: conf}
	d.dbname, _ = dbnameOfDSN(conf.DbDSN())

	conn, err := sql.Open(conf.DbDriver(), conf.DbDSN())
	if err != nil {
		return nil, err
	}
	d.conn = conn

	d.getStmt, err = conn.Prepare(fmt.Sprintf("SELECT id, device_name, time, event, value, description FROM %s.notifications WHERE time>=? ORDER BY time DESC", d.dbname))
	if err != nil {
		log.WithError(err).Fatal("can't prepare getEvents statement")
	}

	d.addStmt, err = conn.Prepare(fmt.Sprintf("INSERT INTO %s.notifications SET id=0, device_name=?, time=?, event=?, value=?, description=?", d.dbname))
	if err != nil {
		log.WithError(err).Fatal("can't prepare addEvent statement")
	}

	d.getDeviceStmt, err = conn.Prepare(fmt.Sprintf("SELECT id, device_name, time, event, value, description FROM %s.notifications WHERE time>=? AND device_name=? ORDER BY time DESC", d.dbname))
	if err != nil {
		log.WithError(err).Fatal("can't prepare getDeviceEvents statement")
	}

	return &d, nil
}
