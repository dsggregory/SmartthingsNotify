package dao

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTimeConvert(t *testing.T) {
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
