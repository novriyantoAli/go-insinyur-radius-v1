package utils

import (
	"time"

	"github.com/sirupsen/logrus"
)

func FreeradiusStringToDate(date string) (res time.Time, err error) {
	wita, err := time.LoadLocation("Asia/Makassar")

	if err != nil {
		logrus.Error(err)
		return
	}

	layoutFormat := "02 Jan 2006 15:04:05"

	res, err = time.ParseInLocation(layoutFormat, date, wita)
	if err != nil {
		logrus.Error(err)
	}

	return

}

func FreeradiusDateToString(t time.Time) (res string, err error) {
	wita, err := time.LoadLocation("Asia/Makassar")

	if err != nil {
		logrus.Error(err)
		return
	}

	layoutFormat := "02 Jan 2006 15:04:05"

	tm, err := time.ParseInLocation(layoutFormat, t.Format(layoutFormat), wita)
	if err != nil {
		logrus.Error(err)
		return
	}

	res = tm.Format(layoutFormat)

	return
}
