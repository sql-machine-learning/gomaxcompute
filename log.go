package gomaxcompute

import (
	"fmt"
	"io"
	"os"
	"path"

	"github.com/sirupsen/logrus"
)

var log *logrus.Entry

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func init() {
	logDir := getEnv("GOMAXCOMPUTE_log_dir", "")
	logLevel := getEnv("GOMAXCOMPUTE_log_level", "info")

	ll, e := logrus.ParseLevel(logLevel)
	if e != nil {
		ll = logrus.InfoLevel
	}
	var f io.Writer
	if logDir != "" {
		e = os.MkdirAll(logDir, 0744)
		if e != nil {
			log.Panicf("create log directory failed: %v", e)
		}

		f, e = os.Create(path.Join(logDir, fmt.Sprintf("gomaxcompute-%d.log", os.Getpid())))
		if e != nil {
			log.Panicf("open log file failed: %v", e)
		}
	} else {
		f = os.Stdout
	}

	lg := logrus.New()
	lg.SetOutput(f)
	lg.SetLevel(ll)
	log = lg.WithFields(logrus.Fields{"driver": "odps"})
}
