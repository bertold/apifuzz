package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"
)

var validOffsets = map[string]bool{
	"d": true,
	"h": true,
	"m": true,
}

func TestDateHandler(t *testing.T) {

	tests := []struct {
		name          string
		offset        int
		unit          string
		errorExpected bool
	}{
		{"no arg", 0, "", false},
		{"unit only", 0, "d", false},
		{"offset only", 1, "", false},
		{"invalid unit", 1, "x", true},
		{"1 year ahead", 365, "d", false},
		{"1 year behind", -365, "d", false},
		{"1 day ahead", 1, "d", false},
		{"1 day behind", -1, "d", false},
		{"2 hours ahead", 2, "h", false},
		{"3 hours behind", -3, "h", false},
		{"4 minutes ahead", 4, "m", false},
		{"5 minutes behind", -5, "h", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			path := "/date"
			if test.offset != 0 {
				path += "?offset=" + strconv.Itoa(test.offset)
			}
			if test.unit != "" {
				if test.offset == 0 {
					path += "?unit=" + test.unit
				} else {
					path += "&unit=" + test.unit
				}
			}
			req, err := http.NewRequest("GET", path, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(DateHandler)
			handler.ServeHTTP(rr, req)
			status := rr.Code
			if (status == http.StatusOK && test.errorExpected) || (status != http.StatusOK && !test.errorExpected) {
				t.Errorf("handler returned wrong status code: got %v", status)
				t.Fail()
				return
			}

			if test.errorExpected {
				return
			}

			var expandedDate ExpandedDate
			err = json.Unmarshal(rr.Body.Bytes(), &expandedDate)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				t.Fail()
				return
			}

			multiplier := time.Duration(1 * time.Hour)
			if test.unit != "" {
				switch test.unit {
				case "d":
					multiplier = 24 * time.Hour
				case "h":
					multiplier = time.Hour
				case "m":
					multiplier = time.Minute
				case "s":
					multiplier = time.Second
				}
			}
			expectedDate := time.Now().Add(time.Duration(test.offset) * time.Duration(multiplier))
			if expandedDate.Year != expectedDate.Year() {
				t.Error("year mismatch")
				t.Fail()
			}
			if expandedDate.Month != int(expectedDate.Month()) {
				t.Error("month mismatch")
				t.Fail()
			}
			if expandedDate.Day != expectedDate.Day() {
				t.Error("day mismatch")
				t.Fail()
			}
			if expandedDate.Hour != expectedDate.Hour() {
				t.Error("hour mismatch")
				t.Fail()
			}
			if expandedDate.Minute != expectedDate.Minute() {
				t.Error("minute mismatch")
				t.Fail()
			}
		})
	}
}

func FuzzDateHandler(f *testing.F) {
	f.Fuzz(func(t *testing.T, offset int, unit string) {
		path := "/date"
		if offset != 0 {
			path += "?offset=" + strconv.Itoa(offset)
		}
		if unit != "" {
			if offset == 0 {
				path += "?unit=" + url.QueryEscape(unit)
			} else {
				path += "&unit=" + url.QueryEscape(unit)
			}
		}
		req, err := http.NewRequest("GET", path, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(DateHandler)
		handler.ServeHTTP(rr, req)
		status := rr.Code
		if status != http.StatusOK {
			if _, ok := validOffsets[unit]; ok {
				t.Errorf("handler returned wrong status code: got %v, offset: %d, unit: %s", status, offset, unit)
				t.Fail()
			}
			return
		}

		var expandedDate ExpandedDate
		err = json.Unmarshal(rr.Body.Bytes(), &expandedDate)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
			t.Fail()
			return
		}

		multiplier := time.Duration(1 * time.Hour)
		if unit != "" {
			switch unit {
			case "d":
				multiplier = 24 * time.Hour
			case "h":
				multiplier = time.Hour
			case "m":
				multiplier = time.Minute
			case "s":
				multiplier = time.Second
			}
		}
		expectedDate := time.Now().Add(time.Duration(offset) * time.Duration(multiplier))
		if expandedDate.Year != expectedDate.Year() {
			t.Error("year mismatch")
			t.Fail()
		}
		if expandedDate.Month != int(expectedDate.Month()) {
			t.Error("month mismatch")
			t.Fail()
		}
		if expandedDate.Day != expectedDate.Day() {
			t.Error("day mismatch")
			t.Fail()
		}
		if expandedDate.Hour != expectedDate.Hour() {
			t.Error("hour mismatch")
			t.Fail()
		}
		if expandedDate.Minute != expectedDate.Minute() {
			t.Error("minute mismatch")
			t.Fail()
		}
	})
}
