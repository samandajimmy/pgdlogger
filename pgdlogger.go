package pgdlogger

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"
)

const (
	timestampFormat   = "2006-01-02 15:04:05.000"
	starString        = "**********"
	fieldKeyRequestId = "_requestId"
	fieldKeyData      = "data"
)

var (
	strExclude = []string{"password", "base64", "npwp", "phone", "nik", "ktp", "gaji", "othr",
		"slik"}

	requestId = ""
)

func Init(loglvl string) {
	logrus.SetReportCaller(true)
	formatter := &logrus.JSONFormatter{
		TimestampFormat: timestampFormat,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			tmp := strings.Split(f.File, "/")
			filename := tmp[len(tmp)-1]
			return "", fmt.Sprintf("%s:%d", filename, f.Line)
		},
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyMsg:  "_message",
			logrus.FieldKeyTime: "_time",
		},
	}

	logrus.SetFormatter(formatter)
	// default log level
	logLevel := logrus.DebugLevel
	logrus.SetLevel(logLevel)
	logLevel, err := logrus.ParseLevel(loglvl)

	if err != nil {
		Make().Debug(err)
		return
	}

	logrus.SetLevel(logLevel)
}

func SetRequestId(reqId string) {
	requestId = reqId
}

func GetRequestId() string {
	return requestId
}

func Make(data ...map[string]interface{}) *logrus.Entry {
	dataMap := map[string]interface{}{}
	logField := logrus.Fields{}
	logrus.SetReportCaller(true)

	if requestId != "" {
		logField[fieldKeyRequestId] = requestId
	}

	if len(data) > 0 {
		dataMap = data[0]
		payloadExcluder(&dataMap)
		logField[fieldKeyData] = dataMap
	}

	return logrus.WithFields(logField)
}

func MakeWithoutReportCaller(data ...map[string]interface{}) *logrus.Entry {
	log := Make(data...)
	logrus.SetReportCaller(false)

	return log
}

func Dump(strct ...interface{}) {
	fmt.Println("DEBUGGING ONLY")
	spew.Dump(strct)
	fmt.Println("DEBUGGING ONLY")
}

func reExcludePayload(pl interface{}) (map[string]interface{}, bool) {
	vMap, ok := pl.(map[string]interface{})

	if !ok {
		return map[string]interface{}{}, ok
	}

	payloadExcluder(&vMap)

	return vMap, true
}

func payloadExcluder(pl *map[string]interface{}) {
	var ok bool
	var vMap map[string]interface{}
	plMap := *pl

	for k, v := range plMap {
		vMap, ok = reExcludePayload(v)

		if ok {
			plMap[k] = vMap
			continue
		}

		if contains(strExclude, k) {
			v = starString
		}

		plMap[k] = v
	}

	*pl = plMap
}

func contains(strIncluder []string, str string) bool {
	for _, include := range strIncluder {
		if strings.Contains(strings.ToLower(str), include) {
			return true
		}
	}

	return false
}
