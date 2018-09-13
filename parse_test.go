package main

import (
	"io/ioutil"
	"math"
	"path"
	"testing"
	"time"

	"github.com/prometheus/common/promlog"
)

func init() { // so we don't panic
	level := promlog.AllowedLevel{}
	level.Set("info")
	logger = promlog.New(level)
}

func TestParsing(t *testing.T) {
	cases := []struct {
		filename string
		date     time.Time
	}{
		{filename: "google.cn", date: time.Date(2019, 3, 17, 12, 48, 36, 0, time.UTC)},
		{filename: "google.com", date: time.Date(2020, 9, 14, 4, 0, 0, 0, time.UTC)},
		{filename: "ietf.org", date: time.Date(2020, 3, 12, 5, 0, 0, 0, time.UTC)},
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

		parsedDays, err := parse(cases[i].filename, ans)
		if parsedDays == 0 || err != nil {
			t.Errorf("%s got %.2f (error=%v) ", cases[i].filename, parsedDays, err)
		}

		answerDays := cases[i].date.Sub(time.Now()).Hours() / 24
		d := math.Abs(parsedDays-answerDays)
		if d < 0 || d > 1 {
			t.Errorf("cases[%d]: parsedDays=%.0f answerDays=%.0f", i, parsedDays, answerDays)
		}
	}
}
