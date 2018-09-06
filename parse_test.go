package main

import (
	"io/ioutil"
	"path"
	"testing"

	"github.com/prometheus/common/promlog"
)

func init() { // so we don't panic
	level := promlog.AllowedLevel{}
	level.Set("info")
	logger = promlog.New(level)
}

func TestParsing(t *testing.T) {
	cases := []struct {
		filename, line string
	}{
		{filename: "google.cn", line: "Expiration Time: 2019-03-17 12:48:36"},
	}

	for i := range cases {
		w := path.Join("testdata", cases[i].filename)
		ans, err := ioutil.ReadFile(w)
		if err != nil {
			t.Errorf("problem on %s: %v", cases[i].filename, err)
		}

		n1, e1 := parse(cases[i].filename, ans)
		n2, e2 := parse(cases[i].filename, ans)

		if (n1 <= 0 || n2 <= 0) || n1 != n2 {
			t.Errorf("%s got %.2f (error=%v) and %.2f (error=%v)", cases[i].filename, n1, e1, n2, e2)
		}
	}
}
