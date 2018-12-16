package dao

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMySqlTimeConvert(t *testing.T) {
	type testTime struct {
		ts string
		ti int64
	}
	tests := []testTime{
		testTime{"1970-01-01 00:00:00", int64(0)},
		testTime{"2018-12-09 21:54:20", int64(1544392460)},
	}

	assert := assert.New(t)

	for i := range tests {
		ts := UnixToMysqlTime(tests[i].ti)
		assert.Equal(tests[i].ts, ts)

		ti := MysqlTimeToUnix(tests[i].ts)
		assert.Equal(int64(tests[i].ti), ti)
	}
}

func TestSinceTimeConvert(t *testing.T) {
	type testTime struct {
		ts string
		ti int64
	}
	tests := []testTime{
		testTime{"01/01/1970 00:00:00", int64(0)},
		testTime{"1/1/1970 0:00:00", int64(0)},
		testTime{"12/09/2018 21:54:20", int64(1544392460)},
	}

	assert := assert.New(t)

	for i := range tests {
		t, err := SinceFormatToTime(tests[i].ts)
		assert.Nil(err)
		assert.Equal(tests[i].ti, t.Unix())
	}

	now := time.Now().Unix()
	tm, err := SinceFormatToTime("1h")
	assert.Nil(err)
	assert.True((tm.Unix()+(60*60)+5) >= now, fmt.Sprintf("%d - %d  = %d > (60*60)", now, tm.Unix(), now-tm.Unix()))
}
