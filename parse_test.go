package main

import (
	"io/ioutil"
	"path"
	"testing"
	"time"

	"github.com/prometheus/common/promlog"
)

func init() { // so we don't panic
	allowedLevel := promlog.AllowedLevel{}
	allowedLevel.Set("debug")
	allowedFormat := promlog.AllowedFormat{}
	allowedFormat.Set("logfmt")
	promlogConfig := promlog.Config{
		Level:  &allowedLevel,
		Format: &allowedFormat,
	}
	logger = promlog.New(&promlogConfig)
}

func TestParsing(t *testing.T) {
	cases := []struct {
		filename string
		date     time.Time
	}{
		{filename: "google.cn", date: time.Date(2019, 3, 17, 12, 48, 36, 0, time.UTC)},
		{filename: "google.jp", date: time.Date(2021, 5, 31, 0, 0, 0, 0, time.UTC)},
		{filename: "google.com", date: time.Date(2020, 9, 14, 4, 0, 0, 0, time.UTC)},
		{filename: "ietf.org", date: time.Date(2020, 3, 12, 5, 0, 0, 0, time.UTC)},
		{filename: "unisportstore.fi", date: time.Date(2019, 3, 20, 17, 13, 49, 0, time.UTC)},
		// This is a .com WHOIS response that has a different format.
		{filename: "com", date: time.Date(2018, 7, 3, 19, 6, 9, 0, time.UTC)},
		{filename: "io", date: time.Date(2018, 12, 21, 17, 35, 22, 0, time.UTC)},
	}

	for i := range cases {
		w := path.Join("testdata", cases[i].filename)
		ans, err := ioutil.ReadFile(w)
		if err != nil {
			t.Errorf("problem on %s: %v", cases[i].filename, err)
		}

		parsedDate, err := parse(cases[i].filename, ans)
		if parsedDate == 0 || err != nil {
			t.Errorf("%s got %.2f (error=%v) ", cases[i].filename, parsedDate, err)
		}

		answerDate := float64(cases[i].date.Unix())
		if answerDate != parsedDate {
			t.Errorf("cases[%d]: parsedDate=%.0f answerDate=%.0f", i, parsedDate, answerDate)
		}
	}
}
