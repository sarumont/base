// Copyright 2018 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package base

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"runtime"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	t1 := Now()

	// Verify time.Time methods work
	if diff := t1.Sub(t1.Time); diff != 0 {
		t.Errorf("got %v", diff)
	}
	if tt := time.Now().Add(1 * time.Second); t1.Sub(tt) == 0 {
		t.Error("expected difference in timing")
	}
}

func TestTime__JSON(t *testing.T) {
	// marshal and then unmarshal
	t1 := Now()

	bs, err := t1.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	var t2 Time
	if err := json.Unmarshal(bs, &t2); err != nil {
		t.Fatal(err)
	}
	if !t1.Equal(t2) {
		t.Errorf("unequal: t1=%q t2=%q", t1, t2)
	}

	in := []byte(`"2018-11-27T00:54:53Z"`)
	var t3 Time
	if err := json.Unmarshal(in, &t3); err != nil {
		t.Fatal(err)
	}
	if t3.IsZero() {
		t.Error("t3 shouldn't be zero time")
	}
}

var quote = []byte(`"`)

// TestTime__ruby will attempt to parse an ISO 8601 time generated by this library
func TestTime__ruby(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping ruby ISO 8601 test on windows")
	}

	bin, err := exec.LookPath("ruby")
	if err != nil || bin == "" {
		if inCI := os.Getenv("TRAVIS_OS_NAME") != ""; inCI {
			t.Fatal("ruby not found")
		} else {
			t.Skip("ruby not found")
		}
	}

	tt, err := time.Parse(iso8601Format, "2018-11-18T09:04:23-08:00")
	if err != nil {
		t.Fatal(err)
	}
	t1 := Time{
		Time: tt,
	}

	bs, err := t1.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	bs = bytes.TrimPrefix(bytes.TrimSuffix(bs, quote), quote)

	// Check with ruby
	cmd := exec.Command(bin, "time.rb", string(bs))
	bs, err = cmd.CombinedOutput()
	if err != nil {
		t.Errorf("err=%v\nOutput: %v", err, string(bs))
	}

	// Validate ruby output
	if !bytes.Contains(bs, []byte(`Date: 2018-11-18`)) {
		t.Errorf("no Date: %v", string(bs))
	}
	if !bytes.Contains(bs, []byte(`Time: 09:04:23`)) {
		t.Errorf("no Time: %v", string(bs))
	}
}
