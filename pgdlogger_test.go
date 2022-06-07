package pgdlogger_test

import (
	"io/ioutil"
	"testing"

	"github.com/samandajimmy/pgdlogger"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestSetRequestId(t *testing.T) {
	pgdlogger.SetRequestId("1122334455")
	actual := pgdlogger.GetRequestId()

	assert.Equal(t, "1122334455", actual)
}

func TestMake(t *testing.T) {
	pgdlogger.Init("debug")
	pgdlogger.SetRequestId("1122334455")
	logrus.SetOutput(ioutil.Discard) // Send all logs to nowhere by default
	entry := pgdlogger.Make()
	expected := logrus.Fields{
		"_requestId": "1122334455",
	}

	assert.Equal(t, expected, entry.Data)

	pgdlogger.SetRequestId("")
	entry = pgdlogger.Make()
	expected = logrus.Fields{}

	assert.Equal(t, expected, entry.Data)

	mapObj := map[string]interface{}{
		"requestPayload": "ini dia",
		"password":       "1234",
	}
	entry = pgdlogger.Make(mapObj)
	expected = logrus.Fields{
		"data": mapObj,
	}

	assert.Equal(t, expected, entry.Data)

	mapObj = map[string]interface{}{
		"requestPayload": map[string]interface{}{
			"key": "value",
		},
	}
	entry = pgdlogger.MakeWithoutReportCaller(mapObj)
	expected = logrus.Fields{
		"data": mapObj,
	}

	assert.Equal(t, expected, entry.Data)
}

func TestInit(t *testing.T) {
	pgdlogger.Init("something")
	pgdlogger.Dump()

	assert.Equal(t, logrus.DebugLevel, logrus.GetLevel())
}
